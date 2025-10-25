package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"

	"github.com/go-mockingcode/data/internal/config"
	"github.com/go-mockingcode/data/internal/database"
	"github.com/go-mockingcode/data/internal/handler"
	"github.com/go-mockingcode/data/internal/middleware"
	"github.com/go-mockingcode/data/internal/repository"
	"github.com/go-mockingcode/data/internal/service"
	applogger "github.com/go-mockingcode/logger"
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

	// Initialize logger
	logger := applogger.FromEnv()
	slog.SetDefault(logger)

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
	generatorHandler := handler.NewGeneratorHandler()

	// Route Settings
	mux := http.NewServeMux()

	// Public Handlers
	mux.HandleFunc("/health", healthHandler)
	mux.HandleFunc("/generate", generatorHandler.HandleGenerate)

	mux.Handle("/", docHandler)

	mux.Handle("/swagger/", httpSwagger.WrapHandler)

	// Middleware Settings - extract user ID from X-User-ID header (set by Gateway)
	handlerWithUserID := middleware.ProjectInfoMiddleware()(mux)

	port := cfg.ServerPort
	logger.Info("Data service starting",
		slog.String("port", port),
		slog.String("mode", "API Gateway Pattern - trusting X-User-ID header"),
	)
	log.Fatal(http.ListenAndServe(":"+port, handlerWithUserID))
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status": "ok", "service": "data"}`))
}
