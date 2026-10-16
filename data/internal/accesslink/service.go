package accesslink

import (
	"context"
	"database/sql"
	"errors"
	"time"

	apperrors "github.com/pkab00/dockloud/data/internal/errors"
)

type Service struct {
	repo Repository
}

type ExternalCreateLinkRequest struct {
	DocumentID  int               `json:"document_id"`
	Permissions AccessPermissions `json:"permissions"`
	ExpiresAt   *time.Time        `json:"expires_at"`
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetAll(ctx context.Context, creatorID string) ([]AccessLink, error) {
	res, err := s.repo.GetAll(ctx, creatorID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperrors.NoLinksFound
		} else {
			return nil, err
		}
	}
	return res, nil
}

func (s *Service) GetAllForDocument(ctx context.Context, creatorID string, docID int) ([]AccessLink, error) {
	res, err := s.repo.GetAllForDocument(ctx, creatorID, docID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperrors.NoLinksFound
		} else {
			return nil, err
		}
	}
	return res, nil

}

func (s *Service) GetByToken(ctx context.Context, token string) (*AccessLink, error) {
	res, err := s.repo.GetByToken(ctx, token)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperrors.NoLinksFound
		} else {
			return nil, err
		}
	}
	return res, nil
}

func (s *Service) Create(ctx context.Context, reqID string, req *CreateLinkRequest) (*AccessLink, error) {
	return s.repo.Create(ctx, reqID, req)
}

func (s *Service) Revoke(ctx context.Context, reqID string, token string) (*AccessLink, error) {
	return s.repo.Revoke(ctx, reqID, token)
}
