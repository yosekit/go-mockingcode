package handler

import (
	"encoding/json"
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
	switch r.Method {
	case http.MethodGet:
		h.GetUserProjects(w, r)
	case http.MethodPost:
		h.CreateProject(w, r)
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
func (h *ProjectHandler) GetUserProjects(w http.ResponseWriter, r *http.Request) {
	userID, err := context.GetUserID(r.Context())
	if err != nil {
		writeErrorJson(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	projects, err := h.projectService.GetUserProjects(userID)
	if err != nil {
		writeErrorJson(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeSuccessJson(w, http.StatusOK, map[string]interface{}{
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
func (h *ProjectHandler) CreateProject(w http.ResponseWriter, r *http.Request) {
	userID, err := context.GetUserID(r.Context())
	if err != nil {
		writeErrorJson(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

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

// HandleProjectByID godoc
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
func (h *ProjectHandler) HandleProjectByID(w http.ResponseWriter, r *http.Request) {
	userID, err := context.GetUserID(r.Context())
	if err != nil {
		writeErrorJson(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	// Извлекаем ID из URL
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 3 {
		writeErrorJson(w, http.StatusBadRequest, "Invalid project ID")
		return
	}

	projectID, err := strconv.ParseInt(pathParts[2], 10, 64)
	if err != nil {
		writeErrorJson(w, http.StatusBadRequest, "Invalid project ID")
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.getProject(w, r, userID, projectID)
	// case http.MethodPut:
	//     h.updateProject(w, r, userID, projectID)
	// case http.MethodDelete:
	//     h.deleteProject(w, r, userID, projectID)
	default:
		writeErrorJson(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

// getProject возвращает конкретный проект
func (h *ProjectHandler) getProject(w http.ResponseWriter, r *http.Request, userID, projectID int64) {
	project, err := h.projectService.GetProject(projectID, userID)
	if err != nil {
		writeErrorJson(w, http.StatusNotFound, err.Error())
		return
	}

	writeSuccessJson(w, http.StatusOK, project)
}

// HandleProjectCollections godoc
// @Summary Get or create collections for a project
// @Description Get list of collections for a project or create a new collection
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
// @Router /projects/{id}/collections [post]
func (h *ProjectHandler) HandleProjectCollections(w http.ResponseWriter, r *http.Request) {
	userID, err := context.GetUserID(r.Context())
	if err != nil {
		writeErrorJson(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	// Извлекаем ID из URL
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 4 {
		writeErrorJson(w, http.StatusBadRequest, "Invalid project ID")
		return
	}

	projectID, err := strconv.ParseInt(pathParts[2], 10, 64)
	if err != nil {
		writeErrorJson(w, http.StatusBadRequest, "Invalid project ID")
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
		h.getProjectCollections(w, r, projectID, userID)
	case http.MethodPost:
		h.createProjectCollection(w, r, projectID, userID)
	default:
		writeErrorJson(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

func (h *ProjectHandler) getProjectCollections(w http.ResponseWriter, r *http.Request, projectID int64, userID int64) {
	collections, err := h.collectionService.GetProjectCollections(projectID, userID)
	if err != nil {
		writeErrorJson(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeSuccessJson(w, http.StatusOK, map[string]interface{}{
		"collections": collections,
		"count":       len(collections),
	})
}

func (h *ProjectHandler) createProjectCollection(w http.ResponseWriter, r *http.Request, projectID int64, userID int64) {
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
