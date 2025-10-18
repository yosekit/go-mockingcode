package database

import (
	"context"
	"fmt"
	"log"

	"github.com/go-mockingcode/data/internal/config"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func NewMongoDB(cfg *config.DataConfig) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), cfg.MongoDBTimeout)
	defer cancel()

	// TODO Get application host
	uri := "mongodb://localhost:" + cfg.MongoDBPort

	clientOptions := options.Client().
		ApplyURI(uri).
		SetTimeout(cfg.MongoDBTimeout)

	client, err := mongo.Connect(clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %v", err)
	}

	// Проверяем соединение
	if err := client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %v", err)
	}

	log.Println("Data service: Successfully connected to MongoDB")
	return client, nil
}
