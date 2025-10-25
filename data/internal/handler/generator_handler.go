package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-mockingcode/data/internal/service"
	"github.com/go-mockingcode/models"
)

type GeneratorHandler struct {
	generator *service.DataGenerator
}

func NewGeneratorHandler() *GeneratorHandler {
	return &GeneratorHandler{
		generator: service.NewDataGenerator(0), // random seed
	}
}

// GenerateRequest запрос на генерацию данных
type GenerateRequest struct {
	Fields []models.FieldTemplate `json:"fields" binding:"required"`
	Count  int                    `json:"count" example:"10"`
	Seed   *uint64                `json:"seed,omitempty" example:"12345"`
}

// GenerateResponse ответ с сгенерированными данными
type GenerateResponse struct {
	Documents []map[string]interface{} `json:"documents"`
	Count     int                      `json:"count"`
}

// HandleGenerate godoc
// @Summary Generate mock data
// @Description Generate mock data based on field schema
// @Tags generator
// @Accept json
// @Produce json
// @Param request body GenerateRequest true "Generation request"
// @Success 200 {object} GenerateResponse
// @Failure 400 {object} map[string]string
// @Router /generate [post]
func (h *GeneratorHandler) HandleGenerate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeErrorJson(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var req GenerateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErrorJson(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if len(req.Fields) == 0 {
		writeErrorJson(w, http.StatusBadRequest, "Fields are required")
		return
	}

	if req.Count <= 0 {
		req.Count = 10 // default
	}

	// Устанавливаем seed если указан
	if req.Seed != nil {
		h.generator = service.NewDataGenerator(*req.Seed)
	}

	// Генерируем данные
	documents := h.generator.GenerateDocuments(req.Fields, req.Count)

	response := GenerateResponse{
		Documents: documents,
		Count:     len(documents),
	}

	writeOrderedJson(w, http.StatusOK, response)
}
