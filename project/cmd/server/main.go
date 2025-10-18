package main

import (
	"log"
	"log/slog"
	"net"
	"net/http"
	"os"

	"github.com/go-mockingcode/project/internal/config"
	"github.com/go-mockingcode/project/internal/database"
	projectgrpc "github.com/go-mockingcode/project/internal/grpc"
	"github.com/go-mockingcode/project/internal/handler"
	"github.com/go-mockingcode/project/internal/middleware"
	"github.com/go-mockingcode/project/internal/repository"
	"github.com/go-mockingcode/project/internal/service"
	applogger "github.com/go-mockingcode/logger"
	pb "github.com/go-mockingcode/proto"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"

	_ "github.com/go-mockingcode/project/docs"

	httpSwagger "github.com/swaggo/http-swagger"
)

// @title MockingCode Project Service API
// @version 1.0
// @description Project management service for MockingCode platform
// @termsOfService http://mockingcode.dev/terms/

// @contact.name API Support
// @contact.url http://mockingcode.dev/support
// @contact.email support@mockingcode.dev

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8082
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Enter the token with the `Bearer ` prefix, e.g. "Bearer abcde12345"
func main() {
	// DEV
	if err := godotenv.Load(".env"); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Load Configuration
	cfg := config.Load()

	// Initialize logger
	logger := applogger.FromEnv()
	slog.SetDefault(logger)

	// Connect to DB
	db, err := database.NewPostgresDB(cfg)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Init Repositories
	projectRepo := repository.NewProjectRepository(db)
	if err := projectRepo.InitSchema(); err != nil {
		log.Fatal("Failed to init projects schema:", err)
	}
	collectionRepo := repository.NewCollectionRepository(db)
	if err := collectionRepo.InitSchema(); err != nil {
		log.Fatal("Failed to init collections schema:", err)
	}

	// Init Services
	projectService := service.NewProjectService(
		projectRepo,
		cfg.MaxProjectsPerUser,
		cfg.BaseURLFormat,
	)
	collectionService := service.NewCollectionService(
		projectRepo,
		collectionRepo,
		cfg.MaxSchemasPerProject,
	)

	// Init Handlers
	projectHandler := handler.NewProjectHandler(projectService)
	collectionHandler := handler.NewCollectionHandler(projectService, collectionService)
	apiKeyHandler := handler.NewAPIKeyHandler(projectService)

	// Route Settings
	mux := http.NewServeMux()

	// Public Handlers
	mux.HandleFunc("/health", healthHandler)
	mux.Handle("/swagger/", httpSwagger.WrapHandler)

	// Application Handlers (protected by Gateway)
	mux.HandleFunc("/projects", projectHandler.HandleProjects)
	mux.HandleFunc("/projects/{id}", projectHandler.HandleProjectByID)
	mux.HandleFunc("/projects/{id}/collections", collectionHandler.HandleProjectCollections)
	mux.HandleFunc("/projects/{id}/collections/{collectionId}", collectionHandler.HandleProjectCollectionByID)

	// API Keys validation (used by Data service, not through Gateway)
	mux.HandleFunc("/api-keys/", apiKeyHandler.ValidateAPIKey)

	// Middleware Settings - extract user ID from X-User-ID header (set by Gateway)
	handlerWithUserID := middleware.UserIDMiddleware(mux)

	port := cfg.ServerPort
	grpcPort := os.Getenv("GRPC_PORT")
	if grpcPort == "" {
		grpcPort = "9082"
	}

	logger.Info("Project service starting",
		slog.String("http_port", port),
		slog.String("grpc_port", grpcPort),
		slog.String("mode", "API Gateway Pattern - trusting X-User-ID header"),
	)

	// Start gRPC server in goroutine
	go func() {
		lis, err := net.Listen("tcp", ":"+grpcPort)
		if err != nil {
			log.Fatalf("Failed to listen on gRPC port: %v", err)
		}

		grpcServer := grpc.NewServer()
		pb.RegisterProjectServiceServer(grpcServer, projectgrpc.NewProjectGRPCServer(projectService, collectionService))

		slog.Info("gRPC server starting", slog.String("port", grpcPort))
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Failed to serve gRPC: %v", err)
		}
	}()

	// Start HTTP server
	log.Fatal(http.ListenAndServe(":"+port, handlerWithUserID))
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status": "ok", "service": "project"}`))
}
