package document

import (
	"context"

	"github.com/pkab00/dockloud/api/internal/client/data"
	apperror "github.com/pkab00/dockloud/api/internal/error"
	"github.com/pkab00/dockloud/api/internal/middleware"
	"github.com/pkab00/dockloud/api/internal/model"
)

type ExternalCreateDocumentRequest struct {
	Name      string `json:"name"`
	Content   string `json:"content"`
	IsPrivate bool   `json:"is_private"`
}

type Service struct {
	client *data.Client
}

func getRequesterID(ctx context.Context) (*string, error) {
	reqID, ok := middleware.UIDFromContext(ctx)
	if !ok {
		return nil, apperror.FailedToGetRequesterID
	}
	return &reqID, nil
}

func NewService(client *data.Client) *Service {
	return &Service{client: client}
}

func (s *Service) Create(ctx context.Context, userReq ExternalCreateDocumentRequest) (*model.Document, error) {
	var err error

	reqID, err := getRequesterID(ctx)
	if err != nil {
		return nil, err
	}

	doc, err := s.client.CreateDocument(ctx, &data.InternalCreateDocumentRequest{
		OwnerID:   *reqID,
		Name:      userReq.Name,
		Content:   userReq.Content,
		IsPrivate: userReq.IsPrivate,
	})
	if err != nil {
		return nil, err
	}
	return doc, nil
}

func (s *Service) GetAll(ctx context.Context) ([]model.Document, error) {
	var err error

	reqID, err := getRequesterID(ctx)
	if err != nil {
		return nil, err
	}

	docs, err := s.client.GetAllDocuments(ctx, *reqID)
	if err != nil {
		return nil, err
	}
	return docs, nil
}

func (s *Service) GetByAuthor(ctx context.Context, authID string) ([]model.Document, error) {
	var err error

	reqID, err := getRequesterID(ctx)
	if err != nil {
		return nil, err
	}

	docs, err := s.client.GetDocumentsByAuthor(ctx, *reqID, authID)
	if err != nil {
		return nil, err
	}
	return docs, nil
}
