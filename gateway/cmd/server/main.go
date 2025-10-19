package main

import (
	"log"
	"log/slog"
	"net/http"
	"strings"

	"github.com/go-mockingcode/gateway/internal/client"
	"github.com/go-mockingcode/gateway/internal/config"
	"github.com/go-mockingcode/gateway/internal/handler"
	"github.com/go-mockingcode/gateway/internal/middleware"
	applogger "github.com/go-mockingcode/logger"
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

	// Initialize logger
	logger := applogger.FromEnv()
	slog.SetDefault(logger)

	// Initialize service clients
	// Using gRPC for auth token validation (faster, type-safe)
	authGRPCClient, err := client.NewAuthGRPCClient(cfg.AuthGRPCURL)
	if err != nil {
		log.Fatalf("Failed to connect to auth gRPC service: %v", err)
	}
	defer authGRPCClient.Close()

	// Using gRPC for API key validation (faster, type-safe)
	projectGRPCClient, err := client.NewProjectGRPCClient(cfg.ProjectGRPCURL)
	if err != nil {
		log.Fatalf("Failed to connect to project gRPC service: %v", err)
	}
	defer projectGRPCClient.Close()

	// HTTP clients for proxying full requests
	projectClient := client.NewProjectClient(cfg.ProjectServiceURL)
	dataClient := client.NewDataClient(cfg.DataServiceURL)

	// Initialize handlers
	// Auth handler still uses HTTP client for register/login/refresh (not performance critical)
	authHTTPClient := client.NewAuthClient(cfg.AuthServiceURL)
	authHandler := handler.NewAuthHandler(authHTTPClient)
	proxyHandler := handler.NewProxyHandler(projectClient, dataClient)
	dataAPIHandler := handler.NewPublicAPIHandler(dataClient)
	adminHandler := handler.NewAdminHandler()

	// Setup routing
	mux := http.NewServeMux()

	// Health check
	mux.HandleFunc("/health", healthHandler)

	// Admin routes (internal - TODO: add admin auth)
	mux.HandleFunc("/admin/log-level", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			adminHandler.GetLogLevel(w, r)
		} else if r.Method == http.MethodPut {
			adminHandler.SetLogLevel(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Auth routes (public)
	mux.HandleFunc("/auth/register", authHandler.Register)
	mux.HandleFunc("/auth/login", authHandler.Login)
	mux.HandleFunc("/auth/refresh", authHandler.Refresh)

	// Frontend API routes (protected - requires JWT)
	// For UI: https://mockingcode.dev/projects/{api_key}
	mux.HandleFunc("/projects", proxyHandler.HandleProjects)    // GET list, POST create
	mux.HandleFunc("/projects/", proxyHandler.HandleProjects)   // GET/PUT/DELETE by api_key, collections

	// Protected routes (JWT-based)
	// ВАЖНО: CORS должен быть ПЕРЕД Auth, чтобы preflight и 401 работали корректно
	protectedHandler := middleware.AuthMiddleware(authGRPCClient)(mux)
	protectedHandler = middleware.CORSMiddleware(cfg)(protectedHandler)

	// Public Data API routes (protected by API key) - for developers
	// Pattern: /{api_key}/{collection}[/{id}]
	// Example: https://{api_key}.mockingcode.dev/users (or http://localhost:8080/{api_key}/users)
	publicAPIWithMiddleware := middleware.CORSMiddleware(cfg)(http.HandlerFunc(dataAPIHandler.HandlePublicAPI))
	publicAPIWithMiddleware = middleware.APIKeyMiddleware(projectGRPCClient)(publicAPIWithMiddleware)

	// Main router with intelligent routing
	mainMux := http.NewServeMux()
	mainMux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Check if this is a public API request (starts with api_key pattern: 16 hex chars)
		// Pattern: /{api_key}/{collection}
		pathParts := strings.Split(strings.TrimPrefix(r.URL.Path, "/"), "/")
		
		if len(pathParts) >= 2 && len(pathParts[0]) == 16 {
			// Looks like API key format - route to Public API
			publicAPIWithMiddleware.ServeHTTP(w, r)
			return
		}
		
		// Otherwise route to protected Frontend API
		protectedHandler.ServeHTTP(w, r)
	})

	// Start server
	logger.Info("API Gateway starting",
		slog.String("port", cfg.ServerPort),
		slog.String("auth_grpc", cfg.AuthGRPCURL),
		slog.String("project_grpc", cfg.ProjectGRPCURL),
		slog.String("data_service", cfg.DataServiceURL),
	)

	if err := http.ListenAndServe(":"+cfg.ServerPort, mainMux); err != nil {
		log.Fatal("Failed to start gateway:", err)
	}
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok","service":"gateway"}`))
}
