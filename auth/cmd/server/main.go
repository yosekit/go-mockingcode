package main

import (
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/go-mockingcode/auth/internal/config"
	"github.com/go-mockingcode/auth/internal/database"
	"github.com/go-mockingcode/auth/internal/handler"
	"github.com/go-mockingcode/auth/internal/repository"
	"github.com/go-mockingcode/auth/internal/service"
	applogger "github.com/go-mockingcode/logger"
	"github.com/joho/godotenv"

	_ "github.com/go-mockingcode/auth/docs"

	httpSwagger "github.com/swaggo/http-swagger"
)

// @title MockingCode Auth Service API
// @version 1.0
// @description Authentication service for MockingCode platform
// @termsOfService http://mockingcode.dev/terms/

// @contact.name API Support
// @contact.url http://mockingcode.dev/support
// @contact.email support@mockingcode.dev

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8081
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
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

	// Init DB schemas
	userRepo := repository.NewUserRepository(db)
	if err := userRepo.InitSchema(); err != nil {
		log.Fatal("Failed to init users schema:", err)
	}
	tokenRepo := repository.NewTokenRepository(db)
	if err := tokenRepo.InitSchema(); err != nil {
		log.Fatal("Failed to init tokens schema:", err)
	}

	// Init Services
	authService := service.NewAuthService(userRepo, tokenRepo, cfg)
	authHandler := handler.NewAuthHandler(authService)

	// Route Settings
	mux := http.NewServeMux()

	// Health Check
	mux.HandleFunc("/health", healthHandler)
	// Swagger
	mux.Handle("/swagger/", httpSwagger.WrapHandler)

	// Auth Handlers
	mux.HandleFunc("/register", authHandler.Register)
	mux.HandleFunc("/login", authHandler.Login)
	mux.HandleFunc("/refresh", authHandler.Refresh)
	mux.HandleFunc("/validate", authHandler.Validate)
	mux.HandleFunc("/logout", authHandler.Logout)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	logger.Info("Auth service starting", slog.String("port", port))
	log.Fatal(http.ListenAndServe(":"+port, mux))
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	w.Write([]byte(`{"status": "ok", "service": "auth"}`))
}
