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
	"github.com/pkab00/dockloud/shared/contract"
	sharedserviceauth "github.com/pkab00/dockloud/shared/serviceauth"
)

type Client struct {
	baseURL    string
	httpClient *http.Client
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
