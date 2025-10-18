package handler

import (
	"net/http"
	"strings"

	"github.com/go-mockingcode/project/internal/service"
)

type APIKeyHandler struct {
	projectService *service.ProjectService
}

func NewAPIKeyHandler(projectService *service.ProjectService) *APIKeyHandler {
	return &APIKeyHandler{
		projectService: projectService,
	}
}

// ValidateAPIKey godoc
// @Summary Validate API Key
// @Description Validate API Key and return project information
// @Tags api-keys
// @Produce json
// @Param apiKey path string true "API Key"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api-keys/{apiKey}/validate [get]
func (h *APIKeyHandler) ValidateAPIKey(w http.ResponseWriter, r *http.Request) {
	// Извлекаем API Key из URL
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 4 {
		writeErrorJson(w, http.StatusBadRequest, "Invalid API Key")
		return
	}

	apiKey := pathParts[3]
	if apiKey == "" {
		writeErrorJson(w, http.StatusBadRequest, "API Key is required")
		return
	}

	project, err := h.projectService.GetProjectByAPIKey(apiKey)
	if err != nil || project == nil {
		writeSuccessJson(w, http.StatusOK, map[string]interface{}{
			"valid":   false,
			"project": nil,
		})
		return
	}

	writeSuccessJson(w, http.StatusOK, map[string]interface{}{
		"valid": true,
		"project": map[string]interface{}{
			"id":       project.ID,
			"user_id":  project.UserID,
			"name":     project.Name,
			"api_key":  project.APIKey,
			"base_url": project.BaseURL,
		},
	})
}
