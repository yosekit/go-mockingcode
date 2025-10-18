package middleware

import (
	"context"
	"log/slog"
	"net/http"
)

type contextKey string

const UserIDKey contextKey = "user_id"

// UserIDMiddleware extracts user ID from X-User-ID header (set by API Gateway)
// and adds it to the request context
func UserIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := r.Header.Get("X-User-ID")
		
		if userID != "" {
			slog.Debug("extracted user_id from header", slog.String("user_id", userID))
			ctx := context.WithValue(r.Context(), UserIDKey, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		} else {
			// No user ID header - for public endpoints like /health
			slog.Debug("no X-User-ID header", slog.String("path", r.URL.Path))
			next.ServeHTTP(w, r)
		}
	})
}

// GetUserID extracts user ID from context
func GetUserID(ctx context.Context) (string, bool) {
	userID, ok := ctx.Value(UserIDKey).(string)
	return userID, ok
}

