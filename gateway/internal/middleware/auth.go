package middleware

import (
	"context"
	"log/slog"
	"net/http"
	"strings"

	"github.com/go-mockingcode/gateway/internal/client"
)

type contextKey string

const UserIDKey contextKey = "user_id"

// AuthMiddleware validates JWT token for protected routes
// Accepts any client that implements TokenValidator interface (HTTP or gRPC)
func AuthMiddleware(authClient client.TokenValidator) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			slog.Debug("auth middleware",
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
			)
			
			// Skip auth for public endpoints
			if isPublicEndpoint(r.URL.Path) {
				slog.Debug("public endpoint, skipping auth", slog.String("path", r.URL.Path))
				next.ServeHTTP(w, r)
				return
			}
			
			slog.Debug("protected endpoint, checking auth")

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
				slog.Error("token validation failed", slog.String("error", err.Error()))
				writeError(w, http.StatusUnauthorized, "Failed to validate token")
				return
			}

			slog.Debug("token validated",
				slog.Bool("valid", validateResp.Valid),
				slog.String("user_id", validateResp.UserID),
			)

			if !validateResp.Valid {
				slog.Warn("invalid token, blocking request")
				writeError(w, http.StatusUnauthorized, "Invalid token")
				return
			}

			// Add user ID to request header for downstream services
			slog.Debug("adding X-User-ID header", slog.String("user_id", validateResp.UserID))
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
		"/admin/",  // Admin endpoints are public (TODO: add admin auth)
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

