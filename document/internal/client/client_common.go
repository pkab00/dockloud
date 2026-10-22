package client

import (
	"net/http"
	"os"
	"time"

	"github.com/pkab00/dockloud/shared/serviceauth"
)

type Client struct {
	BaseURL    string
	HTTPClient *http.Client
}

func NewClient(baseURL string) (*Client, error) {
	apiClaims := serviceauth.Claims{
		Subject:  "dockloud-document",
		Issuer:   "dockloud-document",
		Audience: []string{"dockloud-data", "dockloud-link"},
	}

	token, err := serviceauth.GenerateToken(
		os.Getenv("DOCUMENT_SERVICE_KEY"),
		apiClaims,
	)
	if err != nil {
		return nil, err
	}

	transport := &serviceauth.AuthTransport{
		Base:   http.DefaultTransport,
		Token:  token,
		Claims: apiClaims,
	}

	return &Client{
		BaseURL: baseURL,
		HTTPClient: &http.Client{
			Transport: transport,
			Timeout:   5 * time.Second,
		},
	}, nil
}
