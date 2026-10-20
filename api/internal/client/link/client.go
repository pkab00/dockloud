package link

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/pkab00/dockloud/api/internal/client"
	"github.com/pkab00/dockloud/api/internal/client/data"
	"github.com/pkab00/dockloud/api/internal/model"
)

type Client struct {
	*client.Client
}

func NewClient(baseURL string) (*Client, error) {
	c, err := client.NewClient(baseURL)
	if err != nil {
		return nil, err
	}
	return &Client{c}, nil
}

type CreateLinkRequest struct {
	DocumentID  int               `json:"document_id"`
	Permissions model.Permissions `json:"permissions"`
	ExpiresAt   *time.Time        `json:"expires_at"`
}

// ========== LINKS ==========

func (c *Client) CreateLink(ctx context.Context, req *CreateLinkRequest, reqID string) (*model.Link, error) {
	var err error
	body, err := client.MarshalRequestBody(req)
	if err != nil {
		return nil, err
	}

	clientReq, err := c.NewRequest(
		ctx,
		http.MethodPost,
		fmt.Sprintf("%s/links?requester_id=%s", c.BaseURL, reqID),
		body,
	)
	if err != nil {
		return nil, err
	}

	resp, err := c.HTTPClient.Do(clientReq)
	if err != nil {
		return nil, fmt.Errorf("send create link request: %w", err)
	}
	defer resp.Body.Close()

	if err := client.CheckResponseStatus(resp, http.StatusCreated); err != nil {
		return nil, fmt.Errorf("check create link response: %w", err)
	}

	var link model.Link
	if err := json.NewDecoder(resp.Body).Decode(&link); err != nil {
		return nil, fmt.Errorf("decode create link response: %w", err)
	}

	return &link, nil
}

func (c *Client) RevokeLink(ctx context.Context, token string, reqID string) (*model.Link, error) {
	var err error

	clientReq, err := c.NewRequest(
		ctx,
		http.MethodDelete,
		fmt.Sprintf("%s/links/%s?requester_id=%s", c.BaseURL, token, reqID),
		nil,
	)
	if err != nil {
		return nil, err
	}

	resp, err := c.HTTPClient.Do(clientReq)
	if err != nil {
		return nil, fmt.Errorf("send delete link request: %w", err)
	}
	defer resp.Body.Close()

	if err := client.CheckResponseStatus(resp, http.StatusOK); err != nil {
		return nil, fmt.Errorf("check delete link response: %w", err)
	}

	var link model.Link
	if err := json.NewDecoder(resp.Body).Decode(&link); err != nil {
		return nil, fmt.Errorf("decode delete link response: %w", err)
	}

	return &link, nil
}

func (c *Client) GetLinks(ctx context.Context, creatorID string, documentID string) ([]model.Link, error) {
	var err error

	clientReq, err := c.NewRequest(
		ctx,
		http.MethodGet,
		fmt.Sprintf("%s/links?creator_id=%s&document_id=%s", c.BaseURL, creatorID, documentID),
		nil,
	)
	if err != nil {
		return nil, err
	}

	resp, err := c.HTTPClient.Do(clientReq)
	if err != nil {
		return nil, fmt.Errorf("send get link request: %w", err)
	}
	defer resp.Body.Close()

	if err := client.CheckResponseStatus(resp, http.StatusOK); err != nil {
		return nil, fmt.Errorf("check get link response: %w", err)
	}

	var links []model.Link
	if err := json.NewDecoder(resp.Body).Decode(&links); err != nil {
		return nil, fmt.Errorf("decode get link response: %w", err)
	}

	return links, nil
}

// ========== DOCUMENTS ==========

func (c *Client) GetDocument(ctx context.Context, token string) (*model.Document, error) {
	var err error

	clientReq, err := c.NewRequest(
		ctx,
		http.MethodGet,
		fmt.Sprintf("%s/documents/%s", c.BaseURL, token),
		nil,
	)
	if err != nil {
		return nil, err
	}

	resp, err := c.HTTPClient.Do(clientReq)
	if err != nil {
		return nil, fmt.Errorf("send get link request: %w", err)
	}
	defer resp.Body.Close()

	if err := client.CheckResponseStatus(resp, http.StatusOK); err != nil {
		return nil, fmt.Errorf("check get link response: %w", err)
	}

	var document model.Document
	if err := json.NewDecoder(resp.Body).Decode(&document); err != nil {
		return nil, fmt.Errorf("decode get link response: %w", err)
	}

	return &document, nil
}

func (c *Client) UpdateDocument(ctx context.Context, token string, request *data.UpdateDocumentRequest) (*model.Document, error) {
	var err error

	body, err := client.MarshalRequestBody(request)
	if err != nil {
		return nil, err
	}

	clientReq, err := c.NewRequest(
		ctx,
		http.MethodPatch,
		fmt.Sprintf("%s/documents/%s", c.BaseURL, token),
		body,
	)
	if err != nil {
		return nil, err
	}

	resp, err := c.HTTPClient.Do(clientReq)
	if err != nil {
		return nil, fmt.Errorf("send get link request: %w", err)
	}
	defer resp.Body.Close()

	if err := client.CheckResponseStatus(resp, http.StatusOK); err != nil {
		return nil, fmt.Errorf("check get link response: %w", err)
	}

	var document model.Document
	if err := json.NewDecoder(resp.Body).Decode(&document); err != nil {
		return nil, fmt.Errorf("decode get link response: %w", err)
	}

	return &document, nil
}

func (c *Client) DeleteDocument(ctx context.Context, token string) (*model.Document, error) {
	var err error

	clientReq, err := c.NewRequest(
		ctx,
		http.MethodDelete,
		fmt.Sprintf("%s/documents/%s", c.BaseURL, token),
		nil,
	)
	if err != nil {
		return nil, err
	}

	resp, err := c.HTTPClient.Do(clientReq)
	if err != nil {
		return nil, fmt.Errorf("send get link request: %w", err)
	}
	defer resp.Body.Close()

	if err := client.CheckResponseStatus(resp, http.StatusOK); err != nil {
		return nil, fmt.Errorf("check get link response: %w", err)
	}

	var document model.Document
	if err := json.NewDecoder(resp.Body).Decode(&document); err != nil {
		return nil, fmt.Errorf("decode get link response: %w", err)
	}

	return &document, nil
}
