package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	projectctx "github.com/go-mockingcode/data/internal/pkg/context"
	"github.com/go-mockingcode/data/internal/pkg/project"
)

var (
	PublicPaths = []string{
		"/health",
		"/swagger/",
	}
)

func AuthMiddleware(projectClient *project.ProjectClient) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Пропускаем публичные эндпоинты
			for _, path := range PublicPaths {
				if r.URL.Path == path || strings.HasPrefix(r.URL.Path, path) {
					next.ServeHTTP(w, r)
					return
				}
			}

			// Извлекаем API Key из пути URL
			// Формат: /{api_key}/{collection}/...
			pathParts := strings.Split(r.URL.Path, "/")
			if len(pathParts) < 2 {
				writeErrorJson(w, http.StatusUnauthorized, "Invalid URL format")
				return
			}

			apiKey := pathParts[1]
			if apiKey == "" {
				writeErrorJson(w, http.StatusUnauthorized, "API Key is required")
				return
			}

			// Валидируем API Key через Project Service
			projectInfo, err := projectClient.ValidateAPIKey(apiKey)
			if err != nil {
				writeErrorJson(w, http.StatusUnauthorized, "Invalid API Key")
				return
			}

			// Добавляем информацию о проекте в контекст
			ctx := context.WithValue(r.Context(), projectctx.ProjectKey, projectInfo)

			// Убираем API Key из пути для следующих handlers
			newPath := "/" + strings.Join(pathParts[2:], "/")
			r.URL.Path = newPath

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// TODO refactor fot common utils
func writeErrorJson(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}
