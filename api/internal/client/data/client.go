package data

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/pkab00/dockloud/api/internal/model"
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

type CreateDocumentRequest struct {
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

func NewClient(baseURL string) *Client {
	return &Client{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

// ========================== USERS ==========================

func (c *Client) GetUser(ctx context.Context, uid string) (*model.User, error) {
	var u *model.User
	var err error

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		fmt.Sprintf("%s/users/%s", c.baseURL, uid),
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("client error: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("client error: %w", err)
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&u); err != nil {
		return nil, fmt.Errorf("client error: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("client returned status %d", resp.StatusCode)
	}
	return u, nil
}

func (c *Client) GetAllUsers(ctx context.Context) ([]model.User, error) {
	var u []model.User
	var err error

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		fmt.Sprintf("%s/users", c.baseURL),
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("client error: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("client error: %w", err)
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&u); err != nil {
		return nil, fmt.Errorf("client error: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("client returned status %d", resp.StatusCode)
	}
	return u, nil
}

func (c *Client) CreateUser(ctx context.Context, inReq *InternalCreateUserRequest) (*model.User, error) {
	var u *model.User
	var err error

	body, err := json.Marshal(inReq)
	if err != nil {
		return nil, fmt.Errorf("client error: %w", err)
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		fmt.Sprintf("%s/users", c.baseURL),
		bytes.NewReader(body),
	)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("client error: %w", err)
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&u); err != nil {
		return nil, fmt.Errorf("client error: %w", err)
	}
	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("client returned status %d", resp.StatusCode)
	}
	return u, nil
}

func (c *Client) UpdateUser(ctx context.Context, uid string, inReq *InternalUpdateUserRequest) (*model.User, error) {
	var u *model.User
	var err error

	body, err := json.Marshal(inReq)
	if err != nil {
		return nil, fmt.Errorf("client error: %w", err)
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPatch,
		fmt.Sprintf("%s/users/%s", c.baseURL, uid),
		bytes.NewReader(body),
	)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("client error: %w", err)
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&u); err != nil {
		return nil, fmt.Errorf("client error: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("client returned status %d", resp.StatusCode)
	}
	return u, nil
}

func (c *Client) DeleteUser(ctx context.Context, uid string) (*model.User, error) {
	var u *model.User
	var err error

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodDelete,
		fmt.Sprintf("%s/users/%s", c.baseURL, uid),
		nil,
	)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("client error: %w", err)
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&u); err != nil {
		return nil, fmt.Errorf("client error: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("client returned status %d", resp.StatusCode)
	}
	return u, nil
}

// ========================== DOCUMENTS ==========================

func (c *Client) GetAllDocuments(ctx context.Context) ([]model.Document, error) {
	var d []model.Document
	var err error

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		fmt.Sprintf("%s/documents", c.baseURL),
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("client error: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("client error: %w", err)
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&d); err != nil {
		return nil, fmt.Errorf("client error: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("client returned status %d", resp.StatusCode)
	}
	return d, nil
}

func (c *Client) GetDocumentsByUID(ctx context.Context, uid string) ([]model.Document, error) {
	var d []model.Document
	var err error

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		fmt.Sprintf("%s/documents/user/%s", c.baseURL, uid),
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("client error: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("client error: %w", err)
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&d); err != nil {
		return nil, fmt.Errorf("client error: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("client returned status %d", resp.StatusCode)
	}
	return d, nil
}

func (c *Client) GetDocumentByID(ctx context.Context, id int) (*model.Document, error) {
	var d *model.Document
	var err error

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		fmt.Sprintf("%s/documents/%d", c.baseURL, id),
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("client error: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("client error: %w", err)
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&d); err != nil {
		return nil, fmt.Errorf("client error: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("client returned status %d", resp.StatusCode)
	}
	return d, nil
}

func (c *Client) CreateDocument(ctx context.Context, req *CreateDocumentRequest) (*model.Document, error) {
	var d *model.Document
	var err error

	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("client error: %w", err)
	}

	dataReq, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		fmt.Sprintf("%s/documents", c.baseURL),
		bytes.NewReader(body),
	)

	resp, err := c.httpClient.Do(dataReq)
	if err != nil {
		return nil, fmt.Errorf("client error: %w", err)
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&d); err != nil {
		return nil, fmt.Errorf("client error: %w", err)
	}
	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("client returned status %d", resp.StatusCode)
	}
	return d, nil
}

func (c *Client) UpdateDocument(ctx context.Context, id int, req *CreateDocumentRequest) (*model.Document, error) {
	var d *model.Document
	var err error

	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("client error: %w", err)
	}

	dataReq, err := http.NewRequestWithContext(
		ctx,
		http.MethodPatch,
		fmt.Sprintf("%s/documents/%d", c.baseURL, id),
		bytes.NewReader(body),
	)

	resp, err := c.httpClient.Do(dataReq)
	if err != nil {
		return nil, fmt.Errorf("client error: %w", err)
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&d); err != nil {
		return nil, fmt.Errorf("client error: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("client returned status %d", resp.StatusCode)
	}
	return d, nil
}

func (c *Client) DeleteDocument(ctx context.Context, id int) (*model.Document, error) {
	var d *model.Document
	var err error

	dataReq, err := http.NewRequestWithContext(
		ctx,
		http.MethodDelete,
		fmt.Sprintf("%s/documents/%d", c.baseURL, id),
		nil,
	)

	resp, err := c.httpClient.Do(dataReq)
	if err != nil {
		return nil, fmt.Errorf("client error: %w", err)
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&d); err != nil {
		return nil, fmt.Errorf("client error: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("client returned status %d", resp.StatusCode)
	}
	return d, nil
}
