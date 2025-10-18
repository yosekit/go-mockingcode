package project

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type ProjectClient struct {
	baseURL string
	client  *http.Client
}

func NewProjectClient(baseURL string) *ProjectClient {
	return &ProjectClient{
		baseURL: baseURL,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

type ValidateAPIKeyResponse struct {
	Project *ProjectInfo `json:"project"`
	Valid   bool         `json:"valid"`
}

type ProjectInfo struct {
	ID      int64  `json:"id"`
	UserID  int64  `json:"user_id"`
	Name    string `json:"name"`
	APIKey  string `json:"api_key"`
	BaseURL string `json:"base_url"`
}

func (c *ProjectClient) ValidateAPIKey(apiKey string) (*ProjectInfo, error) {
	req, err := http.NewRequest("GET", c.baseURL+"/api-keys/"+apiKey+"/validate", nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("project service returned status: %d", resp.StatusCode)
	}

	var result ValidateAPIKeyResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	if !result.Valid {
		return nil, fmt.Errorf("invalid API key")
	}

	return result.Project, nil
}
