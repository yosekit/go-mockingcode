package handler

import (
	"log/slog"
	"net/http"
	"strings"

	"github.com/go-mockingcode/gateway/internal/client"
)

type PublicAPIHandler struct {
	dataClient *client.DataClient
}

func NewPublicAPIHandler(dataClient *client.DataClient) *PublicAPIHandler {
	return &PublicAPIHandler{
		dataClient: dataClient,
	}
}

// HandlePublicAPI handles public data API requests: /{api_key}/{collection}[/{id}]
func (h *PublicAPIHandler) HandlePublicAPI(w http.ResponseWriter, r *http.Request) {
	// Extract path: /{api_key}/{collection}[/{id}]
	// Remove /{api_key} prefix to get /{collection}[/{id}]
	pathParts := strings.Split(strings.TrimPrefix(r.URL.Path, "/"), "/")
	
	if len(pathParts) < 2 {
		writeError(w, http.StatusBadRequest, "Invalid path")
		return
	}

	// pathParts: ["{api_key}", "collection", ...optional id]
	// We need: /{collection}[/{id}] for Data Service
	dataPath := "/" + strings.Join(pathParts[1:], "/")
	
	slog.Debug("proxying to data service (public API)",
		slog.String("original_path", r.URL.Path),
		slog.String("data_path", dataPath),
	)
	
	// Add query parameters
	if r.URL.RawQuery != "" {
		dataPath += "?" + r.URL.RawQuery
	}

	// Proxy to Data Service
	resp, err := h.dataClient.ProxyRequest(r, dataPath)
	if err != nil {
		slog.Error("failed to proxy to data service", slog.String("error", err.Error()))
		writeError(w, http.StatusBadGateway, "Failed to reach data service")
		return
	}
	defer resp.Body.Close()

	slog.Debug("received response from data service", slog.Int("status", resp.StatusCode))

	if err := client.CopyResponse(w, resp); err != nil {
		slog.Error("failed to copy response", slog.String("error", err.Error()))
	}
}

