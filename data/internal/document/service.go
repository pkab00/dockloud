package document

import (
	"context"
	"strings"

	apperrors "github.com/pkab00/dockloud/data/internal/errors"
)

type Service struct {
	documentRepo Repository
}

func NewService(documentRepo Repository) *Service {
	return &Service{
		documentRepo: documentRepo,
	}
}

func (s *Service) GetByAuthorID(ctx context.Context, authorID string, requesterID string) ([]Response, error) {
	var mapped []Response
	res, err := s.documentRepo.GetByAuthorID(ctx, authorID, requesterID)
	for _, doc := range res {
		mapped = append(mapped, *doc.toResponse())
	}
	return mapped, err
}

func (s *Service) GetAccessible(ctx context.Context, uid string) ([]Response, error) {
	var mapped []Response
	res, err := s.documentRepo.GetAccessible(ctx, uid)
	for _, doc := range res {
		mapped = append(mapped, *doc.toResponse())
	}
	return mapped, err
}

func (s *Service) GetByID(ctx context.Context, documentID int, requesterID string) (*Response, error) {
	res, err := s.documentRepo.GetByID(ctx, documentID, requesterID)
	return res.toResponse(), err
}

func (s *Service) Create(ctx context.Context, req *CreateRequest) (*Response, error) {
	if strings.TrimSpace(req.Name) == "" {
		return nil, apperrors.InvalidDocumentName
	}

	res, err := s.documentRepo.Create(ctx, req)
	return res.toResponse(), err
}

func (s *Service) Delete(ctx context.Context, documentID int, requesterID string) (*Response, error) {
	res, err := s.documentRepo.Delete(ctx, documentID, requesterID)
	return res.toResponse(), err
}

func (s *Service) Update(ctx context.Context, documentID int, requesterID string, req *UpdateRequest) (*Response, error) {
	res, err := s.documentRepo.Update(ctx, documentID, requesterID, req)
	return res.toResponse(), err
}

func (s *Service) GetUnsafe(ctx context.Context, documentID int) (*Response, error) {
	res, err := s.documentRepo.GetUnsafe(ctx, documentID)
	return res.toResponse(), err
}

func (s *Service) DeleteUnsafe(ctx context.Context, documentID int) (*Response, error) {
	res, err := s.documentRepo.DeleteUnsafe(ctx, documentID)
	return res.toResponse(), err
}

func (s *Service) UpdateUnsafe(ctx context.Context, documentID int, req *UpdateRequest) (*Response, error) {
	res, err := s.documentRepo.UpdateUnsafe(ctx, documentID, req)
	return res.toResponse(), err
}
