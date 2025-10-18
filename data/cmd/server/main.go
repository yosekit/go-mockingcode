package main

import (
	"context"
	"log"
	"net/http"

	"github.com/go-mockingcode/data/internal/config"
	"github.com/go-mockingcode/data/internal/database"
	"github.com/go-mockingcode/data/internal/handler"
	"github.com/go-mockingcode/data/internal/middleware"
	"github.com/go-mockingcode/data/internal/pkg/project"
	"github.com/go-mockingcode/data/internal/repository"
	"github.com/go-mockingcode/data/internal/service"
	"github.com/joho/godotenv"

	_ "github.com/go-mockingcode/data/docs"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title MockingCode Data Service API
// @version 1.0
// @description Data service for MockingCode platform
// @termsOfService http://mockingcode.dev/terms/

// @contact.name API Support
// @contact.url http://mockingcode.dev/support
// @contact.email support@mockingcode.dev

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8083
// @BasePath /
func main() {
	// Load .env
	if err := godotenv.Load(".env"); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Load Config
	cfg := config.Load()

	// Connect to MongoDB
	client, err := database.NewMongoDB(cfg)
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}
	defer func() {
		if err := client.Disconnect(context.Background()); err != nil {
			log.Printf("Error disconnecting from MongoDB: %v", err)
		}
	}()

	// Init Repositories
	docRepo := repository.NewDocumentRepository(client, cfg.MongoDBName)

	// Init Services
	docService := service.NewDocumentService(docRepo, cfg.MaxDocumentsPerCollection)

	// Init Handlers
	docHandler := handler.NewDocumentHandler(docService)

	// Init Project Client
	projectClient := project.NewProjectClient("http://localhost:" + cfg.ProjectPort) // TODO

	// Route Settings
	mux := http.NewServeMux()

	// Public Handlers
	mux.HandleFunc("/health", healthHandler)

	mux.Handle("/", docHandler)

	mux.Handle("/swagger/", httpSwagger.WrapHandler)

	// Application Handlers
	// mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	// 	pathParts := strings.Split(r.URL.Path, "/")

	// 	// После middleware путь будет: /{collection} или /{collection}/{id}
	// 	if len(pathParts) == 2 && pathParts[1] != "" {
	// 		// /{collection}
	// 		docHandler.HandleCollection(w, r)
	// 	} else if len(pathParts) == 3 && pathParts[1] != "" && pathParts[2] != "" {
	// 		// /{collection}/{id}
	// 		docHandler.HandleDocument(w, r)
	// 	} else {
	// 		writeErrorJson(w, http.StatusNotFound, "Endpoint not found")
	// 	}
	// })

	// Middleware Settings
	handlerWithAuth := middleware.AuthMiddleware(projectClient)(mux)

	port := cfg.ServerPort
	log.Printf("Data service starting on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, handlerWithAuth))
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status": "ok", "service": "data"}`))
}
