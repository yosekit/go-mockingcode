package main

import (
	"log"
	"net/http"
	"strings"

	"github.com/go-mockingcode/gateway/internal/client"
	"github.com/go-mockingcode/gateway/internal/config"
	"github.com/go-mockingcode/gateway/internal/handler"
	"github.com/go-mockingcode/gateway/internal/middleware"
	"github.com/joho/godotenv"
)

// @title MockingCode API Gateway
// @version 1.0
// @description API Gateway for MockingCode platform - unified entry point for all services
// @termsOfService http://mockingcode.dev/terms/

// @contact.name API Support
// @contact.url http://mockingcode.dev/support
// @contact.email support@mockingcode.dev

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description Enter the token with the `Bearer ` prefix, e.g. "Bearer abcde12345"
func main() {
	// Load .env
	if err := godotenv.Load(".env"); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Load configuration
	cfg := config.Load()

	// Initialize service clients
	authClient := client.NewAuthClient(cfg.AuthServiceURL)
	projectClient := client.NewProjectClient(cfg.ProjectServiceURL)
	dataClient := client.NewDataClient(cfg.DataServiceURL)

	// Initialize handlers
	authHandler := handler.NewAuthHandler(authClient)
	proxyHandler := handler.NewProxyHandler(projectClient, dataClient)

	// Setup routing
	mux := http.NewServeMux()

	// Health check
	mux.HandleFunc("/health", healthHandler)

	// Auth routes (public)
	mux.HandleFunc("/auth/register", authHandler.Register)
	mux.HandleFunc("/auth/login", authHandler.Login)
	mux.HandleFunc("/auth/refresh", authHandler.Refresh)

	// API routes (protected) - proxy to project service
	mux.HandleFunc("/api/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("[Router] /api/ handler called, path: %s", r.URL.Path)
		if strings.HasPrefix(r.URL.Path, "/api/projects") {
			log.Printf("[Router] Routing to HandleProjects")
			proxyHandler.HandleProjects(w, r)
		} else {
			log.Printf("[Router] No match for path %s, returning 404", r.URL.Path)
			http.NotFound(w, r)
		}
	})

	// Data routes (protected) - proxy to data service
	mux.HandleFunc("/data/", proxyHandler.HandleData)

	// Apply middleware chain
	handler := middleware.CORSMiddleware(cfg)(mux)
	handler = middleware.AuthMiddleware(authClient)(handler)

	// Start server
	log.Printf("🚀 API Gateway starting on port %s", cfg.ServerPort)
	log.Printf("📡 Auth Service: %s", cfg.AuthServiceURL)
	log.Printf("📡 Project Service: %s", cfg.ProjectServiceURL)
	log.Printf("📡 Data Service: %s", cfg.DataServiceURL)
	
	if err := http.ListenAndServe(":"+cfg.ServerPort, handler); err != nil {
		log.Fatal("Failed to start gateway:", err)
	}
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok","service":"gateway"}`))
}

