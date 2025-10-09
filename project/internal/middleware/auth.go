package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/go-mockingcode/project/internal/pkg/auth"
	authctx "github.com/go-mockingcode/project/internal/pkg/context"
)

var (
	PublicPaths = []string{
		"/health",
		"/swagger/",
	}
)

func AuthMiddleware(authClient *auth.AuthClient) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Пропускаем публичные эндпоинты
			for _, path := range PublicPaths {
				if r.URL.Path == path || strings.HasPrefix(r.URL.Path, path) {
					next.ServeHTTP(w, r)
					return
				}
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
			user, err := authClient.ValidateToken(token)
			if err != nil {
				writeErrorJson(w, http.StatusUnauthorized, "Invalid token")
				return
			}

			// Добавляем user_id в контекст
			ctx := context.WithValue(r.Context(), authctx.UserIDKey, user.UserID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func writeErrorJson(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}
