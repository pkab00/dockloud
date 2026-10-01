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
	GetAll(ctx context.Context) ([]User, error)
	GetByUID(ctx context.Context, uid string) (*User, error)
	Create(ctx context.Context, req *UserCreateRequest) (*User, error)
}

type UserCreateRequest struct {
	Username     string `json:"username"`
	Email        string `json:"email"`
	PasswordHash string `json:"password_hash"`
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

func (r *RepositoryWithCache) GetAll(ctx context.Context) ([]User, error) {
	query := "SELECT * FROM users"
	row, err := r.db.QueryContext(ctx, query)

	if err != nil {
		return nil, fmt.Errorf("error getting all users: %w", err)
	}

	var users []User
	for row.Next() {
		var u User
		if err := row.Scan(&u.UID, &u.Username, &u.Email, &u.PasswordHash, &u.CreatedAt); err != nil {
			return nil, fmt.Errorf("error getting all users: %w", err)
		}
		users = append(users, u)
	}

	if err := row.Err(); err != nil {
		return nil, fmt.Errorf("error getting all users: %w", err)
	}

	return users, nil
}

func (r *RepositoryWithCache) GetByUID(ctx context.Context, uid string) (*User, error) {
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
