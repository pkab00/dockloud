package link

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"time"

	"github.com/pkab00/dockloud/link/internal/client/data"
	"github.com/pkab00/dockloud/link/internal/model"
)

type ExternalCreateLinkRequest struct {
	DocumentID  int               `json:"document_id"`
	Permissions model.Permissions `json:"permissions"`
	ExpiresAt   *time.Time        `json:"expires_at"`
}

func generateToken() (string, error) {
	b := make([]byte, 16)

	if _, err := rand.Read(b); err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(b), nil
}

type Service struct {
	client *data.Client
}

func NewService(client *data.Client) *Service {
	return &Service{client: client}
}

func (s *Service) Create(ctx context.Context, reqID string, request ExternalCreateLinkRequest) (*model.Link, error) {
	token, err := generateToken()
	if err != nil {
		return nil, err
	}
	return s.client.CreateLink(ctx, reqID, data.InternalCreateLinkRequest{
		Token:       token,
		DocumentID:  request.DocumentID,
		Permissions: request.Permissions,
		ExpiresAt:   request.ExpiresAt,
	})
}

func (s *Service) Revoke(ctx context.Context, token string, reqID string) (*model.Link, error) {
	return s.client.RevokeLink(ctx, token, reqID)
}

func (s *Service) GetList(ctx context.Context, creatorID string, documentID string) ([]model.Link, error) {
	return s.client.GetLinksList(ctx, creatorID, documentID)
}
