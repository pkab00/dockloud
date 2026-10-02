package document

import (
	"context"
	"strings"

	apperrors "github.com/pkab00/dockloud/data/internal/errors"
	"github.com/pkab00/dockloud/data/internal/user"
)

type Service struct {
	documentRepo Repository
	userRepo     user.Repository
}

func NewService(documentRepo Repository, userRepo user.Repository) *Service {
	return &Service{
		documentRepo: documentRepo,
		userRepo:     userRepo,
	}
}

func (s *Service) GetAll(ctx context.Context) ([]Document, error) {
	return s.documentRepo.GetAll(ctx)
}

func (s *Service) GetByUID(ctx context.Context, uid string) ([]Document, error) {
	_, err := s.userRepo.GetByUID(ctx, uid)
	if err != nil {
		return nil, apperrors.UserNotFound
	}
	return s.documentRepo.GetByUID(ctx, uid)
}

func (s *Service) GetByID(ctx context.Context, id int) (*Document, error) {
	return s.documentRepo.GetByID(ctx, id)
}

func (s *Service) Create(ctx context.Context, req *DocumentCreateRequest) (*Document, error) {
	if strings.TrimSpace(req.Name) == "" {
		return nil, apperrors.InvalidDocumentName
	}

	return s.documentRepo.Create(ctx, req)
}

func (s *Service) Delete(ctx context.Context, id int) (*Document, error) {
	return s.documentRepo.Delete(ctx, id)
}

func (s *Service) Update(
	ctx context.Context,
	id int,
	req *DocumentUpdateRequest,
) (*Document, error) {
	return s.documentRepo.Update(ctx, id, req)
}
