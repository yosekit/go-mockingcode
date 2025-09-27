package config

import (
	"time"

	"github.com/go-mockingcode/auth/internal/pkg/env"
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

	JWTSecret string

	AccessTokenExpiry  time.Duration
	RefreshTokenExpiry time.Duration
}

func Load() *Config {
	return &Config{
		DBHost:     env.GetString("DB_HOST", "localhost"),
		DBPort:     env.GetString("DB_PORT", "5432"),
		DBUser:     env.GetString("DB_USER", "mockuser"),
		DBPassword: env.GetString("DB_PASSWORD", "mockpass"),
		DBName:     env.GetString("DB_NAME", "mockdb"),

		DBMaxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 25),
		DBMaxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 25),
		DBMaxLifetime:  env.GetDuration("DB_MAX_LIFETIME", 5*time.Minute),

		ServerPort: env.GetString("PORT", "8081"),
		JWTSecret:  env.GetString("JWT_SECRET", "super-secret-jwt-key"), // TODO

		AccessTokenExpiry:  env.GetDuration("ACCESS_TOKEN_EXPIRY", 15*time.Minute),
		RefreshTokenExpiry: env.GetDuration("REFRESH_TOKEN_EXPIRY", 7*24*time.Hour),
	}
}
