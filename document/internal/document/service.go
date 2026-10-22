package document

import (
	"context"

	"github.com/pkab00/dockloud/document/internal/client/data"
	"github.com/pkab00/dockloud/document/internal/client/link"
	apperror "github.com/pkab00/dockloud/document/internal/error"
	"github.com/pkab00/dockloud/document/internal/middleware"
	"github.com/pkab00/dockloud/document/internal/model"
)

type ExternalCreateDocumentRequest struct {
	Name      string `json:"name"`
	Content   string `json:"content"`
	IsPrivate *bool  `json:"is_private"`
}

type Service struct {
	dataClient *data.Client
	linkClient *link.Client
}

func getRequesterID(ctx context.Context) (*string, error) {
	reqID, ok := middleware.UIDFromContext(ctx)
	if !ok {
		return nil, apperror.FailedToGetRequesterID
	}
	return &reqID, nil
}

func NewService(dataClient *data.Client, linkClient *link.Client) *Service {
	return &Service{
		dataClient: dataClient,
		linkClient: linkClient,
	}
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

	doc, err := s.dataClient.CreateDocument(ctx, &data.InternalCreateDocumentRequest{
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

	doc, err := s.dataClient.UpdateDocument(ctx, docID, *reqID, userReq)
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

	docs, err := s.dataClient.GetAllDocuments(ctx, *reqID)
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

	docs, err := s.dataClient.GetDocumentsByAuthor(ctx, *reqID, authID)
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

	doc, err := s.dataClient.GetDocumentByID(ctx, docID, *reqID)
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

	docs, err := s.dataClient.DeleteDocument(ctx, docID, *reqID)
	if err != nil {
		return nil, err
	}
	return docs, nil
}

func (s *Service) GetShared(ctx context.Context, token string) (*model.Document, error) {
	doc, err := s.linkClient.GetDocument(ctx, token)
	if err != nil {
		return nil, err
	}
	return doc, nil
}

func (s *Service) UpdateShared(ctx context.Context, token string, req *data.UpdateDocumentRequest) (*model.Document, error) {
	doc, err := s.linkClient.UpdateDocument(ctx, token, req)
	if err != nil {
		return nil, err
	}
	return doc, nil
}

func (s *Service) DeleteShared(ctx context.Context, token string) (*model.Document, error) {
	doc, err := s.linkClient.DeleteDocument(ctx, token)
	if err != nil {
		return nil, err
	}
	return doc, nil
}
