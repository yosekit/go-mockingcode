package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type AuthClient struct {
	baseURL string
	client  *http.Client
}

func NewAuthClient(baseURL string) *AuthClient {
	return &AuthClient{
		baseURL: baseURL,
		client:  &http.Client{},
	}
}

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type ValidateResponse struct {
	Valid  bool   `json:"valid"`
	UserID string `json:"user_id"`
}

func (c *AuthClient) Register(email, password string) (*AuthResponse, error) {
	reqBody := RegisterRequest{Email: email, Password: password}
	body, _ := json.Marshal(reqBody)

	resp, err := c.client.Post(c.baseURL+"/register", "application/json", bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("registration failed: %s", string(bodyBytes))
	}

	var authResp AuthResponse
	if err := json.NewDecoder(resp.Body).Decode(&authResp); err != nil {
		return nil, err
	}

	return &authResp, nil
}

func (c *AuthClient) Login(email, password string) (*AuthResponse, error) {
	reqBody := LoginRequest{Email: email, Password: password}
	body, _ := json.Marshal(reqBody)

	resp, err := c.client.Post(c.baseURL+"/login", "application/json", bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("login failed: %s", string(bodyBytes))
	}

	var authResp AuthResponse
	if err := json.NewDecoder(resp.Body).Decode(&authResp); err != nil {
		return nil, err
	}

	return &authResp, nil
}

func (c *AuthClient) ValidateToken(token string) (*ValidateResponse, error) {
	req, err := http.NewRequest("GET", c.baseURL+"/validate", nil)
	if err != nil {
		fmt.Printf("[AuthClient] Error creating request: %v\n", err)
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := c.client.Do(req)
	if err != nil {
		fmt.Printf("[AuthClient] Error executing request: %v\n", err)
		return nil, err
	}
	defer resp.Body.Close()

	fmt.Printf("[AuthClient] Response status: %d\n", resp.StatusCode)

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("[AuthClient] Non-200 status, returning Valid=false\n")
		return &ValidateResponse{Valid: false}, nil
	}

	// Auth service returns {user_id, email}, not {valid, user_id}
	// user_id can be either number or string depending on auth implementation
	var authResponse struct {
		UserID interface{} `json:"user_id"`
		Email  string      `json:"email"`
	}
	
	if err := json.NewDecoder(resp.Body).Decode(&authResponse); err != nil {
		fmt.Printf("[AuthClient] Error decoding response: %v\n", err)
		return nil, err
	}

	fmt.Printf("[AuthClient] Decoded response: %+v\n", authResponse)

	// Convert user_id to string
	var userID string
	switch v := authResponse.UserID.(type) {
	case float64:
		userID = fmt.Sprintf("%.0f", v)
	case string:
		userID = v
	default:
		userID = fmt.Sprintf("%v", v)
	}

	fmt.Printf("[AuthClient] Converted userID: %s, Valid: true\n", userID)

	// Convert to ValidateResponse format
	return &ValidateResponse{
		Valid:  true,
		UserID: userID,
	}, nil
}

func (c *AuthClient) Refresh(refreshToken string) (*AuthResponse, error) {
	reqBody := map[string]string{"refresh_token": refreshToken}
	body, _ := json.Marshal(reqBody)

	resp, err := c.client.Post(c.baseURL+"/refresh", "application/json", bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("token refresh failed: %s", string(bodyBytes))
	}

	var authResp AuthResponse
	if err := json.NewDecoder(resp.Body).Decode(&authResp); err != nil {
		return nil, err
	}

	return &authResp, nil
}

