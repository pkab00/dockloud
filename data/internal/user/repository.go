package user

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	apperrors "github.com/pkab00/dockloud/data/internal/errors"
	"github.com/redis/go-redis/v9"
)

const userKeyPrefix = "user"

var logger = log.Default()

type Repository interface {
	Get(ctx context.Context, uid string) (*User, error)
	Find(ctx context.Context, req *UserFindRequest) (*UserFindResponse, error)
	Authenticate(ctx context.Context, req *UserAuthRequest) (*string, error)
	Create(ctx context.Context, req *UserCreateRequest) (*User, error)
	Delete(ctx context.Context, uid string) (*User, error)
	Update(ctx context.Context, uid string, req *UserUpdateRequest) (*User, error)
}

type UserCreateRequest struct {
	Username     string `json:"username"`
	Email        string `json:"email"`
	PasswordHash string `json:"password_hash"`
}

type UserUpdateRequest struct {
	Username     *string `json:"username"`
	Email        *string `json:"email"`
	PasswordHash *string `json:"password_hash"`
}

type UserFindRequest struct {
	UID      *string `json:"uid"`
	Username *string `json:"username"`
	Email    *string `json:"email"`
}

type UserAuthRequest struct {
	Email        *string `json:"email"`
	PasswordHash *string `json:"password_hash"`
}

type RepositoryWithCache struct {
	db    *sql.DB
	cache *redis.Client
}

func NewRepositoryWithCache(db *sql.DB, cache *redis.Client) *RepositoryWithCache {
	return &RepositoryWithCache{
		db:    db,
		cache: cache,
	}
}

func (r *RepositoryWithCache) Find(ctx context.Context, req *UserFindRequest) (*UserFindResponse, error) {
	var u UserFindResponse
	var err error

	query := `
	SELECT
		EXISTS (
			SELECT 1 FROM users
			WHERE $1::uuid IS NOT NULL AND uid = $1::uuid
		) AS uid_exists,
		EXISTS (
			SELECT 1 FROM users
			WHERE $2::text IS NOT NULL AND username = $2::text
		) AS username_exists,
		EXISTS (
			SELECT 1 FROM users
			WHERE $3::text IS NOT NULL AND email = $3::text
		) AS email_exists;
	`

	err = r.db.
		QueryRowContext(ctx, query, req.UID, req.Username, req.Email).
		Scan(&u.UIDExists, &u.UsernameExists, &u.EmailExists)

	if err != nil {
		return nil, fmt.Errorf("error finding user: %w", err)
	}

	return &u, nil
}

func (r *RepositoryWithCache) Authenticate(ctx context.Context, req *UserAuthRequest) (*string, error) {
	var err error
	var uid string

	query := `
		SELECT uid FROM users
		WHERE 
			$1::text IS NOT NULL 
			AND email = $1::text
			AND $2::text IS NOT NULL 
			AND password_hash = $2::text
	`

	err = r.db.
		QueryRowContext(ctx, query, req.Email, req.PasswordHash).
		Scan(&uid)

	if err != nil {
		return nil, fmt.Errorf("error authenticating user: %w", err)
	}

	return &uid, nil
}

func (r *RepositoryWithCache) Get(ctx context.Context, uid string) (*User, error) {
	var u User
	var err error

	userKey := fmt.Sprintf("%s:%s", userKeyPrefix, uid)
	jsonUser, err := r.cache.Get(ctx, userKey).Result()

	switch {
	case err == nil:
		if err := json.Unmarshal([]byte(jsonUser), &u); err != nil {
			_ = r.cache.Del(ctx, userKey).Err()
			break
		}
		return &u, nil

	case errors.Is(err, redis.Nil):
		// cache miss

	default:
		logger.Println("cache is unavailable")
		break
	}

	query := `
		SELECT uid, username, email, password_hash, created_at
        FROM users
        WHERE uid = $1
	`

	err = r.db.
		QueryRowContext(ctx, query, uid).
		Scan(&u.UID, &u.Username, &u.Email, &u.PasswordHash, &u.CreatedAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperrors.UserNotFound
		}
		return nil, fmt.Errorf("error getting user by uid: %w", err)
	}

	marshaledUser, err := json.Marshal(u)
	if err != nil {
		logger.Println("unable to put value to cache")
		return &u, nil
	}
	err = r.cache.Set(ctx, userKey, string(marshaledUser), 10*time.Minute).Err()
	if err != nil {
		logger.Println("unable to put value to cache")
	}

	return &u, nil
}

func (r *RepositoryWithCache) Create(ctx context.Context, req *UserCreateRequest) (*User, error) {
	var u User
	var err error

	query := `
		INSERT INTO users(username, email, password_hash)
		VALUES($1, $2, $3)
		RETURNING uid, username, email, password_hash, created_at
	`

	err = r.db.
		QueryRowContext(ctx, query, req.Username, req.Email, req.PasswordHash).
		Scan(&u.UID, &u.Username, &u.Email, &u.PasswordHash, &u.CreatedAt)

	if err != nil {
		return nil, fmt.Errorf("error creating user: %w", err)
	}
	return &u, nil
}

func (r *RepositoryWithCache) Delete(ctx context.Context, uid string) (*User, error) {
	var u User
	var err error

	query := `
		DELETE FROM users
		WHERE uid = $1
		RETURNING uid, username, email, password_hash, created_at
	`

	err = r.db.
		QueryRowContext(ctx, query, uid).
		Scan(&u.UID, &u.Username, &u.Email, &u.PasswordHash, &u.CreatedAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperrors.UserNotFound
		}
		return nil, fmt.Errorf("error deleting user: %w", err)
	}

	userKey := fmt.Sprintf("%s:%s", userKeyPrefix, uid)
	err = r.cache.Del(ctx, userKey).Err()
	if err != nil {
		logger.Println("error clearing old cache")
	}

	return &u, nil
}

func (r *RepositoryWithCache) Update(ctx context.Context, uid string, req *UserUpdateRequest) (*User, error) {
	var u User
	var err error

	query := `
        UPDATE users
        SET
            username = COALESCE($1, username),
            email = COALESCE($2, email),
            password_hash = COALESCE($3, password_hash)
        WHERE uid = $4
		RETURNING uid, username, email, password_hash, created_at
    `

	err = r.db.
		QueryRowContext(
			ctx,
			query,
			req.Username,
			req.Email,
			req.PasswordHash,
			uid,
		).
		Scan(&u.UID, &u.Username, &u.Email, &u.PasswordHash, &u.CreatedAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperrors.UserNotFound
		}
		return nil, fmt.Errorf("error updating user: %w", err)
	}

	userKey := fmt.Sprintf("%s:%s", userKeyPrefix, uid)
	err = r.cache.Del(ctx, userKey).Err()
	if err != nil {
		logger.Println("error clearing old cache")
	}

	return &u, nil
}
