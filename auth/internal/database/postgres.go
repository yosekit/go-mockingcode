package database

import (
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/go-mockingcode/auth/internal/config"
	_ "github.com/lib/pq"
)

func NewPostgresDB(cfg *config.Config) (*sql.DB, error) {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
        cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)

	slog.Info("connecting to PostgreSQL",
		slog.String("host", cfg.DBHost),
		slog.String("database", cfg.DBName),
	)
	
	db, err := sql.Open("postgres", connStr)
	if err != nil{
		return nil, fmt.Errorf("failed to open database: %v", err)
	}

	db.SetMaxOpenConns(cfg.DBMaxOpenConns)
	db.SetMaxIdleConns(cfg.DBMaxIdleConns)
	db.SetConnMaxLifetime(cfg.DBMaxLifetime)

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %v", err)
	}

	return db, nil
}