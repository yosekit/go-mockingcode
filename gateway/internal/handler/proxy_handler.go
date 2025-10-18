package handler

import (
	"log"
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
	log.Printf("[ProxyHandler] HandleProjects called for: %s", r.URL.Path)
	
	// Extract path after /api
	path := strings.TrimPrefix(r.URL.Path, "/api")
	log.Printf("[ProxyHandler] Proxying to project service, path: %s", path)
	
	// Add query parameters
	if r.URL.RawQuery != "" {
		path += "?" + r.URL.RawQuery
	}

	resp, err := h.projectClient.ProxyRequest(r, path)
	if err != nil {
		log.Printf("[ProxyHandler] Error proxying to project service: %v", err)
		writeError(w, http.StatusBadGateway, "Failed to reach project service")
		return
	}
	defer resp.Body.Close()

	log.Printf("[ProxyHandler] Got response from project service: %d", resp.StatusCode)

	if err := client.CopyResponse(w, resp); err != nil {
		log.Printf("[ProxyHandler] Error copying response: %v", err)
	}
}

// HandleData proxies requests to data service
func (h *ProxyHandler) HandleData(w http.ResponseWriter, r *http.Request) {
	// Extract path after /data
	path := strings.TrimPrefix(r.URL.Path, "/data")
	
	// Add query parameters
	if r.URL.RawQuery != "" {
		path += "?" + r.URL.RawQuery
	}

	resp, err := h.dataClient.ProxyRequest(r, path)
	if err != nil {
		log.Printf("Error proxying to data service: %v", err)
		writeError(w, http.StatusBadGateway, "Failed to reach data service")
		return
	}
	defer resp.Body.Close()

	if err := client.CopyResponse(w, resp); err != nil {
		log.Printf("Error copying response: %v", err)
	}
}

