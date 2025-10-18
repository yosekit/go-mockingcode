package config

import (
	"time"

	"github.com/go-mockingcode/data/internal/pkg/env"
)

type DataConfig struct {
	ServerPort string

	MongoDBPort    string
	MongoDBName    string
	MongoDBTimeout time.Duration

	ProjectPort string

	MaxDocumentsPerCollection int
	DefaultGenerationCount    int
}

func Load() *DataConfig {
	return &DataConfig{
		// Server
		ServerPort: env.GetString("DATA_PORT", "8083"),

		// MongoDB
		MongoDBPort:    env.GetString("MONGO_PORT", "27017"), // mongodb://localhost:port
		MongoDBName:    env.GetString("MONGO_DB", "mockingcode"),
		MongoDBTimeout: env.GetDuration("MONGO_TIMEOUT", 10*time.Second),

		// External services
		ProjectPort: env.GetString("PROJECT_PORT", "8082"),

		// Application settings
		MaxDocumentsPerCollection: env.GetInt("DATA_MAX_DOCS_PER_COLLECTION", 500),
		DefaultGenerationCount:    env.GetInt("DATA_DEFAULT_GENERATION_COUNT", 10),
	}
}
