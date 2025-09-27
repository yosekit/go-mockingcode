package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
    DBHost     string
    DBPort     string
    DBUser     string
    DBPassword string
    DBName     string
    
    DBMaxOpenConns int
    DBMaxIdleConns int
    DBMaxLifetime  time.Duration

    ServerPort string
    JWTSecret  string
    
    AccessTokenExpiry  time.Duration
    RefreshTokenExpiry time.Duration
}

func Load() *Config {
    return &Config{
        DBHost:     getEnv("DB_HOST", "localhost"),
        DBPort:     getEnv("DB_PORT", "5432"),
        DBUser:     getEnv("DB_USER", "mockuser"),
        DBPassword: getEnv("DB_PASSWORD", "mockpass"),
        DBName:     getEnv("DB_NAME", "mockdb"),

        DBMaxOpenConns: getEnvInt("DB_MAX_OPEN_CONNS", 25),
        DBMaxIdleConns: getEnvInt("DB_MAX_IDLE_CONNS", 25),
        DBMaxLifetime:  getEnvDuration("DB_MAX_LIFETIME", 5*time.Minute),

        ServerPort: getEnv("PORT", "8081"),
        JWTSecret:  getEnv("JWT_SECRET", "super-secret-jwt-key"), // TODO

        AccessTokenExpiry:  getEnvDuration("ACCESS_TOKEN_EXPIRY", 15*time.Minute),
        RefreshTokenExpiry: getEnvDuration("REFRESH_TOKEN_EXPIRY", 7*24*time.Hour),
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

func getEnvDuration(key string, defaultValue time.Duration) time.Duration {
    if value := os.Getenv(key); value != "" {
        if duration, err := time.ParseDuration(value); err == nil {
            return duration
        }
    }
    return defaultValue
}