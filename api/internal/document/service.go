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
	IsPrivate *bool  `json:"is_private"`
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

	var isPrivate bool
	if userReq.IsPrivate == nil {
		isPrivate = true
	} else {
		isPrivate = *userReq.IsPrivate
	}

	doc, err := s.client.CreateDocument(ctx, &data.InternalCreateDocumentRequest{
		OwnerID:   *reqID,
		Name:      userReq.Name,
		Content:   userReq.Content,
		IsPrivate: isPrivate,
	})
	if err != nil {
		return nil, err
	}
	return doc, nil
}

func (s *Service) Update(ctx context.Context, docID int, userReq *data.UpdateDocumentRequest) (*model.Document, error) {
	var err error

	reqID, err := getRequesterID(ctx)
	if err != nil {
		return nil, err
	}

	doc, err := s.client.UpdateDocument(ctx, docID, *reqID, userReq)
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

func (s *Service) GetByID(ctx context.Context, docID int) (*model.Document, error) {
	var err error

	reqID, err := getRequesterID(ctx)
	if err != nil {
		return nil, err
	}

	doc, err := s.client.GetDocumentByID(ctx, docID, *reqID)
	if err != nil {
		return nil, err
	}
	return doc, nil
}

func (s *Service) Delete(ctx context.Context, docID int) (*model.Document, error) {
	var err error

	reqID, err := getRequesterID(ctx)
	if err != nil {
		return nil, err
	}

	docs, err := s.client.DeleteDocument(ctx, docID, *reqID)
	if err != nil {
		return nil, err
	}
	return docs, nil
}
