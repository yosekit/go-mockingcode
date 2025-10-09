package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-mockingcode/auth/internal/config"
	"github.com/go-mockingcode/auth/internal/database"
	"github.com/go-mockingcode/auth/internal/handler"
	"github.com/go-mockingcode/auth/internal/repository"
	"github.com/go-mockingcode/auth/internal/service"
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
	mux.HandleFunc("/health", healthHandler)
	mux.HandleFunc("/register", authHandler.Register)
	mux.HandleFunc("/login", authHandler.Login)
	mux.HandleFunc("/refresh", authHandler.Refresh)
	mux.HandleFunc("/validate", authHandler.Validate)
	mux.HandleFunc("/logout", authHandler.Logout)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	log.Printf("Auth service starting on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	w.Write([]byte(`{"status": "ok", "service": "auth"}`))
}
