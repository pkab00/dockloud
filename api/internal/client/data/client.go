package data

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	apperror "github.com/pkab00/dockloud/api/internal/error"
	"github.com/pkab00/dockloud/api/internal/model"
	"github.com/pkab00/dockloud/api/internal/serviceauth"
	"github.com/pkab00/dockloud/shared/contract"
)

type InternalCreateUserRequest struct {
	Username     string `json:"username"`
	Email        string `json:"email"`
	PasswordHash string `json:"password_hash"`
}

type InternalUpdateUserRequest struct {
	Username     *string `json:"username"`
	Email        *string `json:"email"`
	PasswordHash *string `json:"password_hash"`
}

type InternalAuthUserRequest struct {
	Email        *string `json:"email"`
	PasswordHash *string `json:"password_hash"`
}

type FindUserRequest struct {
	UID      *string `json:"uid"`
	Username *string `json:"username"`
	Email    *string `json:"email"`
}

type InternalCreateDocumentRequest struct {
	OwnerID   string `json:"owner_id"`
	Name      string `json:"name"`
	Content   string `json:"content"`
	IsPrivate bool   `json:"is_private"`
}

type UpdateDocumentRequest struct {
	Name      *string `json:"name"`
	Content   *string `json:"content"`
	IsPrivate *bool   `json:"is_private"`
}

type Client struct {
	baseURL    string
	httpClient *http.Client
}

func NewClient(baseURL string) (*Client, error) {
	token, err := serviceauth.GenerateToken()
	if err != nil {
		return nil, err
	}

	transport := &serviceauth.AuthTransport{
		Base:  http.DefaultTransport,
		Token: token,
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

func (c *Client) newRequest(
	ctx context.Context,
	method string,
	url string,
	body io.Reader,
) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, fmt.Errorf("create data service request: %w", err)
	}
	return req, nil
}

// ========================== USERS ==========================

func (c *Client) FindUser(ctx context.Context, req *FindUserRequest) (*model.UserFindResult, error) {
	body, err := marshalRequestBody(req)
	if err != nil {
		return nil, err
	}

	dataReq, err := c.newRequest(
		ctx,
		http.MethodGet,
		fmt.Sprintf("%s/users/find", c.baseURL),
		body,
	)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(dataReq)
	if err != nil {
		return nil, fmt.Errorf("send find user request to data service: %w", err)
	}
	defer resp.Body.Close()

	if err := checkResponseStatus(resp, http.StatusOK); err != nil {
		log.Default().Printf("data client: %v", err)
		return nil, err
	}

	var result model.UserFindResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode find user response: %w", err)
	}

	return &result, nil
}

func (c *Client) AuthUser(ctx context.Context, req *InternalAuthUserRequest) (*model.User, error) {
	body, err := marshalRequestBody(req)
	if err != nil {
		return nil, err
	}

	dataReq, err := c.newRequest(
		ctx,
		http.MethodGet,
		fmt.Sprintf("%s/users/auth", c.baseURL),
		body,
	)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(dataReq)
	if err != nil {
		return nil, fmt.Errorf("send authenticate user request to data service: %w", err)
	}
	defer resp.Body.Close()

	if err := checkResponseStatus(resp, http.StatusOK); err != nil {
		log.Default().Printf("data client: %v", err)
		return nil, err
	}

	var user model.User
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, fmt.Errorf("decode authenticate user response: %w", err)
	}

	return &user, nil
}

func (c *Client) CreateUser(ctx context.Context, inReq *InternalCreateUserRequest) (*model.User, error) {
	body, err := marshalRequestBody(inReq)
	if err != nil {
		return nil, err
	}

	req, err := c.newRequest(
		ctx,
		http.MethodPost,
		fmt.Sprintf("%s/users", c.baseURL),
		body,
	)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("send create user request to data service: %w", err)
	}
	defer resp.Body.Close()

	if err := checkResponseStatus(resp, http.StatusCreated); err != nil {
		log.Default().Printf("data client: %v", err)
		return nil, err
	}

	var user model.User
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, fmt.Errorf("decode create user response: %w", err)
	}

	return &user, nil
}

func (c *Client) UpdateUser(ctx context.Context, uid string, inReq *InternalUpdateUserRequest) (*model.User, error) {
	body, err := marshalRequestBody(inReq)
	if err != nil {
		return nil, err
	}

	req, err := c.newRequest(
		ctx,
		http.MethodPatch,
		fmt.Sprintf("%s/users/%s", c.baseURL, uid),
		body,
	)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("send update user request to data service: %w", err)
	}
	defer resp.Body.Close()

	if err := checkResponseStatus(resp, http.StatusOK); err != nil {
		log.Default().Printf("data client: %v", err)
		return nil, err
	}

	var user model.User
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, fmt.Errorf("decode update user response: %w", err)
	}

	return &user, nil
}

func (c *Client) DeleteUser(ctx context.Context, uid string) (*model.User, error) {
	req, err := c.newRequest(
		ctx,
		http.MethodDelete,
		fmt.Sprintf("%s/users/%s", c.baseURL, uid),
		nil,
	)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("send delete user request to data service: %w", err)
	}
	defer resp.Body.Close()

	if err := checkResponseStatus(resp, http.StatusOK); err != nil {
		log.Default().Printf("data client: %v", err)
		return nil, err
	}

	var user model.User
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, fmt.Errorf("decode delete user response: %w", err)
	}

	return &user, nil
}

// ========================== DOCUMENTS ==========================

func (c *Client) GetAllDocuments(ctx context.Context, reqID string) ([]model.Document, error) {
	req, err := c.newRequest(
		ctx,
		http.MethodGet,
		fmt.Sprintf("%s/documents?requester_id=%s", c.baseURL, reqID),
		nil,
	)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("send get all documents request to data service: %w", err)
	}
	defer resp.Body.Close()

	if err := checkResponseStatus(resp, http.StatusOK); err != nil {
		log.Default().Printf("data client: %v", err)
		return nil, err
	}

	var documents []model.Document
	if err := json.NewDecoder(resp.Body).Decode(&documents); err != nil {
		return nil, fmt.Errorf("decode get all documents response: %w", err)
	}

	return documents, nil
}

func (c *Client) GetDocumentsByAuthor(ctx context.Context, reqID string, authID string) ([]model.Document, error) {
	req, err := c.newRequest(
		ctx,
		http.MethodGet,
		fmt.Sprintf("%s/documents?requester_id=%s&author_id=%s", c.baseURL, reqID, authID),
		nil,
	)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("send get documents by user request to data service: %w", err)
	}
	defer resp.Body.Close()

	if err := checkResponseStatus(resp, http.StatusOK); err != nil {
		log.Default().Printf("data client: %v", err)
		return nil, err
	}

	var documents []model.Document
	if err := json.NewDecoder(resp.Body).Decode(&documents); err != nil {
		return nil, fmt.Errorf("decode get documents by user response: %w", err)
	}

	return documents, nil
}

func (c *Client) GetDocumentByID(ctx context.Context, docID int, reqID string) (*model.Document, error) {
	req, err := c.newRequest(
		ctx,
		http.MethodGet,
		fmt.Sprintf("%s/documents/%d?requester_id=%s", c.baseURL, docID, reqID),
		nil,
	)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("send get document request to data service: %w", err)
	}
	defer resp.Body.Close()

	if err := checkResponseStatus(resp, http.StatusOK); err != nil {
		log.Default().Printf("data client: %v", err)
		return nil, err
	}

	var document model.Document
	if err := json.NewDecoder(resp.Body).Decode(&document); err != nil {
		return nil, fmt.Errorf("decode get document response: %w", err)
	}

	return &document, nil
}

func (c *Client) CreateDocument(ctx context.Context, inReq *InternalCreateDocumentRequest) (*model.Document, error) {
	body, err := marshalRequestBody(inReq)
	if err != nil {
		return nil, err
	}

	req, err := c.newRequest(
		ctx,
		http.MethodPost,
		fmt.Sprintf("%s/documents?requester_id=%s", c.baseURL, inReq.OwnerID),
		body,
	)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("send create document request to data service: %w", err)
	}
	defer resp.Body.Close()

	if err := checkResponseStatus(resp, http.StatusCreated); err != nil {
		log.Default().Printf("data client: %v", err)
		return nil, err
	}

	var document model.Document
	if err := json.NewDecoder(resp.Body).Decode(&document); err != nil {
		return nil, fmt.Errorf("decode create document response: %w", err)
	}

	return &document, nil
}

func (c *Client) UpdateDocument(ctx context.Context, docID int, reqID string, inReq *UpdateDocumentRequest) (*model.Document, error) {
	body, err := marshalRequestBody(inReq)
	if err != nil {
		return nil, err
	}

	req, err := c.newRequest(
		ctx,
		http.MethodPatch,
		fmt.Sprintf("%s/documents/%d?requester_id=%s", c.baseURL, docID, reqID),
		body,
	)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("send update document request to data service: %w", err)
	}
	defer resp.Body.Close()

	if err := checkResponseStatus(resp, http.StatusOK); err != nil {
		log.Default().Printf("data client: %v", err)
		return nil, err
	}

	var document model.Document
	if err := json.NewDecoder(resp.Body).Decode(&document); err != nil {
		return nil, fmt.Errorf("decode update document response: %w", err)
	}

	return &document, nil
}

func (c *Client) DeleteDocument(ctx context.Context, docID int, reqID string) (*model.Document, error) {
	req, err := c.newRequest(
		ctx,
		http.MethodDelete,
		fmt.Sprintf("%s/documents/%d?requester_id=%s", c.baseURL, docID, reqID),
		nil,
	)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("send delete document request to data service: %w", err)
	}
	defer resp.Body.Close()

	if err := checkResponseStatus(resp, http.StatusOK); err != nil {
		log.Default().Printf("data client: %v", err)
		return nil, err
	}

	var document model.Document
	if err := json.NewDecoder(resp.Body).Decode(&document); err != nil {
		return nil, fmt.Errorf("decode delete document response: %w", err)
	}

	return &document, nil
}
