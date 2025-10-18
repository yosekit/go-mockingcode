package middleware

import (
	"log/slog"
	"net/http"
)

type contextKey string

const ProjectInfoKey contextKey = "project_info"

// ProjectInfoMiddleware extracts user ID from X-User-ID header (set by API Gateway)
// For direct API key access in the future, this middleware can be extended
func ProjectInfoMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Skip for public endpoints
			if isPublicEndpoint(r.URL.Path) {
				next.ServeHTTP(w, r)
				return
			}

			// For now, data service is accessed through Gateway
			// X-User-ID header is set by Gateway after authentication
			userID := r.Header.Get("X-User-ID")
			
			if userID != "" {
				slog.Debug("request from gateway", slog.String("user_id", userID))
				// In future, we might need to extract project ID from path/body
				// For now, just pass through - document handlers will handle project validation
				next.ServeHTTP(w, r)
				return
			}

			// No user ID - this shouldn't happen if accessing through Gateway
			slog.Warn("no X-User-ID header, might be direct access", slog.String("path", r.URL.Path))
			// TODO: Add API key validation for direct access
			next.ServeHTTP(w, r)
		})
	}
}

func isPublicEndpoint(path string) bool {
	publicPaths := []string{
		"/health",
		"/swagger",
	}

	for _, publicPath := range publicPaths {
		if path == publicPath || len(path) > len(publicPath) && path[:len(publicPath)+1] == publicPath+"/" {
			return true
		}
	}

	return false
}

