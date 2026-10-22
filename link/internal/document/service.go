package document

import (
	"context"

	"github.com/pkab00/dockloud/link/internal/client/data"
	apperror "github.com/pkab00/dockloud/link/internal/error"
	"github.com/pkab00/dockloud/link/internal/model"
)

type Service struct {
	client *data.Client
}

func NewService(client *data.Client) *Service {
	return &Service{client: client}
}

func (s *Service) Get(ctx context.Context, token string) (*model.Document, error) {
	link, err := s.client.GetLink(ctx, token)
	if err != nil {
		return nil, err
	}

	switch {
	case link.Expired():
		return nil, apperror.ErrTokenExpired
	case link.Revoked():
		return nil, apperror.ErrTokenRevoked
	case !*link.Permissions.Read:
		return nil, apperror.ErrNoPermission
	default:
		doc, err := s.client.GetDocument(ctx, link.DocumentID)
		if err != nil {
			return nil, err
		}
		return doc, nil
	}
}

func (s *Service) Update(ctx context.Context, token string, req data.UpdateDocumentRequest) (*model.Document, error) {
	link, err := s.client.GetLink(ctx, token)
	if err != nil {
		return nil, err
	}

	switch {
	case link.Expired():
		return nil, apperror.ErrTokenExpired
	case link.Revoked():
		return nil, apperror.ErrTokenRevoked
	case !*link.Permissions.Write:
		return nil, apperror.ErrNoPermission
	default:
		doc, err := s.client.UpdateDocument(ctx, req, link.DocumentID)
		if err != nil {
			return nil, err
		}
		return doc, nil
	}
}

func (s *Service) Delete(ctx context.Context, token string) (*model.Document, error) {
	link, err := s.client.GetLink(ctx, token)
	if err != nil {
		return nil, err
	}

	switch {
	case link.Expired():
		return nil, apperror.ErrTokenExpired
	case link.Revoked():
		return nil, apperror.ErrTokenRevoked
	case !*link.Permissions.Delete:
		return nil, apperror.ErrNoPermission
	default:
		doc, err := s.client.DeleteDocument(ctx, link.DocumentID)
		if err != nil {
			return nil, err
		}
		return doc, nil
	}
}
