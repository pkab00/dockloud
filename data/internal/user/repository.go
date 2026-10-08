package user

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"

	apperrors "github.com/pkab00/dockloud/data/internal/errors"
	"github.com/redis/go-redis/v9"
)

const userKeyPrefix = "user"

var logger = log.Default()

type Repository interface {
	Find(ctx context.Context, req *UserFindRequest) (*UserFindResult, error)
	Authenticate(ctx context.Context, req *UserAuthRequest) (*bool, error)
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

func (r *RepositoryWithCache) Find(ctx context.Context, req *UserFindRequest) (*UserFindResult, error) {
	var u UserFindResult
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

func (r *RepositoryWithCache) Authenticate(ctx context.Context, req *UserAuthRequest) (*bool, error) {
	var err error
	var res bool

	query := `
	SELECT
		EXISTS (
			SELECT 1 FROM users
			WHERE 
				$1::text IS NOT NULL 
				AND email = $1::text
				AND $2::text IS NOT NULL 
				AND password_hash = $2::text
		)
	`

	err = r.db.
		QueryRowContext(ctx, query, req.Email, req.PasswordHash).
		Scan(&res)

	if err != nil {
		return nil, fmt.Errorf("error authenticating user: %w", err)
	}

	return &res, nil
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
