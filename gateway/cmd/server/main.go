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
	authClient := client.NewAuthClient(cfg.AuthServiceURL)
	projectClient := client.NewProjectClient(cfg.ProjectServiceURL)
	dataClient := client.NewDataClient(cfg.DataServiceURL)

	// Initialize handlers
	authHandler := handler.NewAuthHandler(authClient)
	proxyHandler := handler.NewProxyHandler(projectClient, dataClient)
	adminHandler := handler.NewAdminHandler()

	// Setup routing
	mux := http.NewServeMux()

	// Health check
	mux.HandleFunc("/health", healthHandler)

	// Admin routes (public for now - TODO: add admin auth)
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

	// API routes (protected) - proxy to project service
	mux.HandleFunc("/api/", func(w http.ResponseWriter, r *http.Request) {
		slog.Debug("api route handler", slog.String("path", r.URL.Path))
		if strings.HasPrefix(r.URL.Path, "/api/projects") {
			proxyHandler.HandleProjects(w, r)
		} else {
			slog.Warn("no route match", slog.String("path", r.URL.Path))
			http.NotFound(w, r)
		}
	})

	// Data routes (protected) - proxy to data service
	mux.HandleFunc("/data/", proxyHandler.HandleData)

	// Apply middleware chain
	handler := middleware.CORSMiddleware(cfg)(mux)
	handler = middleware.AuthMiddleware(authClient)(handler)

	// Start server
	logger.Info("API Gateway starting",
		slog.String("port", cfg.ServerPort),
		slog.String("auth_service", cfg.AuthServiceURL),
		slog.String("project_service", cfg.ProjectServiceURL),
		slog.String("data_service", cfg.DataServiceURL),
	)
	
	if err := http.ListenAndServe(":"+cfg.ServerPort, handler); err != nil {
		log.Fatal("Failed to start gateway:", err)
	}
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok","service":"gateway"}`))
}

