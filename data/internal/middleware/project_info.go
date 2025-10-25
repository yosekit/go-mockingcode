package middleware

import (
	stdcontext "context"
	"log/slog"
	"net/http"
	"strconv"

	appcontext "github.com/go-mockingcode/data/internal/pkg/context"
	"github.com/go-mockingcode/data/internal/pkg/project"
)

type contextKey string

const ProjectInfoKey contextKey = "project_info"

// ProjectInfoMiddleware extracts user ID and project ID from headers (set by API Gateway)
func ProjectInfoMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Skip for public endpoints
			if isPublicEndpoint(r.URL.Path) {
				next.ServeHTTP(w, r)
				return
			}

			// Extract headers set by Gateway
			userIDStr := r.Header.Get("X-User-ID")
			projectIDStr := r.Header.Get("X-Project-ID")
			
			if userIDStr != "" && projectIDStr != "" {
				// Parse IDs
				userID, err := strconv.ParseInt(userIDStr, 10, 64)
				if err != nil {
					slog.Error("failed to parse user ID", slog.String("user_id", userIDStr), slog.String("error", err.Error()))
					http.Error(w, "Invalid user ID", http.StatusBadRequest)
					return
				}

				projectID, err := strconv.ParseInt(projectIDStr, 10, 64)
				if err != nil {
					slog.Error("failed to parse project ID", slog.String("project_id", projectIDStr), slog.String("error", err.Error()))
					http.Error(w, "Invalid project ID", http.StatusBadRequest)
					return
				}

				// Create ProjectInfo and add to context
				projectInfo := &project.ProjectInfo{
					ID:     projectID,
					UserID: userID,
				}

				ctx := stdcontext.WithValue(r.Context(), appcontext.ProjectKey, projectInfo)

				slog.Debug("request from gateway", 
					slog.Int64("user_id", userID),
					slog.Int64("project_id", projectID),
				)

				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}

			// No headers - this shouldn't happen if accessing through Gateway
			slog.Warn("no X-User-ID or X-Project-ID headers", slog.String("path", r.URL.Path))
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
		})
	}
}

func isPublicEndpoint(path string) bool {
	publicPaths := []string{
		"/health",
		"/swagger",
		"/generate",
	}

	for _, publicPath := range publicPaths {
		if path == publicPath || len(path) > len(publicPath) && path[:len(publicPath)+1] == publicPath+"/" {
			return true
		}
	}

	return false
}

