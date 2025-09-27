package main

import (
	"log"
	"net/http"

	"github.com/go-mockingcode/project/internal/config"
	"github.com/go-mockingcode/project/internal/database"
	"github.com/go-mockingcode/project/internal/repository"
	"github.com/joho/godotenv"
)

func main() {
	// DEV
	if err := godotenv.Load("../.env"); err != nil {
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

	// Route Settings
	mux := http.NewServeMux()
	mux.HandleFunc("/health", healthHandler)

	port := cfg.ServerPort
	log.Printf("Project service starting on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status": "ok", "service": "project"}`))
}
