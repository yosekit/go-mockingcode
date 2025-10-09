package main

import (
	"log"
	"net/http"

	"github.com/go-mockingcode/project/internal/config"
	"github.com/go-mockingcode/project/internal/database"
	"github.com/go-mockingcode/project/internal/handler"
	"github.com/go-mockingcode/project/internal/middleware"
	"github.com/go-mockingcode/project/internal/pkg/auth"
	"github.com/go-mockingcode/project/internal/repository"
	"github.com/go-mockingcode/project/internal/service"
	"github.com/joho/godotenv"
)

func main() {
	// DEV
	if err := godotenv.Load(".env"); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Load Configuration
	cfg := config.Load()

	// Connect to DB
	db, err := database.NewPostgresDB(cfg)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Init DB schemas
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
		collectionRepo,
		cfg.MaxSchemasPerProject,
	)

	// Init Handlers
	projectHandler := handler.NewProjectHandler(projectService, collectionService)

	// Init Auth Client
	authClient := auth.NewAuthClient(cfg.AuthServiceURL)

	// Route Settings
	mux := http.NewServeMux()
	mux.HandleFunc("/health", healthHandler)

	// Middleware Settings
	handlerWithAuth := middleware.AuthMiddleware(authClient)(mux)

	// Handler Settings
	mux.HandleFunc("/projects", projectHandler.HandleProjects)
	mux.HandleFunc("/projects/", projectHandler.HandleProjectByID)
	mux.HandleFunc("/projects/{id}/collections", projectHandler.HandleProjectCollections)

	port := cfg.ServerPort
	log.Printf("Project service starting on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, handlerWithAuth))
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status": "ok", "service": "project"}`))
}
