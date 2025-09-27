package config

import (
	"fmt"
	"time"

	"github.com/go-mockingcode/project/internal/pkg/env"
)

type Config struct {
	ServerPort string

	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string

	DBMaxOpenConns int
	DBMaxIdleConns int
	DBMaxLifetime  time.Duration

	MaxProjectsPerUser   int
	MaxSchemasPerProject int

	// External Service
	AuthServiceURL string
}

func Load() *Config {
	return &Config{
		ServerPort: env.GetString("PROJECT_PORT", "8082"),

		DBHost:     env.GetString("DB_HOST", "localhost"),
		DBPort:     env.GetString("DB_PORT", "5432"),
		DBUser:     env.GetString("DB_USER", "mockuser"),
		DBPassword: env.GetString("DB_PASSWORD", "mockpass"),
		DBName:     env.GetString("DB_NAME", "mockdb"),

		DBMaxOpenConns: env.GetInt("PROJECT_DB_MAX_OPEN_CONNS", 25),
		DBMaxIdleConns: env.GetInt("PROJECT_DB_MAX_IDLE_CONNS", 25),
		DBMaxLifetime:  env.GetDuration("PROJECT_DB_MAX_LIFETIME", 5*time.Minute),

		MaxProjectsPerUser:   env.GetInt("MAX_PROJECTS_PER_USER", 10),
		MaxSchemasPerProject: env.GetInt("MAX_SCHEMAS_PER_PROJECT", 50),

		AuthServiceURL: fmt.Sprintf("http://localhost:%s", env.GetString("AUTH_PORT", "8081")),
	}
}
