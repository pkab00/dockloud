package accesslink

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	apperrors "github.com/pkab00/dockloud/data/internal/errors"
)

var logger = log.Default()

func scanLinkFromRows(rows *sql.Rows, al *AccessLink) error {
	return rows.Scan(
		&al.Token, &al.DocumentID, &al.CreatorID,
		&al.Permissions.Read, &al.Permissions.Write, &al.Permissions.Delete,
		&al.CreatedAt, &al.ExpiresAt, &al.RevokedAt,
	)
}

func scanLinkFromRow(row *sql.Row, al *AccessLink) error {
	return row.Scan(
		&al.Token, &al.DocumentID, &al.CreatorID,
		&al.Permissions.Read, &al.Permissions.Write, &al.Permissions.Delete,
		&al.CreatedAt, &al.ExpiresAt, &al.RevokedAt,
	)
}

type CreateLinkRequest struct {
	Token       string            `json:"token"`
	DocumentID  int               `json:"document_id"`
	Permissions AccessPermissions `json:"permissions"`
	ExpiresAt   *time.Time        `json:"expires_at"`
}

type Repository interface {
	GetAll(ctx context.Context, creatorID string) ([]AccessLink, error)
	GetAllForDocument(ctx context.Context, creatorID string, docID int) ([]AccessLink, error)
	GetByToken(ctx context.Context, token string) (*AccessLink, error)
	Create(ctx context.Context, reqID string, req *CreateLinkRequest) (*AccessLink, error)
	Revoke(ctx context.Context, reqID string, token string) (*AccessLink, error)
}

type PostgresRepository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{
		db: db,
	}
}

func (r *PostgresRepository) GetAll(ctx context.Context, creatorID string) ([]AccessLink, error) {
	var query string
	var err error

	query = `
		SELECT token, document_id, creator_id,
		read_allowed, write_allowed, delete_allowed,
		created_at, expires_at, revoked_at
		FROM access_links
		WHERE creator_id = $1
		AND (revoked_at IS NULL OR expires_at > now())
	`

	rows, err := r.db.QueryContext(ctx, query, creatorID)
	if err != nil {
		return nil, fmt.Errorf("error getting accessable links: %w", err)
	}
	defer rows.Close()

	var links []AccessLink

	for rows.Next() {
		var al AccessLink

		if err := scanLinkFromRows(rows, &al); err != nil {
			return nil, fmt.Errorf("error getting accessable links: %w", err)
		}

		links = append(links, al)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error getting accessable links: %w", err)
	}

	return links, nil
}

func (r *PostgresRepository) GetAllForDocument(ctx context.Context, creatorID string, docID int) ([]AccessLink, error) {
	var query string
	var err error

	query = `
		SELECT token, document_id, creator_id,
		read_allowed, write_allowed, delete_allowed,
		created_at, expires_at, revoked_at
		FROM access_links
		WHERE creator_id = $1 AND document_id = $2
		AND (revoked_at IS NULL OR expires_at > now())
	`

	rows, err := r.db.QueryContext(ctx, query, creatorID, docID)
	if err != nil {
		return nil, fmt.Errorf("error getting accessable links: %w", err)
	}
	defer rows.Close()

	var links []AccessLink

	for rows.Next() {
		var al AccessLink

		if err := scanLinkFromRows(rows, &al); err != nil {
			return nil, fmt.Errorf("error getting accessable links: %w", err)
		}

		links = append(links, al)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error getting accessable links: %w", err)
	}

	return links, nil
}

func (r *PostgresRepository) GetByToken(ctx context.Context, token string) (*AccessLink, error) {
	var al AccessLink
	var query string
	var err error

	query = `
		SELECT token, document_id, creator_id,
		read_allowed, write_allowed, delete_allowed,
		created_at, expires_at, revoked_at
		FROM access_links
		WHERE token = $1
		LIMIT 1
	`

	err = scanLinkFromRow(r.db.QueryRowContext(ctx, query, token), &al)
	if err != nil {
		return nil, fmt.Errorf("error getting link by token: %w", err)
	}

	if al.isActiveLink() {
		return &al, nil
	}

	return nil, apperrors.AccessLinkExpired
}

func (r *PostgresRepository) Create(ctx context.Context, reqID string, req *CreateLinkRequest) (*AccessLink, error) {
	var al AccessLink
	var err error

	query := `
		INSERT INTO access_links(token, document_id, creator_id,
		read_allowed, write_allowed, delete_allowed, expires_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING token, document_id, creator_id,
		read_allowed, write_allowed, delete_allowed,
		created_at, expires_at, revoked_at
	`

	err = scanLinkFromRow(
		r.db.QueryRowContext(
			ctx, query, req.Token, req.DocumentID, reqID,
			req.Permissions.Read, req.Permissions.Write, req.Permissions.Delete, req.ExpiresAt,
		), &al,
	)
	if err != nil {
		return nil, fmt.Errorf("error creating link: %w", err)
	}
	return &al, nil
}

func (r *PostgresRepository) Revoke(ctx context.Context, reqID string, token string) (*AccessLink, error) {
	var al AccessLink
	var err error

	query := `
		UPDATE access_links
		SET revoked_at = now()
		WHERE token = $1 AND creator_id = $2
		RETURNING token, document_id, creator_id,
		read_allowed, write_allowed, delete_allowed,
		created_at, expires_at, revoked_at
	`

	err = scanLinkFromRow(r.db.QueryRowContext(ctx, query, token, reqID), &al)
	if err != nil {
		return nil, fmt.Errorf("error deleting link: %w", err)
	}

	return &al, nil
}
