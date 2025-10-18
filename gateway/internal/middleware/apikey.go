package middleware

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/go-mockingcode/gateway/internal/client"
)

type projectContextKey string

const ProjectInfoKey projectContextKey = "project_info"

// APIKeyMiddleware validates API key for public data access
func APIKeyMiddleware(projectClient *client.ProjectGRPCClient) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			slog.Debug("apikey middleware", 
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
			)

			// Extract API key from path: /{api_key}/{collection}
			pathParts := strings.Split(strings.TrimPrefix(r.URL.Path, "/"), "/")
			
			if len(pathParts) < 2 {
				writeError(w, http.StatusBadRequest, "Invalid API endpoint format. Use: /{api_key}/{collection}")
				return
			}

			apiKey := pathParts[0]
			if apiKey == "" {
				writeError(w, http.StatusBadRequest, "API key is required")
				return
			}

			// Mask API key for logging (show first 8 chars or less)
			maskedKey := apiKey
			if len(apiKey) > 8 {
				maskedKey = apiKey[:8] + "..."
			}
			slog.Debug("validating API key", slog.String("api_key", maskedKey))

			// Validate API key via gRPC
			projectInfo, err := projectClient.ValidateAPIKey(apiKey)
			if err != nil {
				slog.Error("failed to validate API key", slog.String("error", err.Error()))
				writeError(w, http.StatusUnauthorized, "Failed to validate API key")
				return
			}

			if !projectInfo.Valid {
				slog.Warn("invalid API key", slog.String("api_key", maskedKey))
				writeError(w, http.StatusUnauthorized, "Invalid API key")
				return
			}

			slog.Debug("API key valid",
				slog.Int64("project_id", projectInfo.ProjectID),
				slog.Int64("user_id", projectInfo.UserID),
				slog.String("project_name", projectInfo.ProjectName),
			)

			// Extract collection name from path: /{api_key}/{collection}
			collectionName := ""
			if len(pathParts) >= 2 {
				collectionName = pathParts[1]
			}

			// Check if collection schema exists (optional validation)
			if collectionName != "" {
				schema, err := projectClient.GetCollectionSchema(projectInfo.ProjectID, collectionName)
				if err != nil {
					slog.Error("failed to check collection", slog.String("error", err.Error()))
					writeError(w, http.StatusInternalServerError, "Failed to validate collection")
					return
				}

				// If schema exists but inactive - deny access
				if schema.Found && !schema.IsActive {
					slog.Warn("collection is inactive",
						slog.String("collection", collectionName),
						slog.Int64("project_id", projectInfo.ProjectID),
					)
					writeError(w, http.StatusForbidden, "Collection is not active")
					return
				}

				// Log schema status
				if schema.Found {
					slog.Debug("collection schema found (validation enabled)",
						slog.String("collection", collectionName),
						slog.Int64("collection_id", schema.CollectionID),
					)
				} else {
					slog.Debug("collection schema not found (schema-less mode)",
						slog.String("collection", collectionName),
					)
				}
			}

			// Add project info to context
			ctx := context.WithValue(r.Context(), ProjectInfoKey, projectInfo)
			
			// Add headers for Data Service
			r.Header.Set("X-Project-ID", fmt.Sprintf("%d", projectInfo.ProjectID))
			r.Header.Set("X-User-ID", fmt.Sprintf("%d", projectInfo.UserID))
			
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

