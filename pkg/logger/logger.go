package logger

import (
	"log/slog"
	"os"
	"strings"
	"sync"
)

// Config holds logger configuration
type Config struct {
	Level  string // debug, info, warn, error
	Format string // json, text
}

var (
	defaultLevel = new(slog.LevelVar)
	levelMutex   sync.RWMutex
)

// New creates a new configured slog.Logger with dynamic level support
func New(cfg Config) *slog.Logger {
	// Parse and set initial log level
	level := parseLevel(cfg.Level)
	defaultLevel.Set(level)

	// Create handler based on format
	var handler slog.Handler
	opts := &slog.HandlerOptions{
		Level: defaultLevel, // Use LevelVar for dynamic level changes
	}

	if strings.ToLower(cfg.Format) == "json" {
		handler = slog.NewJSONHandler(os.Stdout, opts)
	} else {
		handler = slog.NewTextHandler(os.Stdout, opts)
	}

	return slog.New(handler)
}

// parseLevel converts string to slog.Level
func parseLevel(levelStr string) slog.Level {
	switch strings.ToLower(levelStr) {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn", "warning":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

// Default creates logger with default settings (INFO level, text format)
func Default() *slog.Logger {
	return New(Config{
		Level:  "info",
		Format: "text",
	})
}

// FromEnv creates logger from environment variables
// LOG_LEVEL - debug|info|warn|error (default: info)
// LOG_FORMAT - json|text (default: text)
func FromEnv() *slog.Logger {
	return New(Config{
		Level:  getEnv("LOG_LEVEL", "info"),
		Format: getEnv("LOG_FORMAT", "text"),
	})
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// SetLevel changes the log level at runtime
// level: "debug", "info", "warn", "error"
func SetLevel(level string) {
	levelMutex.Lock()
	defer levelMutex.Unlock()
	
	newLevel := parseLevel(level)
	defaultLevel.Set(newLevel)
	slog.Info("log level changed",
		slog.String("new_level", level),
		slog.Int("level_value", int(newLevel)),
	)
}

// GetLevel returns current log level as string
func GetLevel() string {
	levelMutex.RLock()
	defer levelMutex.RUnlock()
	
	level := defaultLevel.Level()
	switch level {
	case slog.LevelDebug:
		return "debug"
	case slog.LevelInfo:
		return "info"
	case slog.LevelWarn:
		return "warn"
	case slog.LevelError:
		return "error"
	default:
		return "info"
	}
}

