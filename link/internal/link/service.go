package link

import (
	"context"

	"github.com/pkab00/dockloud/link/internal/client/data"
	"github.com/pkab00/dockloud/link/internal/model"
)

type Service struct {
	client data.Client
}

func NewService(client data.Client) *Service {
	return &Service{client: client}
}

func (s *Service) Create(ctx context.Context, reqID string, request data.CreateLinkRequest) (*model.Link, error) {
	return s.client.CreateLink(ctx, reqID, request)
}

func (s *Service) Revoke(ctx context.Context, token string, reqID string) (*model.Link, error) {
	return s.client.RevokeLink(ctx, token, reqID)
}

func (s *Service) Get(ctx context.Context, reqID string) (*model.Link, error) {
	return s.client.GetLink(ctx, reqID)
}

func (s *Service) ListForDocument(ctx context.Context, creatorID string, documentID int) ([]model.Link, error) {
	return s.client.GetLinksForDocument(ctx, creatorID, documentID)
}

func (s *Service) ListByCreator(ctx context.Context, creatorID string) ([]model.Link, error) {
	return s.client.GetLinksByCreator(ctx, creatorID)
}
