package config

import (
	"os"
	"strconv"
)

type Config struct {
    DBHost     string
    DBPort     string
    DBUser     string
    DBPassword string
    DBName     string
    
    ServerPort string
    JWTSecret  string
    
    RedisURL   string
}

func Load() *Config {
    return &Config{
        DBHost:     getEnv("DB_HOST", "localhost"),
        DBPort:     getEnv("DB_PORT", "5432"),
        DBUser:     getEnv("DB_USER", "mockuser"),
        DBPassword: getEnv("DB_PASSWORD", "mockpass"),
        DBName:     getEnv("DB_NAME", "mockdb"),
        
        ServerPort: getEnv("AUTH_PORT", "8081"),
        JWTSecret:  getEnv("JWT_SECRET", "super-secret-jwt-key"), // TODO
        
        RedisURL:   getEnv("REDIS_URL", "localhost:6379"),
    }
}

func getEnv(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
    if value := os.Getenv(key); value != "" {
        if intValue, err := strconv.Atoi(value); err == nil {
            return intValue
        }
    }
    return defaultValue
}