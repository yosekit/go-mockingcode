package config

import (
	"fmt"

	"github.com/go-mockingcode/gateway/internal/pkg/env"
)

type Config struct {
	ServerPort string

	// Service URLs (HTTP - для backward compatibility)
	AuthServiceURL    string
	ProjectServiceURL string
	DataServiceURL    string

	// gRPC URLs
	AuthGRPCURL    string
	ProjectGRPCURL string

	// CORS Settings
	CORSAllowedOrigins []string
	CORSAllowedMethods []string
	CORSAllowedHeaders []string

	// Rate Limiting
	RateLimitEnabled bool
	RateLimitPerMin  int
}

func Load() *Config {
	// Support both direct URL env vars (for Docker) and port-based (for local dev)
	authURL := env.GetString("AUTH_SERVICE_URL", "")
	if authURL == "" {
		authURL = fmt.Sprintf("http://localhost:%s", env.GetString("AUTH_PORT", "8081"))
	}

	projectURL := env.GetString("PROJECT_SERVICE_URL", "")
	if projectURL == "" {
		projectURL = fmt.Sprintf("http://localhost:%s", env.GetString("PROJECT_PORT", "8082"))
	}

	dataURL := env.GetString("DATA_SERVICE_URL", "")
	if dataURL == "" {
		dataURL = fmt.Sprintf("http://localhost:%s", env.GetString("DATA_PORT", "8083"))
	}

	// gRPC URLs
	authGRPCURL := env.GetString("AUTH_GRPC_URL", "localhost:9081")
	projectGRPCURL := env.GetString("PROJECT_GRPC_URL", "localhost:9082")

	return &Config{
		ServerPort: env.GetString("GATEWAY_PORT", "8080"),

		AuthServiceURL:    authURL,
		ProjectServiceURL: projectURL,
		DataServiceURL:    dataURL,

		AuthGRPCURL:    authGRPCURL,
		ProjectGRPCURL: projectGRPCURL,

		CORSAllowedOrigins: []string{"*"}, // TODO: configure properly
		CORSAllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		CORSAllowedHeaders: []string{"Authorization", "Content-Type"},

		RateLimitEnabled: env.GetBool("RATE_LIMIT_ENABLED", false),
		RateLimitPerMin:  env.GetInt("RATE_LIMIT_PER_MIN", 100),
	}
}

