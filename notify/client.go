package notify

import (
	"net/http"
	"time"
)

type Client struct {
	baseURL    string
	apiVersion string
	HTTPClient *http.Client
}

// NewClient creates new Facest.io client with given API key
func NewClient(url string) *Client {
	return &Client{
		HTTPClient: &http.Client{
			Timeout: 5 * time.Minute,
		},
		baseURL:    url,
		apiVersion: "api/v1",
	}
}
