package auth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type ValidateTokenResponse struct {
	UserID int64  `json:"user_id"`
	Email  string `json:"email"`
}

type AuthClient struct {
	baseURL string
	client  *http.Client
}

func NewAuthClient(baseURL string) *AuthClient {
	return &AuthClient{
		baseURL: baseURL,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *AuthClient) ValidateToken(token string) (*ValidateTokenResponse, error) {
	req, err := http.NewRequest("GET", c.baseURL+"/validate", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("auth service returned status: %d", resp.StatusCode)
	}

	var result ValidateTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}
