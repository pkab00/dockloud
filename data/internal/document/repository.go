package document

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

const documentKeyPrefix = "document"

var logger = log.Default()

type Repository interface {
	GetByAuthorID(ctx context.Context, authorID string, requesterID string) ([]Document, error)
	GetAccessible(ctx context.Context, uid string) ([]Document, error)
	GetByID(ctx context.Context, documentID int, requesterID string) (*Document, error)
	Create(ctx context.Context, req *CreateRequest) (*Document, error)
	Delete(ctx context.Context, documentID int, requesterID string) (*Document, error)
	Update(ctx context.Context, documentID int, requesterID string, req *UpdateRequest) (*Document, error)
}

type CreateRequest struct {
	OwnerID   string `json:"owner_id"`
	Name      string `json:"name"`
	Content   string `json:"content"`
	IsPrivate bool   `json:"is_private"`
}

type UpdateRequest struct {
	Name      *string `json:"name"`
	Content   *string `json:"content"`
	IsPrivate *bool   `json:"is_private"`
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

func (r *RepositoryWithCache) GetByAuthorID(
	ctx context.Context,
	authorID string,
	requesterID string,
) ([]Document, error) {
	var query string
	var err error

	query = `
		SELECT id, owner_id, name, content, is_private, is_deleted, created_at, edited_at
		FROM documents
		WHERE (owner_id = $1) AND (owner_id = $2 OR is_private = FALSE)
		AND is_deleted = FALSE
	`

	rows, err := r.db.QueryContext(ctx, query, authorID, requesterID)
	if err != nil {
		return nil, fmt.Errorf("error getting documents by author id: %w", err)
	}
	defer rows.Close()

	var documents []Document

	for rows.Next() {
		var d Document

		if err := rows.Scan(
			&d.ID,
			&d.OwnerID,
			&d.Name,
			&d.Content,
			&d.IsPrivate,
			&d.IsDeleted,
			&d.CreatedAt,
			&d.EditedAt,
		); err != nil {
			return nil, fmt.Errorf("error getting documents by author id: %w", err)
		}

		documents = append(documents, d)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error getting documents by author id: %w", err)
	}

	return documents, nil
}

func (r *RepositoryWithCache) GetAccessible(ctx context.Context, uid string) ([]Document, error) {
	var query string
	var err error

	query = `
		SELECT id, owner_id, name, content, is_private, is_deleted, created_at, edited_at
		FROM documents
		WHERE (owner_id = $1 OR is_private = FALSE)
		AND is_deleted = FALSE
	`

	rows, err := r.db.QueryContext(ctx, query, uid)
	if err != nil {
		return nil, fmt.Errorf("error getting accessable documents: %w", err)
	}
	defer rows.Close()

	var documents []Document

	for rows.Next() {
		var d Document

		if err := rows.Scan(
			&d.ID,
			&d.OwnerID,
			&d.Name,
			&d.Content,
			&d.IsPrivate,
			&d.IsDeleted,
			&d.CreatedAt,
			&d.EditedAt,
		); err != nil {
			return nil, fmt.Errorf("error getting accessable documents: %w", err)
		}

		documents = append(documents, d)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error getting accessable documents: %w", err)
	}

	return documents, nil
}

func (r *RepositoryWithCache) GetByID(ctx context.Context, documentID int, requesterID string) (*Document, error) {
	var d Document
	var err error

	documentKey := fmt.Sprintf("%s:%d", documentKeyPrefix, documentID)

	jsonDocument, err := r.cache.Get(ctx, documentKey).Result()

	switch {
	case err == nil:
		if err := json.Unmarshal([]byte(jsonDocument), &d); err != nil {
			_ = r.cache.Del(ctx, documentKey).Err()
			break
		}

	case errors.Is(err, redis.Nil):
		// cache miss

	default:
		logger.Println("cache is unavailable")
	}

	query := `
		SELECT id, owner_id, name, content, is_private, is_deleted, created_at, edited_at
		FROM documents
		WHERE id = $1
		  AND is_deleted = FALSE
	`

	err = r.db.
		QueryRowContext(ctx, query, documentID).
		Scan(
			&d.ID,
			&d.OwnerID,
			&d.Name,
			&d.Content,
			&d.IsPrivate,
			&d.IsDeleted,
			&d.CreatedAt,
			&d.EditedAt,
		)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperrors.DocumentNotFound
		}

		return nil, fmt.Errorf("error getting document by id: %w", err)
	}

	if !d.IsPrivate || d.OwnerID == requesterID {
		marshaledDocument, err := json.Marshal(d)
		if err != nil {
			logger.Println("unable to put value to cache")
			return &d, nil
		}

		err = r.cache.Set(
			ctx,
			documentKey,
			string(marshaledDocument),
			10*time.Minute,
		).Err()

		if err != nil {
			logger.Println("unable to put value to cache")
		}

		return &d, nil
	}

	return nil, apperrors.AccessDeclined
}

func (r *RepositoryWithCache) Create(ctx context.Context, req *CreateRequest) (*Document, error) {
	var d Document

	query := `
		INSERT INTO documents(owner_id, name, content, is_private)
		VALUES($1, $2, $3, $4)
		RETURNING id, owner_id, name, content, is_private, is_deleted, created_at, edited_at
	`

	err := r.db.
		QueryRowContext(
			ctx,
			query,
			req.OwnerID,
			req.Name,
			req.Content,
			req.IsPrivate,
		).
		Scan(
			&d.ID,
			&d.OwnerID,
			&d.Name,
			&d.Content,
			&d.IsPrivate,
			&d.IsDeleted,
			&d.CreatedAt,
			&d.EditedAt,
		)

	if err != nil {
		return nil, fmt.Errorf("error creating document: %w", err)
	}

	return &d, nil
}

func (r *RepositoryWithCache) Delete(ctx context.Context, documentID int, requesterID string) (*Document, error) {
	var d Document

	query := `
		UPDATE documents
		SET
			is_deleted = TRUE,
			edited_at = NOW()
		WHERE id = $1
		  AND owner_id = $2
		  AND is_deleted = FALSE
		RETURNING id, owner_id, name, content, is_private, is_deleted, created_at, edited_at
	`

	err := r.db.
		QueryRowContext(ctx, query, documentID, requesterID).
		Scan(
			&d.ID,
			&d.OwnerID,
			&d.Name,
			&d.Content,
			&d.IsPrivate,
			&d.IsDeleted,
			&d.CreatedAt,
			&d.EditedAt,
		)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperrors.DocumentNotFound
		}

		return nil, fmt.Errorf("error deleting document: %w", err)
	}

	documentKey := fmt.Sprintf("%s:%d", documentKeyPrefix, documentID)

	if err := r.cache.Del(ctx, documentKey).Err(); err != nil {
		logger.Println("error clearing old cache")
	}

	return &d, nil
}

func (r *RepositoryWithCache) Update(
	ctx context.Context,
	documentID int,
	requesterID string,
	req *UpdateRequest,
) (*Document, error) {
	var d Document

	query := `
		UPDATE documents
		SET
			name = COALESCE($1, name),
			content = COALESCE($2, content),
			is_private = COALESCE($3, is_private),
			edited_at = NOW()
		WHERE id = $4
		  AND owner_id = $5
		  AND is_deleted = FALSE
		RETURNING id, owner_id, name, content, is_private, is_deleted, created_at, edited_at
	`

	err := r.db.
		QueryRowContext(
			ctx,
			query,
			req.Name,
			req.Content,
			req.IsPrivate,
			documentID,
			requesterID,
		).
		Scan(
			&d.ID,
			&d.OwnerID,
			&d.Name,
			&d.Content,
			&d.IsPrivate,
			&d.IsDeleted,
			&d.CreatedAt,
			&d.EditedAt,
		)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperrors.DocumentNotFound
		}

		return nil, fmt.Errorf("error updating document: %w", err)
	}

	documentKey := fmt.Sprintf("%s:%d", documentKeyPrefix, documentID)

	if err := r.cache.Del(ctx, documentKey).Err(); err != nil {
		logger.Println("error clearing old cache")
	}

	return &d, nil
}
