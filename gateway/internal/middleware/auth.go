package middleware

import (
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/go-mockingcode/gateway/internal/client"
)

type contextKey string

const UserIDKey contextKey = "user_id"

// AuthMiddleware validates JWT token for protected routes
func AuthMiddleware(authClient *client.AuthClient) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Printf("[AuthMiddleware] Request: %s %s", r.Method, r.URL.Path)
			
			// Skip auth for public endpoints
			if isPublicEndpoint(r.URL.Path) {
				log.Printf("[AuthMiddleware] Public endpoint, skipping auth")
				next.ServeHTTP(w, r)
				return
			}
			
			log.Printf("[AuthMiddleware] Protected endpoint, checking auth")

			// Extract token from Authorization header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				writeError(w, http.StatusUnauthorized, "Missing authorization header")
				return
			}

			// Check Bearer prefix
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				writeError(w, http.StatusUnauthorized, "Invalid authorization header format")
				return
			}

			token := parts[1]

			// Validate token with auth service
			validateResp, err := authClient.ValidateToken(token)
			if err != nil {
				log.Printf("[AuthMiddleware] Error validating token: %v\n", err)
				writeError(w, http.StatusUnauthorized, "Failed to validate token")
				return
			}

			log.Printf("[AuthMiddleware] ValidateResp: Valid=%v, UserID=%s\n", validateResp.Valid, validateResp.UserID)

			if !validateResp.Valid {
				log.Printf("[AuthMiddleware] Token is not valid, blocking request\n")
				writeError(w, http.StatusUnauthorized, "Invalid token")
				return
			}

			// Add user ID to request header for downstream services
			log.Printf("[AuthMiddleware] Token valid, adding X-User-ID header: %s", validateResp.UserID)
			r.Header.Set("X-User-ID", validateResp.UserID)
			
			// Also add to context for potential use in gateway handlers
			ctx := context.WithValue(r.Context(), UserIDKey, validateResp.UserID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func isPublicEndpoint(path string) bool {
	publicPaths := []string{
		"/auth/register",
		"/auth/login",
		"/auth/refresh",
		"/health",
		"/swagger",
	}

	for _, publicPath := range publicPaths {
		if strings.HasPrefix(path, publicPath) {
			return true
		}
	}

	return false
}

func writeError(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write([]byte(`{"error":"` + message + `"}`))
}

