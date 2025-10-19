package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-mockingcode/project/internal/model"
	"github.com/go-mockingcode/project/internal/service"
)

type CollectionHandler struct {
	projectService    *service.ProjectService
	collectionService *service.CollectionService
}

func NewCollectionHandler(projectService *service.ProjectService, collectionService *service.CollectionService) *CollectionHandler {
	return &CollectionHandler{
		projectService:    projectService,
		collectionService: collectionService,
	}
}

// HandleProjectCollections handles /projects/{id}/collections endpoint for GET and POST methods
func (h *CollectionHandler) HandleProjectCollections(w http.ResponseWriter, r *http.Request) {
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
func (h *CollectionHandler) GetProjectCollections(w http.ResponseWriter, r *http.Request, projectID int64, userID int64) {
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
func (h *CollectionHandler) CreateProjectCollection(w http.ResponseWriter, r *http.Request, projectID int64, userID int64) {
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

// HandleProjectCollections handles /projects/{id}/collections/{collectionId} endpoints for GET, PUT and DELETE methods
func (h *CollectionHandler) HandleProjectCollectionByID(w http.ResponseWriter, r *http.Request) {
	userID, err := extractUserID(w, r)
	if err != nil {
		return
	}

	projectID, collectionID, err := extractProjectAndCollectionID(w, r)
	if err != nil {
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.GetCollection(w, r, projectID, collectionID, userID)
	case http.MethodPut:
		h.UpdateCollection(w, r, projectID, collectionID, userID)
	case http.MethodDelete:
		h.DeleteCollection(w, r, projectID, collectionID, userID)
	default:
		writeErrorJson(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

// GetCollection godoc
// @Summary Get collection by ID
// @Description Get specific collection by ID
// @Tags collections
// @Produce json
// @Security BearerAuth
// @Param projectId path int true "Project ID"
// @Param collectionId path int true "Collection ID"
// @Success 200 {object} model.Collection
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /projects/{projectId}/collections/{collectionId} [get]
func (h *CollectionHandler) GetCollection(w http.ResponseWriter, r *http.Request, projectID, collectionID, userID int64) {
	collection, err := h.collectionService.GetCollection(collectionID, projectID, userID)
	if err != nil {
		writeErrorJson(w, http.StatusNotFound, err.Error())
		return
	}

	writeSuccessJson(w, http.StatusOK, collection)
}

// UpdateCollection godoc
// @Summary Update collection
// @Description Update collection name, description, fields or config
// @Tags collections
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param projectId path int true "Project ID"
// @Param collectionId path int true "Collection ID"
// @Param request body model.UpdateCollectionRequest true "Collection data"
// @Success 200 {object} model.Collection
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /projects/{projectId}/collections/{collectionId} [put]
func (h *CollectionHandler) UpdateCollection(w http.ResponseWriter, r *http.Request, projectID, collectionID, userID int64) {
	var req model.UpdateCollectionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErrorJson(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	collection, err := h.collectionService.UpdateCollection(collectionID, projectID, userID, &req)
	if err != nil {
		writeErrorJson(w, http.StatusBadRequest, err.Error())
		return
	}

	writeSuccessJson(w, http.StatusOK, collection)
}

// DeleteCollection godoc
// @Summary Delete collection
// @Description Delete collection
// @Tags collections
// @Produce json
// @Security BearerAuth
// @Param projectId path int true "Project ID"
// @Param collectionId path int true "Collection ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /projects/{projectId}/collections/{collectionId} [delete]
func (h *CollectionHandler) DeleteCollection(w http.ResponseWriter, r *http.Request, projectID, collectionID, userID int64) {
	if err := h.collectionService.DeleteCollection(collectionID, projectID, userID); err != nil {
		writeErrorJson(w, http.StatusBadRequest, err.Error())
		return
	}

	writeSuccessJson(w, http.StatusOK, map[string]string{"message": "Collection deleted successfully"})
}

func extractProjectAndCollectionID(w http.ResponseWriter, r *http.Request) (int64, int64, error) {
	pathParts := strings.Split(r.URL.Path, "/")

	if len(pathParts) < 5 {
		writeErrorJson(w, http.StatusBadRequest, "Invalid project ID or collection ID")
		return 0, 0, nil
	}

	projectID, err := strconv.ParseInt(pathParts[2], 10, 64)
	if err != nil {
		writeErrorJson(w, http.StatusBadRequest, "Invalid project ID")
		return 0, 0, err
	}

	collectionID, err := strconv.ParseInt(pathParts[4], 10, 64)
	if err != nil {
		writeErrorJson(w, http.StatusBadRequest, "Invalid collection ID")
		return 0, 0, err
	}

	return projectID, collectionID, nil
}
