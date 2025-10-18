package middleware

import (
	"context"
	"log"
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
			log.Printf("[UserIDMiddleware] Extracted user ID from header: %s", userID)
			ctx := context.WithValue(r.Context(), UserIDKey, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		} else {
			// No user ID header - for public endpoints like /health
			log.Printf("[UserIDMiddleware] No X-User-ID header for: %s", r.URL.Path)
			next.ServeHTTP(w, r)
		}
	})
}

// GetUserID extracts user ID from context
func GetUserID(ctx context.Context) (string, bool) {
	userID, ok := ctx.Value(UserIDKey).(string)
	return userID, ok
}

