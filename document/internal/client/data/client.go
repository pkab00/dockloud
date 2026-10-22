package data

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/pkab00/dockloud/document/internal/client"
	"github.com/pkab00/dockloud/document/internal/model"
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
	*client.Client
}

func NewClient(baseURL string) (*Client, error) {
	cl, err := client.NewClient(baseURL)
	if err != nil {
		return nil, err
	}
	return &Client{cl}, nil
}

// ========================== USERS ==========================

func (c *Client) FindUser(ctx context.Context, req *FindUserRequest) (*model.UserFindResult, error) {
	body, err := client.MarshalRequestBody(req)
	if err != nil {
		return nil, err
	}

	dataReq, err := c.NewRequest(
		ctx,
		http.MethodGet,
		fmt.Sprintf("%s/users/find", c.BaseURL),
		body,
	)
	if err != nil {
		return nil, err
	}

	resp, err := c.HTTPClient.Do(dataReq)
	if err != nil {
		return nil, fmt.Errorf("send find user request to data service: %w", err)
	}
	defer resp.Body.Close()

	if err := client.CheckResponseStatus(resp, http.StatusOK); err != nil {
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
	body, err := client.MarshalRequestBody(req)
	if err != nil {
		return nil, err
	}

	dataReq, err := c.NewRequest(
		ctx,
		http.MethodGet,
		fmt.Sprintf("%s/users/auth", c.BaseURL),
		body,
	)
	if err != nil {
		return nil, err
	}

	resp, err := c.HTTPClient.Do(dataReq)
	if err != nil {
		return nil, fmt.Errorf("send authenticate user request to data service: %w", err)
	}
	defer resp.Body.Close()

	if err := client.CheckResponseStatus(resp, http.StatusOK); err != nil {
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
	body, err := client.MarshalRequestBody(inReq)
	if err != nil {
		return nil, err
	}

	req, err := c.NewRequest(
		ctx,
		http.MethodPost,
		fmt.Sprintf("%s/users", c.BaseURL),
		body,
	)
	if err != nil {
		return nil, err
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("send create user request to data service: %w", err)
	}
	defer resp.Body.Close()

	if err := client.CheckResponseStatus(resp, http.StatusCreated); err != nil {
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
	body, err := client.MarshalRequestBody(inReq)
	if err != nil {
		return nil, err
	}

	req, err := c.NewRequest(
		ctx,
		http.MethodPatch,
		fmt.Sprintf("%s/users/%s", c.BaseURL, uid),
		body,
	)
	if err != nil {
		return nil, err
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("send update user request to data service: %w", err)
	}
	defer resp.Body.Close()

	if err := client.CheckResponseStatus(resp, http.StatusOK); err != nil {
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
	req, err := c.NewRequest(
		ctx,
		http.MethodDelete,
		fmt.Sprintf("%s/users/%s", c.BaseURL, uid),
		nil,
	)
	if err != nil {
		return nil, err
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("send delete user request to data service: %w", err)
	}
	defer resp.Body.Close()

	if err := client.CheckResponseStatus(resp, http.StatusOK); err != nil {
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
	req, err := c.NewRequest(
		ctx,
		http.MethodGet,
		fmt.Sprintf("%s/documents?requester_id=%s", c.BaseURL, reqID),
		nil,
	)
	if err != nil {
		return nil, err
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("send get all documents request to data service: %w", err)
	}
	defer resp.Body.Close()

	if err := client.CheckResponseStatus(resp, http.StatusOK); err != nil {
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
	req, err := c.NewRequest(
		ctx,
		http.MethodGet,
		fmt.Sprintf("%s/documents?requester_id=%s&author_id=%s", c.BaseURL, reqID, authID),
		nil,
	)
	if err != nil {
		return nil, err
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("send get documents by user request to data service: %w", err)
	}
	defer resp.Body.Close()

	if err := client.CheckResponseStatus(resp, http.StatusOK); err != nil {
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
	req, err := c.NewRequest(
		ctx,
		http.MethodGet,
		fmt.Sprintf("%s/documents/%d?requester_id=%s", c.BaseURL, docID, reqID),
		nil,
	)
	if err != nil {
		return nil, err
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("send get document request to data service: %w", err)
	}
	defer resp.Body.Close()

	if err := client.CheckResponseStatus(resp, http.StatusOK); err != nil {
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
	body, err := client.MarshalRequestBody(inReq)
	if err != nil {
		return nil, err
	}

	req, err := c.NewRequest(
		ctx,
		http.MethodPost,
		fmt.Sprintf("%s/documents?requester_id=%s", c.BaseURL, inReq.OwnerID),
		body,
	)
	if err != nil {
		return nil, err
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("send create document request to data service: %w", err)
	}
	defer resp.Body.Close()

	if err := client.CheckResponseStatus(resp, http.StatusCreated); err != nil {
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
	body, err := client.MarshalRequestBody(inReq)
	if err != nil {
		return nil, err
	}

	req, err := c.NewRequest(
		ctx,
		http.MethodPatch,
		fmt.Sprintf("%s/documents/%d?requester_id=%s", c.BaseURL, docID, reqID),
		body,
	)
	if err != nil {
		return nil, err
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("send update document request to data service: %w", err)
	}
	defer resp.Body.Close()

	if err := client.CheckResponseStatus(resp, http.StatusOK); err != nil {
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
	req, err := c.NewRequest(
		ctx,
		http.MethodDelete,
		fmt.Sprintf("%s/documents/%d?requester_id=%s", c.BaseURL, docID, reqID),
		nil,
	)
	if err != nil {
		return nil, err
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("send delete document request to data service: %w", err)
	}
	defer resp.Body.Close()

	if err := client.CheckResponseStatus(resp, http.StatusOK); err != nil {
		log.Default().Printf("data client: %v", err)
		return nil, err
	}

	var document model.Document
	if err := json.NewDecoder(resp.Body).Decode(&document); err != nil {
		return nil, fmt.Errorf("decode delete document response: %w", err)
	}

	return &document, nil
}
