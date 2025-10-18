package client

import (
	"fmt"
	"io"
	"net/http"
)

type ProjectClient struct {
	baseURL string
	client  *http.Client
}

func NewProjectClient(baseURL string) *ProjectClient {
	return &ProjectClient{
		baseURL: baseURL,
		client:  &http.Client{},
	}
}

// ProxyRequest forwards the request to the project service
func (c *ProjectClient) ProxyRequest(r *http.Request, path string) (*http.Response, error) {
	// Create new request to project service
	url := c.baseURL + path
	fmt.Printf("[ProjectClient] Proxying %s to %s\n", r.Method, url)
	
	proxyReq, err := http.NewRequest(r.Method, url, r.Body)
	if err != nil {
		fmt.Printf("[ProjectClient] Error creating request: %v\n", err)
		return nil, fmt.Errorf("failed to create proxy request: %w", err)
	}

	// Copy headers
	for key, values := range r.Header {
		for _, value := range values {
			proxyReq.Header.Add(key, value)
		}
	}
	fmt.Printf("[ProjectClient] Copied %d headers\n", len(r.Header))

	// Execute request
	resp, err := c.client.Do(proxyReq)
	if err != nil {
		fmt.Printf("[ProjectClient] Error executing request: %v\n", err)
		return nil, fmt.Errorf("failed to execute proxy request: %w", err)
	}

	fmt.Printf("[ProjectClient] Got response: %d\n", resp.StatusCode)
	return resp, nil
}

// CopyResponse copies response from service to client
func CopyResponse(w http.ResponseWriter, resp *http.Response) error {
	// Copy headers
	for key, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}

	// Copy status code
	w.WriteHeader(resp.StatusCode)

	// Copy body
	_, err := io.Copy(w, resp.Body)
	return err
}

