package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/go-mockingcode/auth/internal/model"
	"github.com/go-mockingcode/auth/internal/service"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// Register godoc
// @Summary Register a new user
// @Description Create a new user account and automatically login
// @Tags auth
// @Accept json
// @Produce json
// @Param request body model.RegisterRequest true "Registration data"
// @Success 201 {object} model.AuthResponse
// @Failure 400 {object} map[string]string
// @Failure 409 {object} map[string]string
// @Router /register [post]
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeErrorJson(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var req model.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErrorJson(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// 1. Create User
	if err := h.authService.Register(&req); err != nil {
		writeErrorJson(w, http.StatusBadRequest, err.Error())
		return
	}

	// 2. Auto Login
	response, err := h.authService.AutoLogin(req.Email)
	if err != nil {
		writeErrorJson(w, http.StatusInternalServerError, "Registration successful but auto-login failed")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// Login godoc
// @Summary Login user
// @Description Authenticate user and return JWT tokens
// @Tags auth
// @Accept json
// @Produce json
// @Param request body model.LoginRequest true "Login credentials"
// @Success 200 {object} model.AuthResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /login [post]
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeErrorJson(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var req model.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErrorJson(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	response, err := h.authService.Login(&req)
	if err != nil {
		writeErrorJson(w, http.StatusUnauthorized, err.Error())
		return
	}

	writeSuccessJson(w, http.StatusOK, response)
}

// Refresh godoc
// @Summary Refresh JWT tokens
// @Description Get new access and refresh tokens using refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body map[string]string true "Refresh token"
// @Success 200 {object} model.AuthResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /refresh [post]
func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeErrorJson(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var req struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErrorJson(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	response, err := h.authService.RefreshTokens(req.RefreshToken)
	if err != nil {
		writeErrorJson(w, http.StatusUnauthorized, err.Error())
		return
	}

	writeSuccessJson(w, http.StatusOK, response)
}

// Validate godoc
// @Summary Validate JWT token
// @Description Validate JWT token and return user information
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]string
// @Router /validate [get]
func (h *AuthHandler) Validate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeErrorJson(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		writeErrorJson(w, http.StatusUnauthorized, "Authorization header required")
		return
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		writeErrorJson(w, http.StatusUnauthorized, "Invalid authorization format")
		return
	}

	token := parts[1]
	user, err := h.authService.ValidateAccessToken(token)
	if err != nil {
		writeErrorJson(w, http.StatusUnauthorized, "Invalid token")
		return
	}

	writeSuccessJson(w, http.StatusOK, map[string]any{
		"user_id": user.ID,
		"email":   user.Email,
	})
}

// Logout godoc
// @Summary Logout user
// @Description Invalidate refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body map[string]string true "Refresh token"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /logout [post]
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeErrorJson(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var req struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErrorJson(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := h.authService.Logout(req.RefreshToken); err != nil {
		writeErrorJson(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeSuccessJson(w, http.StatusOK, map[string]string{"status": "logged out"})
}

// TODO refactor to common utils
func writeErrorJson(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

func writeSuccessJson(w http.ResponseWriter, statusCode int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}
