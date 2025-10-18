package client

import (
	"fmt"
	"net/http"
)

type DataClient struct {
	baseURL string
	client  *http.Client
}

func NewDataClient(baseURL string) *DataClient {
	return &DataClient{
		baseURL: baseURL,
		client:  &http.Client{},
	}
}

// ProxyRequest forwards the request to the data service
func (c *DataClient) ProxyRequest(r *http.Request, path string) (*http.Response, error) {
	// Create new request to data service
	url := c.baseURL + path

	proxyReq, err := http.NewRequest(r.Method, url, r.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to create proxy request: %w", err)
	}

	// Copy headers
	for key, values := range r.Header {
		for _, value := range values {
			proxyReq.Header.Add(key, value)
		}
	}

	// Execute request
	resp, err := c.client.Do(proxyReq)
	if err != nil {
		return nil, fmt.Errorf("failed to execute proxy request: %w", err)
	}

	return resp, nil
}

