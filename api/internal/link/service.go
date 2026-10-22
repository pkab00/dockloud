package link

import (
	"context"

	"github.com/pkab00/dockloud/api/internal/client/link"
	apperror "github.com/pkab00/dockloud/api/internal/error"
	"github.com/pkab00/dockloud/api/internal/middleware"
	"github.com/pkab00/dockloud/api/internal/model"
)

type Service struct {
	client *link.Client
}

func NewService(client *link.Client) *Service {
	return &Service{client: client}
}

func (s *Service) Create(ctx context.Context, req *link.CreateLinkRequest) (*model.Link, error) {
	reqID, ok := middleware.UIDFromContext(ctx)
	if !ok {
		return nil, apperror.FailedToGetRequesterID
	}
	res, err := s.client.CreateLink(ctx, req, reqID)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (s *Service) Revoke(ctx context.Context, token string) (*model.Link, error) {
	reqID, ok := middleware.UIDFromContext(ctx)
	if !ok {
		return nil, apperror.FailedToGetRequesterID
	}
	res, err := s.client.RevokeLink(ctx, token, reqID)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (s *Service) Get(ctx context.Context, documentID string) ([]model.Link, error) {
	reqID, ok := middleware.UIDFromContext(ctx)
	if !ok {
		return nil, apperror.FailedToGetRequesterID
	}
	res, err := s.client.GetLinks(ctx, reqID, documentID)
	if err != nil {
		return nil, err
	}
	return res, nil
}
