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

type InternalCreateRequest struct {
	Username     string `json:"username"`
	Email        string `json:"email"`
	PasswordHash string `json:"password_hash"`
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

func (c *Client) CreateUser(ctx context.Context, inReq *InternalCreateRequest) (*model.User, error) {
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
