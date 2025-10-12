package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-mockingcode/project/internal/model"
	"github.com/go-mockingcode/project/internal/pkg/context"
	"github.com/go-mockingcode/project/internal/service"
)

type ProjectHandler struct {
	projectService    *service.ProjectService
	collectionService *service.CollectionService
}

func NewProjectHandler(projectService *service.ProjectService, collectionService *service.CollectionService) *ProjectHandler {
	return &ProjectHandler{
		projectService:    projectService,
		collectionService: collectionService,
	}
}

// HandlerProjects handles /projects endpoint for GET and POST methods
func (h *ProjectHandler) HandlerProjects(w http.ResponseWriter, r *http.Request) {
	userID, err := extractUserID(w, r)
	if err != nil {
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.GetUserProjects(w, r, userID)
	case http.MethodPost:
		h.CreateProject(w, r, userID)
	default:
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{"error": "Method not allowed"})
	}
}

// GetUserProjects godoc
// @Summary Get user projects
// @Description Get list of all projects for authenticated user
// @Tags projects
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]string
// @Router /projects [get]
func (h *ProjectHandler) GetUserProjects(w http.ResponseWriter, r *http.Request, userID int64) {
	projects, err := h.projectService.GetUserProjects(userID)
	if err != nil {
		writeErrorJson(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeSuccessJson(w, http.StatusOK, map[string]any{
		"projects": projects,
		"count":    len(projects),
	})
}

// CreateProject godoc
// @Summary Create new project
// @Description Create a new project for authenticated user
// @Tags projects
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body model.CreateProjectRequest true "Project data"
// @Success 201 {object} model.Project
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /projects [post]
func (h *ProjectHandler) CreateProject(w http.ResponseWriter, r *http.Request, userID int64) {
	var req model.CreateProjectRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErrorJson(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Базовая валидация
	if strings.TrimSpace(req.Name) == "" {
		writeErrorJson(w, http.StatusBadRequest, "Project name is required")
		return
	}

	project, err := h.projectService.CreateProject(userID, &req)
	if err != nil {
		writeErrorJson(w, http.StatusBadRequest, err.Error())
		return
	}

	writeSuccessJson(w, http.StatusCreated, project)
}

// HandleProjectByID handles /projects/{id} endpoint for GET, PUT and DELETE methods
func (h *ProjectHandler) HandleProjectByID(w http.ResponseWriter, r *http.Request) {
	userID, err := extractUserID(w, r)
	if err != nil {
		return
	}

	projectID, err := extractProjectID(w, r)
	if err != nil {
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.GetProject(w, r, userID, projectID)
	case http.MethodPut:
		h.UpdateProject(w, r, userID, projectID)
	case http.MethodDelete:
		h.DeleteProject(w, r, userID, projectID)
	default:
		writeErrorJson(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

// GetProject godoc
// @Summary Get project by ID
// @Description Get specific project by ID
// @Tags projects
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Project ID"
// @Success 200 {object} model.Project
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /projects/{id} [get]
func (h *ProjectHandler) GetProject(w http.ResponseWriter, r *http.Request, userID int64, projectID int64) {
	project, err := h.projectService.GetProject(projectID, userID)
	if err != nil {
		writeErrorJson(w, http.StatusNotFound, err.Error())
		return
	}

	writeSuccessJson(w, http.StatusOK, project)
}

// UpdateProject godoc
// @Summary Update project
// @Description Update project name and description
// @Tags projects
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Project ID"
// @Param request body model.UpdateProjectRequest true "Project data"
// @Success 200 {object} model.Project
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /projects/{id} [put]
func (h *ProjectHandler) UpdateProject(w http.ResponseWriter, r *http.Request, userID int64, projectID int64) {
	var req model.UpdateProjectRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErrorJson(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	project, err := h.projectService.UpdateProject(projectID, userID, &req)
	if err != nil {
		writeErrorJson(w, http.StatusNotFound, err.Error())
		return
	}

	writeSuccessJson(w, http.StatusOK, project)
}

// DeleteProject godoc
// @Summary Delete project
// @Description Delete project and all its collections
// @Tags projects
// @Produce json
// @Security BearerAuth
// @Param id path int true "Project ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /projects/{id} [delete]
func (h *ProjectHandler) DeleteProject(w http.ResponseWriter, r *http.Request, userID int64, projectID int64) {
	if err := h.projectService.DeleteProject(projectID, userID); err != nil {
		writeErrorJson(w, http.StatusBadRequest, err.Error())
		return
	}

	writeSuccessJson(w, http.StatusOK, map[string]string{"message": "Project deleted successfully"})
}

// HandleProjectCollections handles /projects/{id}/collections endpoint for GET and POST methods
func (h *ProjectHandler) HandleProjectCollections(w http.ResponseWriter, r *http.Request) {
	userID, err := extractUserID(w, r)
	if err != nil {
		return
	}

	projectID, err := extractProjectID(w, r)
	if err != nil {
		return
	}

	// Проверяем, что проект принадлежит пользователю
	project, err := h.projectService.GetProject(projectID, userID)
	if err != nil {
		writeErrorJson(w, http.StatusNotFound, "Project not found")
		return
	}
	if project == nil {
		writeErrorJson(w, http.StatusNotFound, "Project not found")
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.GetProjectCollections(w, r, projectID, userID)
	case http.MethodPost:
		h.CreateProjectCollection(w, r, projectID, userID)
	default:
		writeErrorJson(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

// GetProjectCollections godoc
// @Summary Get collections for a project
// @Description Get list of collections for a project
// @Tags collections
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Project ID"
// @Success 200 {object} map[string]interface{}
// @Success 201 {object} model.Collection
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /projects/{id}/collections [get]
func (h *ProjectHandler) GetProjectCollections(w http.ResponseWriter, r *http.Request, projectID int64, userID int64) {
	collections, err := h.collectionService.GetProjectCollections(projectID, userID)
	if err != nil {
		writeErrorJson(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeSuccessJson(w, http.StatusOK, map[string]any{
		"collections": collections,
		"count":       len(collections),
	})
}

// CreateProjectCollection godoc
// @Summary Create collections for a project
// @Description Create a new collection
// @Tags collections
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Project ID"
// @Param request body model.CreateCollectionRequest true "Collection data"
// @Success 200 {object} map[string]interface{}
// @Success 201 {object} model.Collection
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /projects/{id}/collections [post]
func (h *ProjectHandler) CreateProjectCollection(w http.ResponseWriter, r *http.Request, projectID int64, userID int64) {
	var req model.CreateCollectionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErrorJson(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Базовая валидация
	if strings.TrimSpace(req.Name) == "" {
		writeErrorJson(w, http.StatusBadRequest, "Collection name is required")
		return
	}
	if len(req.Fields) == 0 {
		writeErrorJson(w, http.StatusBadRequest, "At least one field is required")
		return
	}

	collection, err := h.collectionService.CreateCollection(projectID, userID, &req)
	if err != nil {
		writeErrorJson(w, http.StatusBadRequest, err.Error())
		return
	}

	writeSuccessJson(w, http.StatusCreated, collection)
}

func extractUserID(w http.ResponseWriter, r *http.Request) (int64, error) {
	userID, err := context.GetUserID(r.Context())
	if err != nil {
		writeErrorJson(w, http.StatusUnauthorized, "User not authenticated")
		return 0, err
	}

	return userID, nil
}

func extractProjectID(w http.ResponseWriter, r *http.Request) (int64, error) {
	// Извлекаем ID из URL
	pathParts := strings.Split(r.URL.Path, "/") // ["", "projects", "1"]

	// Для путей: /projects/{id} и /projects/{id}/collections
	if len(pathParts) < 3 {
		writeErrorJson(w, http.StatusBadRequest, "Invalid project ID")
		return 0, fmt.Errorf("invalid project ID")
	}

	projectID, err := strconv.ParseInt(pathParts[2], 10, 64)
	if err != nil {
		writeErrorJson(w, http.StatusBadRequest, "Invalid project ID")
		return 0, err
	}

	return projectID, nil
}

// TODO refactor to common utils
func writeErrorJson(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

func writeSuccessJson(w http.ResponseWriter, statusCode int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}
