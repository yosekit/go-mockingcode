package config

import (
	"time"

	"github.com/go-mockingcode/data/internal/pkg/env"
)

type DataConfig struct {
	ServerPort string

	MongoURI       string
	MongoDBName    string
	MongoDBTimeout time.Duration

	ProjectPort string

	MaxDocumentsPerCollection int
	DefaultGenerationCount    int
}

func Load() *DataConfig {
	return &DataConfig{
		// Server
		ServerPort: env.GetString("PORT", "8083"),

		// MongoDB
		MongoURI:       env.GetString("MONGO_URI", "mongodb://localhost:27017"),
		MongoDBName:    env.GetString("MONGO_DB_NAME", "mockingcode"),
		MongoDBTimeout: env.GetDuration("MONGO_TIMEOUT", 10*time.Second),

		// External services
		ProjectPort: env.GetString("PROJECT_PORT", "8082"),

		// Application settings
		MaxDocumentsPerCollection: env.GetInt("DATA_MAX_DOCS_PER_COLLECTION", 500),
		DefaultGenerationCount:    env.GetInt("DATA_DEFAULT_GENERATION_COUNT", 10),
	}
}
