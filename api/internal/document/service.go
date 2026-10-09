package document

import (
	"context"

	"github.com/pkab00/dockloud/api/internal/client/data"
	"github.com/pkab00/dockloud/api/internal/model"
)

type CreateRequest struct {
	Name      string `json:"name"`
	Content   string `json:"content"`
	IsPrivate bool   `json:"is_private"`
}

type Service struct {
	client *data.Client
}

func NewService(client *data.Client) *Service {
	return &Service{client: client}
}

func (s *Service) Create(ctx context.Context, uid string, userReq CreateRequest) (*model.Document, error) {
	doc, err := s.client.CreateDocument(ctx, &data.CreateDocumentRequest{
		OwnerID:   uid,
		Name:      userReq.Name,
		Content:   userReq.Content,
		IsPrivate: userReq.IsPrivate,
	})
	if err != nil {
		return nil, err
	}
	return doc, nil
}
