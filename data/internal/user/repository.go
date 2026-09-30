package user

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/redis/go-redis/v9"
)

type Repository interface {
	GetAll(ctx context.Context) ([]User, error)
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
