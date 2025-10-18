package handler

import (
	"log/slog"
	"net/http"

	applogger "github.com/go-mockingcode/logger"
)

type AdminHandler struct{}

func NewAdminHandler() *AdminHandler {
	return &AdminHandler{}
}

// SetLogLevel handles runtime log level changes
// PUT /admin/log-level?level=debug
func (h *AdminHandler) SetLogLevel(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	level := r.URL.Query().Get("level")
	if level == "" {
		writeError(w, http.StatusBadRequest, "level parameter is required")
		return
	}

	// Validate level
	validLevels := map[string]bool{
		"debug": true,
		"info":  true,
		"warn":  true,
		"error": true,
	}

	if !validLevels[level] {
		writeError(w, http.StatusBadRequest, "Invalid log level. Use: debug, info, warn, error")
		return
	}

	// Change log level
	applogger.SetLevel(level)

	slog.Info("log level updated via admin endpoint",
		slog.String("new_level", level),
		slog.String("remote_addr", r.RemoteAddr),
	)

	writeJSON(w, http.StatusOK, map[string]string{
		"message": "Log level updated successfully",
		"level":   level,
	})
}

// GetLogLevel returns current log level
// GET /admin/log-level
func (h *AdminHandler) GetLogLevel(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	currentLevel := applogger.GetLevel()

	writeJSON(w, http.StatusOK, map[string]string{
		"level": currentLevel,
	})
}

