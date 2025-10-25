package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"github.com/go-mockingcode/data/internal/model"
	"github.com/go-mockingcode/data/internal/pkg/context"
	"github.com/go-mockingcode/data/internal/pkg/project"
	"github.com/go-mockingcode/data/internal/service"
)

type DocumentHandler struct {
	docService *service.DocumentService
}

func NewDocumentHandler(docService *service.DocumentService) *DocumentHandler {
	return &DocumentHandler{
		docService: docService,
	}
}

func (h *DocumentHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	pathParts := strings.Split(r.URL.Path, "/")

	// После middleware путь будет: /{collection} или /{collection}/{id}
	if len(pathParts) == 2 && pathParts[1] != "" {
		// /{collection}
		h.HandleCollection(w, r)
	} else if len(pathParts) == 3 && pathParts[1] != "" && pathParts[2] != "" {
		// /{collection}/{id}
		h.HandleDocument(w, r)
	} else {
		writeErrorJson(w, http.StatusNotFound, "Endpoint not found")
	}
}

// HandleCollection godoc
// @Summary Handle collection operations
// @Description Handle GET and POST operations for a collection
// @Tags documents
// @Param api_key path string true "API Key"
// @Param collection path string true "Collection Name"
// @Router /{api_key}/{collection} [get]
// @Router /{api_key}/{collection} [post]
func (h *DocumentHandler) HandleCollection(w http.ResponseWriter, r *http.Request) {
	project, collectionName, err := extractProjectAndCollection(w, r)
	if err != nil {
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.GetDocuments(w, r, project, collectionName)
	case http.MethodPost:
		h.CreateDocument(w, r, project, collectionName)
	case http.MethodDelete:
		h.FlushCollection(w, r, project, collectionName)
	}
}

// GetDocuments godoc
// @Summary Get collection documents
// @Description Get all documents from a collection with pagination
// @Tags documents
// @Produce json
// @Param api_key path string true "API Key"
// @Param collection path string true "Collection Name"
// @Param limit query int false "Limit" default(10)
// @Param offset query int false "Offset" default(0)
// @Param sort query string false "Sort field" default(created_at)
// @Param order query string false "Sort order" default(desc)
// @Success 200 {object} model.DocumentsResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /{api_key}/{collection} [get]
func (h *DocumentHandler) GetDocuments(w http.ResponseWriter, r *http.Request, project *project.ProjectInfo, collectionName string) {
	// Парсим query parameters
	opts := parseQueryOptions(r)

	response, err := h.docService.GetDocuments(project.ID, collectionName, opts)
	if err != nil {
		writeErrorJson(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Преобразуем в чистый формат для публичного API
	cleanDocs := make([]map[string]interface{}, len(response.Documents))
	for i, doc := range response.Documents {
		cleanDocs[i] = doc.ToClean()
	}

	writeOrderedJson(w, http.StatusOK, cleanDocs)
}

// CreateDocument godoc
// @Summary Create document
// @Description Create new document in collection
// @Tags documents
// @Accept json
// @Produce json
// @Param api_key path string true "API Key"
// @Param collection path string true "Collection Name"
// @Param request body map[string]interface{} true "Document data"
// @Success 201 {object} model.DocumentResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /{api_key}/{collection} [post]
func (h *DocumentHandler) CreateDocument(w http.ResponseWriter, r *http.Request, project *project.ProjectInfo, collectionName string) {
	var data map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		writeErrorJson(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if len(data) == 0 {
		writeErrorJson(w, http.StatusBadRequest, "Document data is required")
		return
	}

	document, err := h.docService.CreateDocument(project.ID, collectionName, data)
	if err != nil {
		writeErrorJson(w, http.StatusBadRequest, err.Error())
		return
	}

	// Чистый формат для публичного API
	writeOrderedJson(w, http.StatusCreated, document.ToClean())
}

// HandleDocument godoc
// @Summary Handle document operations
// @Description Handle GET, PUT and DELETE operations for a document
// @Tags documents
// @Param api_key path string true "API Key"
// @Param collection path string true "Collection Name"
// @Param id path string true "Document ID"
// @Router /{api_key}/{collection}/{id} [get]
// @Router /{api_key}/{collection}/{id} [put]
// @Router /{api_key}/{collection}/{id} [delete]
func (h *DocumentHandler) HandleDocument(w http.ResponseWriter, r *http.Request) {
	project, collectionName, err := extractProjectAndCollection(w, r)
	if err != nil {
		return
	}

	documentID := extractDocumentID(r)
	if documentID == "" {
		writeErrorJson(w, http.StatusBadRequest, "Document ID is required")
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.GetDocument(w, r, project, collectionName, documentID)
	case http.MethodPut:
		h.UpdateDocument(w, r, project, collectionName, documentID)
	case http.MethodDelete:
		h.DeleteDocument(w, r, project, collectionName, documentID)
	}
}

// GetDocument godoc
// @Summary Get document by ID
// @Description Get specific document by ID
// @Tags documents
// @Produce json
// @Param api_key path string true "API Key"
// @Param collection path string true "Collection Name"
// @Param id path string true "Document ID"
// @Success 200 {object} model.DocumentResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /{api_key}/{collection}/{id} [get]
func (h *DocumentHandler) GetDocument(w http.ResponseWriter, r *http.Request, project *project.ProjectInfo, collectionName, documentID string) {
	document, err := h.docService.GetDocument(project.ID, collectionName, documentID)
	if err != nil {
		writeErrorJson(w, http.StatusInternalServerError, err.Error())
		return
	}

	if document == nil {
		writeErrorJson(w, http.StatusNotFound, "Document not found")
		return
	}

	// Чистый формат для публичного API
	writeOrderedJson(w, http.StatusOK, document.ToClean())
}

// UpdateDocument godoc
// @Summary Update document
// @Description Update existing document
// @Tags documents
// @Accept json
// @Produce json
// @Param api_key path string true "API Key"
// @Param collection path string true "Collection Name"
// @Param id path string true "Document ID"
// @Param request body map[string]interface{} true "Document data"
// @Success 200 {object} model.DocumentResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /{api_key}/{collection}/{id} [put]
func (h *DocumentHandler) UpdateDocument(w http.ResponseWriter, r *http.Request, project *project.ProjectInfo, collectionName, documentID string) {
	var data map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		writeErrorJson(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	document, err := h.docService.UpdateDocument(project.ID, collectionName, documentID, data)
	if err != nil {
		writeErrorJson(w, http.StatusBadRequest, err.Error())
		return
	}

	if document == nil {
		writeErrorJson(w, http.StatusNotFound, "Document not found")
		return
	}

	// Чистый формат для публичного API
	writeOrderedJson(w, http.StatusOK, document.ToClean())
}

// DeleteDocument godoc
// @Summary Delete document
// @Description Delete document from collection
// @Tags documents
// @Produce json
// @Param api_key path string true "API Key"
// @Param collection path string true "Collection Name"
// @Param id path string true "Document ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /{api_key}/{collection}/{id} [delete]
func (h *DocumentHandler) DeleteDocument(w http.ResponseWriter, r *http.Request, project *project.ProjectInfo, collectionName, documentID string) {
	if err := h.docService.DeleteDocument(project.ID, collectionName, documentID); err != nil {
		writeErrorJson(w, http.StatusBadRequest, err.Error())
		return
	}

	writeSuccessJson(w, http.StatusOK, map[string]string{"message": "Document deleted successfully"})
}

// FlushCollection godoc
// @Summary Flush collection
// @Description Delete all documents from collection
// @Tags documents
// @Produce json
// @Param api_key path string true "API Key"
// @Param collection path string true "Collection Name"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /{api_key}/{collection} [delete]
func (h *DocumentHandler) FlushCollection(w http.ResponseWriter, r *http.Request, project *project.ProjectInfo, collectionName string) {
	deletedCount, err := h.docService.FlushCollection(project.ID, collectionName)
	if err != nil {
		writeErrorJson(w, http.StatusBadRequest, err.Error())
		return
	}

	writeSuccessJson(w, http.StatusOK, map[string]interface{}{
		"message":       "Collection flushed successfully",
		"deleted_count": deletedCount,
	})
}

func extractProjectAndCollection(w http.ResponseWriter, r *http.Request) (*project.ProjectInfo, string, error) {
	project, err := context.GetProjectInfo(r.Context())
	if err != nil {
		writeErrorJson(w, http.StatusUnauthorized, "Project not authenticated")
		return nil, "", err
	}

	pathParts := strings.Split(r.URL.Path, "/")
	
	// Определяем, откуда пришел запрос - от gateway или напрямую
	var collectionName string
	if len(pathParts) == 2 {
		// /{collection} - от gateway
		collectionName = pathParts[1]
	} else if len(pathParts) >= 3 {
		// /{api_key}/{collection} или /{api_key}/{collection}/... - напрямую
		collectionName = pathParts[2]
	} else {
		writeErrorJson(w, http.StatusBadRequest, "Invalid collection name")
		return nil, "", fmt.Errorf("invalid collection name")
	}

	if collectionName == "" {
		writeErrorJson(w, http.StatusBadRequest, "Collection name is required")
		return nil, "", fmt.Errorf("collection name is required")
	}

	return project, collectionName, nil
}

func extractDocumentID(r *http.Request) string {
	pathParts := strings.Split(r.URL.Path, "/")
	
	// Определяем, откуда пришел запрос - от gateway или напрямую
	if len(pathParts) == 3 {
		// /{collection}/{id} - от gateway
		return pathParts[2]
	} else if len(pathParts) >= 4 {
		// /{api_key}/{collection}/{id} - напрямую
		return pathParts[3]
	}
	return ""
}

func parseQueryOptions(r *http.Request) model.QueryOptions {
	opts := model.QueryOptions{}

	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if limit, err := strconv.ParseInt(limitStr, 10, 64); err == nil {
			opts.Limit = &limit
		}
	}

	if offsetStr := r.URL.Query().Get("offset"); offsetStr != "" {
		if offset, err := strconv.ParseInt(offsetStr, 10, 64); err == nil {
			opts.Offset = &offset
		}
	}

	if sort := r.URL.Query().Get("sort"); sort != "" {
		opts.Sort = sort
	}

	if order := r.URL.Query().Get("order"); order != "" {
		opts.Order = order
	}

	return opts
}

func writeErrorJson(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

func writeSuccessJson(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

// writeOrderedJson записывает JSON с гарантированным порядком полей (id первым)
func writeOrderedJson(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	
	// Для массива или одного документа формируем упорядоченный JSON
	var buf bytes.Buffer
	
	switch v := data.(type) {
	case []map[string]interface{}:
		// Массив документов
		buf.WriteString("[")
		for i, doc := range v {
			if i > 0 {
				buf.WriteString(",")
			}
			buf.Write(orderDocument(doc))
		}
		buf.WriteString("]")
	case map[string]interface{}:
		// Один документ
		buf.Write(orderDocument(v))
	default:
		// Fallback на обычный JSON
		json.NewEncoder(w).Encode(data)
		return
	}
	
	w.Write(buf.Bytes())
}

// orderDocument упорядочивает поля документа (id первым, остальные в алфавитном порядке)
func orderDocument(doc map[string]interface{}) []byte {
	// Собираем ключи
	keys := make([]string, 0, len(doc))
	for k := range doc {
		keys = append(keys, k)
	}
	
	// Сортируем с приоритетом для id
	sort.Slice(keys, func(i, j int) bool {
		if keys[i] == "id" {
			return true
		}
		if keys[j] == "id" {
			return false
		}
		return keys[i] < keys[j]
	})
	
	// Строим JSON вручную
	result := "{"
	for i, k := range keys {
		if i > 0 {
			result += ","
		}
		result += fmt.Sprintf(`"%s":`, k)
		valueJSON, _ := json.Marshal(doc[k])
		result += string(valueJSON)
	}
	result += "}"
	
	return []byte(result)
}
