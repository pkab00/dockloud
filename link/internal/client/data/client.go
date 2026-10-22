package data

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	apperror "github.com/pkab00/dockloud/link/internal/error"
	"github.com/pkab00/dockloud/link/internal/model"
	"github.com/pkab00/dockloud/shared/contract"
	sharedserviceauth "github.com/pkab00/dockloud/shared/serviceauth"
)

type Client struct {
	baseURL    string
	httpClient *http.Client
}

type InternalCreateLinkRequest struct {
	Token       string            `json:"token"`
	DocumentID  int               `json:"document_id"`
	Permissions model.Permissions `json:"permissions"`
	ExpiresAt   *time.Time        `json:"expires_at"`
}

type UpdateDocumentRequest struct {
	Name      *string `json:"name"`
	Content   *string `json:"content"`
	IsPrivate *bool   `json:"is_private"`
}

func NewClient(baseURL string) (*Client, error) {
	apiClaims := sharedserviceauth.Claims{
		Subject:  "dockloud-link",
		Issuer:   "dockloud-link",
		Audience: []string{"dockloud-data"},
	}

	token, err := sharedserviceauth.GenerateToken(
		os.Getenv("LINK_SERVICE_KEY"),
		apiClaims,
	)
	if err != nil {
		return nil, err
	}

	transport := &sharedserviceauth.AuthTransport{
		Base:   http.DefaultTransport,
		Token:  token,
		Claims: apiClaims,
	}

	return &Client{
		baseURL: baseURL,
		httpClient: &http.Client{
			Transport: transport,
			Timeout:   5 * time.Second,
		},
	}, nil
}

func checkResponseStatus(resp *http.Response, expectedStatus int) error {
	if resp.StatusCode == expectedStatus {
		return nil
	}

	var apiErr contract.ErrorResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiErr); err != nil {
		return fmt.Errorf(
			"data service returned status %d; failed to decode error response: %w",
			resp.StatusCode,
			err,
		)
	}

	return &apperror.ClientError{
		Status:  resp.StatusCode,
		Code:    apiErr.Code,
		Message: apiErr.Message,
	}
}

func marshalRequestBody(v any) (*bytes.Reader, error) {
	body, err := json.Marshal(v)
	if err != nil {
		return nil, fmt.Errorf("marshal data service request: %w", err)
	}
	return bytes.NewReader(body), nil
}

func (c *Client) newRequest(ctx context.Context, method string, url string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, fmt.Errorf("create data service request: %w", err)
	}
	return req, nil
}

// ========== LINKS ==========

func (c *Client) CreateLink(ctx context.Context, reqID string, request InternalCreateLinkRequest) (*model.Link, error) {
	var err error

	body, err := marshalRequestBody(request)
	if err != nil {
		return nil, fmt.Errorf("create data service request: %w", err)
	}

	httpRequest, err := c.newRequest(
		ctx,
		http.MethodPost,
		fmt.Sprintf("%s/links?requester_id=%s", c.baseURL, reqID),
		body,
	)
	if err != nil {
		return nil, fmt.Errorf("create data service request: %w", err)
	}

	response, err := c.httpClient.Do(httpRequest)
	if err != nil {
		return nil, fmt.Errorf("create data service request: %w", err)
	}
	defer response.Body.Close()

	if err := checkResponseStatus(response, http.StatusCreated); err != nil {
		return nil, fmt.Errorf("create data service request: %w", err)
	}

	var link model.Link
	if err := json.NewDecoder(response.Body).Decode(&link); err != nil {
		return nil, fmt.Errorf("create data service request: %w", err)
	}
	return &link, nil
}

func (c *Client) GetLinksList(ctx context.Context, creatorID string, documentID string) ([]model.Link, error) {
	var err error

	httpRequest, err := c.newRequest(
		ctx,
		http.MethodGet,
		fmt.Sprintf("%s/links?creator_id=%s&document_id=%s", c.baseURL, creatorID, documentID),
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("create data service request: %w", err)
	}

	response, err := c.httpClient.Do(httpRequest)
	if err != nil {
		return nil, fmt.Errorf("create data service request: %w", err)
	}
	defer response.Body.Close()

	if err := checkResponseStatus(response, http.StatusOK); err != nil {
		return nil, fmt.Errorf("create data service request: %w", err)
	}

	var link []model.Link
	if err := json.NewDecoder(response.Body).Decode(&link); err != nil {
		return nil, fmt.Errorf("create data service request: %w", err)
	}
	return link, nil
}

func (c *Client) GetLink(ctx context.Context, token string) (*model.Link, error) {
	var err error

	httpRequest, err := c.newRequest(
		ctx,
		http.MethodGet,
		fmt.Sprintf("%s/links/%s", c.baseURL, token),
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("create data service request: %w", err)
	}

	response, err := c.httpClient.Do(httpRequest)
	if err != nil {
		return nil, fmt.Errorf("create data service request: %w", err)
	}
	defer response.Body.Close()

	if err := checkResponseStatus(response, http.StatusOK); err != nil {
		return nil, fmt.Errorf("create data service request: %w", err)
	}

	var link model.Link
	if err := json.NewDecoder(response.Body).Decode(&link); err != nil {
		return nil, fmt.Errorf("create data service request: %w", err)
	}
	return &link, nil
}

func (c *Client) RevokeLink(ctx context.Context, token string, reqID string) (*model.Link, error) {
	var err error

	httpRequest, err := c.newRequest(
		ctx,
		http.MethodDelete,
		fmt.Sprintf("%s/links/%s?requester_id=%s", c.baseURL, token, reqID),
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("create data service request: %w", err)
	}

	response, err := c.httpClient.Do(httpRequest)
	if err != nil {
		return nil, fmt.Errorf("create data service request: %w", err)
	}
	defer response.Body.Close()

	if err := checkResponseStatus(response, http.StatusOK); err != nil {
		return nil, fmt.Errorf("create data service request: %w", err)
	}

	var link model.Link
	if err := json.NewDecoder(response.Body).Decode(&link); err != nil {
		return nil, fmt.Errorf("create data service request: %w", err)
	}
	return &link, nil
}

// ========== DOCUMENTS ==========

func (c *Client) GetDocument(ctx context.Context, documentID int) (*model.Document, error) {
	var err error

	httpRequest, err := c.newRequest(
		ctx,
		http.MethodGet,
		fmt.Sprintf("%s/unsafe/documents/%d", c.baseURL, documentID),
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("create data service request: %w", err)
	}

	response, err := c.httpClient.Do(httpRequest)
	if err != nil {
		return nil, fmt.Errorf("create data service request: %w", err)
	}
	defer response.Body.Close()

	if err := checkResponseStatus(response, http.StatusOK); err != nil {
		return nil, fmt.Errorf("create data service request: %w", err)
	}

	var document model.Document
	if err := json.NewDecoder(response.Body).Decode(&document); err != nil {
		return nil, fmt.Errorf("create data service request: %w", err)
	}
	return &document, nil
}

func (c *Client) UpdateDocument(ctx context.Context, req UpdateDocumentRequest, documentID int) (*model.Document, error) {
	var err error

	body, err := marshalRequestBody(req)
	if err != nil {
		return nil, fmt.Errorf("create data service request: %w", err)
	}

	httpRequest, err := c.newRequest(
		ctx,
		http.MethodPatch,
		fmt.Sprintf("%s/unsafe/documents/%d", c.baseURL, documentID),
		body,
	)
	if err != nil {
		return nil, fmt.Errorf("create data service request: %w", err)
	}

	response, err := c.httpClient.Do(httpRequest)
	if err != nil {
		return nil, fmt.Errorf("create data service request: %w", err)
	}
	defer response.Body.Close()

	if err := checkResponseStatus(response, http.StatusOK); err != nil {
		return nil, fmt.Errorf("create data service request: %w", err)
	}

	var document model.Document
	if err := json.NewDecoder(response.Body).Decode(&document); err != nil {
		return nil, fmt.Errorf("create data service request: %w", err)
	}
	return &document, nil
}

func (c *Client) DeleteDocument(ctx context.Context, documentID int) (*model.Document, error) {
	var err error

	httpRequest, err := c.newRequest(
		ctx,
		http.MethodDelete,
		fmt.Sprintf("%s/unsafe/documents/%d", c.baseURL, documentID),
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("create data service request: %w", err)
	}

	response, err := c.httpClient.Do(httpRequest)
	if err != nil {
		return nil, fmt.Errorf("create data service request: %w", err)
	}
	defer response.Body.Close()

	if err := checkResponseStatus(response, http.StatusOK); err != nil {
		return nil, fmt.Errorf("create data service request: %w", err)
	}

	var document model.Document
	if err := json.NewDecoder(response.Body).Decode(&document); err != nil {
		return nil, fmt.Errorf("create data service request: %w", err)
	}
	return &document, nil
}
