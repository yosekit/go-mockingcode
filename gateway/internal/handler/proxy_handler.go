package handler

import (
	"log/slog"
	"net/http"
	"strings"

	"github.com/go-mockingcode/gateway/internal/client"
)

type ProxyHandler struct {
	projectClient *client.ProjectClient
	dataClient    *client.DataClient
}

func NewProxyHandler(projectClient *client.ProjectClient, dataClient *client.DataClient) *ProxyHandler {
	return &ProxyHandler{
		projectClient: projectClient,
		dataClient:    dataClient,
	}
}

// HandleProjects proxies requests to project service
func (h *ProxyHandler) HandleProjects(w http.ResponseWriter, r *http.Request) {
	// Gateway path: /projects -> Project Service path: /projects
	path := r.URL.Path
	
	slog.Debug("proxying to project service",
		slog.String("original_path", r.URL.Path),
		slog.String("proxy_path", path),
	)
	
	// Add query parameters
	if r.URL.RawQuery != "" {
		path += "?" + r.URL.RawQuery
	}

	resp, err := h.projectClient.ProxyRequest(r, path)
	if err != nil {
		slog.Error("failed to proxy to project service", slog.String("error", err.Error()))
		writeError(w, http.StatusBadGateway, "Failed to reach project service")
		return
	}
	defer resp.Body.Close()

	slog.Debug("received response from project service", slog.Int("status", resp.StatusCode))

	if err := client.CopyResponse(w, resp); err != nil {
		slog.Error("failed to copy response", slog.String("error", err.Error()))
	}
}

// HandleData proxies requests to data service
func (h *ProxyHandler) HandleData(w http.ResponseWriter, r *http.Request) {
	// Extract path after /data
	path := strings.TrimPrefix(r.URL.Path, "/data")
	
	slog.Debug("proxying to data service",
		slog.String("original_path", r.URL.Path),
		slog.String("proxy_path", path),
	)
	
	// Add query parameters
	if r.URL.RawQuery != "" {
		path += "?" + r.URL.RawQuery
	}

	resp, err := h.dataClient.ProxyRequest(r, path)
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

