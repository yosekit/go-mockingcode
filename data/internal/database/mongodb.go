package database

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/go-mockingcode/data/internal/config"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func NewMongoDB(cfg *config.DataConfig) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), cfg.MongoDBTimeout)
	defer cancel()

	clientOptions := options.Client().
		ApplyURI(cfg.MongoURI).
		SetTimeout(cfg.MongoDBTimeout)

	slog.Info("connecting to MongoDB", slog.String("uri", cfg.MongoURI))

	client, err := mongo.Connect(clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %v", err)
	}

	// Проверяем соединение
	if err := client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %v", err)
	}

	slog.Info("Successfully connected to MongoDB")
	return client, nil
}
