package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-mockingcode/auth/internal/model"
	"github.com/go-mockingcode/auth/internal/service"
)

type AuthHandler struct {
    authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
    return &AuthHandler{authService: authService}
}

// Register - create user + auto login
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


func writeErrorJson(w http.ResponseWriter, statusCode int, message string) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(statusCode)
    json.NewEncoder(w).Encode(map[string]string{"error": message})
}

func writeSuccessJson(w http.ResponseWriter, statusCode int, data interface{}) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(statusCode)
    json.NewEncoder(w).Encode(data)
}