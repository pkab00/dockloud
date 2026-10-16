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

func (s *Service) GetByAuthorID(ctx context.Context, authorID string, requesterID string) ([]DocumentResponse, error) {
	var mapped []DocumentResponse
	res, err := s.documentRepo.GetByAuthorID(ctx, authorID, requesterID)
	for _, doc := range res {
		mapped = append(mapped, *doc.toResponse())
	}
	return mapped, err
}

func (s *Service) GetAccessible(ctx context.Context, uid string) ([]DocumentResponse, error) {
	var mapped []DocumentResponse
	res, err := s.documentRepo.GetAccessible(ctx, uid)
	for _, doc := range res {
		mapped = append(mapped, *doc.toResponse())
	}
	return mapped, err
}

func (s *Service) GetByID(ctx context.Context, documentID int, requesterID string) (*DocumentResponse, error) {
	res, err := s.documentRepo.GetByID(ctx, documentID, requesterID)
	return res.toResponse(), err
}

func (s *Service) Create(ctx context.Context, req *CreateRequest) (*DocumentResponse, error) {
	if strings.TrimSpace(req.Name) == "" {
		return nil, apperrors.InvalidDocumentName
	}

	res, err := s.documentRepo.Create(ctx, req)
	return res.toResponse(), err
}

func (s *Service) Delete(ctx context.Context, documentID int, requesterID string) (*DocumentResponse, error) {
	res, err := s.documentRepo.Delete(ctx, documentID, requesterID)
	return res.toResponse(), err
}

func (s *Service) Update(
	ctx context.Context,
	documentID int,
	requesterID string,
	req *UpdateRequest,
) (*DocumentResponse, error) {
	res, err := s.documentRepo.Update(ctx, documentID, requesterID, req)
	return res.toResponse(), err
}
