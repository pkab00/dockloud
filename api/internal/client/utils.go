package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	apperror "github.com/pkab00/dockloud/api/internal/error"
	"github.com/pkab00/dockloud/shared/contract"
)

func CheckResponseStatus(resp *http.Response, expectedStatus int) error {
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

func MarshalRequestBody(v any) (*bytes.Reader, error) {
	body, err := json.Marshal(v)
	if err != nil {
		return nil, fmt.Errorf("marshal data service request: %w", err)
	}
	return bytes.NewReader(body), nil
}

func (c *Client) NewRequest(
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
