package middleware

import (
	"log"
	"net/http"
	"strings"

	"github.com/go-mockingcode/gateway/internal/config"
)

// CORSMiddleware handles CORS headers
func CORSMiddleware(cfg *config.Config) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Printf("[CORSMiddleware] Request: %s %s", r.Method, r.URL.Path)
			
			// Set CORS headers
			origin := r.Header.Get("Origin")
			if origin != "" && isOriginAllowed(origin, cfg.CORSAllowedOrigins) {
				w.Header().Set("Access-Control-Allow-Origin", origin)
			} else if len(cfg.CORSAllowedOrigins) == 1 && cfg.CORSAllowedOrigins[0] == "*" {
				w.Header().Set("Access-Control-Allow-Origin", "*")
			}

			w.Header().Set("Access-Control-Allow-Methods", strings.Join(cfg.CORSAllowedMethods, ", "))
			w.Header().Set("Access-Control-Allow-Headers", strings.Join(cfg.CORSAllowedHeaders, ", "))
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Max-Age", "3600")

			// Handle preflight requests
			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func isOriginAllowed(origin string, allowedOrigins []string) bool {
	for _, allowed := range allowedOrigins {
		if allowed == "*" || allowed == origin {
			return true
		}
	}
	return false
}

