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

// HandleProjects обрабатывает /projects (GET, POST)
func (h *ProjectHandler) HandleProjects(w http.ResponseWriter, r *http.Request) {
	userID, err := context.GetUserID(r.Context())
	if err != nil {
		writeErrorJson(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.getUserProjects(w, r, userID)
	case http.MethodPost:
		h.createProject(w, r, userID)
	default:
		writeErrorJson(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

// getUserProjects возвращает все проекты пользователя
func (h *ProjectHandler) getUserProjects(w http.ResponseWriter, r *http.Request, userID int64) {
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

// createProject создает новый проект
func (h *ProjectHandler) createProject(w http.ResponseWriter, r *http.Request, userID int64) {
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

// HandleProjectByID обрабатывает /projects/{id} (GET, PUT, DELETE)
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

// HandleProjectCollections обрабатывает /projects/{id}/collections (GET, POST)
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
