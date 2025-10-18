package service

import (
	"errors"
	"fmt"

	"github.com/go-mockingcode/data/internal/model"
	"github.com/go-mockingcode/data/internal/repository"
	"github.com/go-mockingcode/models"
)

type DocumentService struct {
	docRepo              *repository.DocumentRepository
	generator            *DataGenerator
	maxDocsPerCollection int
}

func NewDocumentService(docRepo *repository.DocumentRepository, maxDocsPerCollection int) *DocumentService {
	return &DocumentService{
		docRepo:              docRepo,
		generator:            NewDataGenerator(0), // random seed
		maxDocsPerCollection: maxDocsPerCollection,
	}
}

// GetDocuments возвращает документы коллекции
func (s *DocumentService) GetDocuments(projectID int64, collectionName string, opts model.QueryOptions) (*model.DocumentsResponse, error) {
	return s.docRepo.GetDocuments(projectID, collectionName, opts)
}

// GetDocument возвращает документ по ID
func (s *DocumentService) GetDocument(projectID int64, collectionName, documentID string) (*model.MockDocument, error) {
	return s.docRepo.GetDocumentByID(projectID, collectionName, documentID)
}

// CreateDocument создает новый документ
func (s *DocumentService) CreateDocument(projectID int64, collectionName string, data map[string]interface{}) (*model.MockDocument, error) {
	// Проверяем лимит документов
	count, err := s.docRepo.CountDocuments(projectID, collectionName)
	if err != nil {
		return nil, err
	}

	if count >= int64(s.maxDocsPerCollection) {
		return nil, fmt.Errorf("maximum documents limit reached: %d", s.maxDocsPerCollection)
	}

	return s.docRepo.CreateDocument(projectID, collectionName, data)
}

// UpdateDocument обновляет документ
func (s *DocumentService) UpdateDocument(projectID int64, collectionName, documentID string, data map[string]interface{}) (*model.MockDocument, error) {
	return s.docRepo.UpdateDocument(projectID, collectionName, documentID, data)
}

// DeleteDocument удаляет документ
func (s *DocumentService) DeleteDocument(projectID int64, collectionName, documentID string) error {
	return s.docRepo.DeleteDocument(projectID, collectionName, documentID)
}

// FlushCollection очищает коллекцию
func (s *DocumentService) FlushCollection(projectID int64, collectionName string) (int64, error) {
	return s.docRepo.DeleteAllDocuments(projectID, collectionName)
}

// GenerateDocuments генерирует данные по шаблону коллекции
func (s *DocumentService) GenerateDocuments(projectID int64, collection *models.Collection, req *model.GenerateRequest) ([]*model.MockDocument, error) {
	if len(collection.Fields) == 0 {
		return nil, errors.New("collection has no fields defined")
	}

	count := collection.Config.Count
	if req.Count != nil {
		count = *req.Count
	}

	if count <= 0 {
		return nil, errors.New("count must be positive")
	}

	// Проверяем лимит
	currentCount, err := s.docRepo.CountDocuments(projectID, collection.Name)
	if err != nil {
		return nil, err
	}

	if currentCount+int64(count) > int64(s.maxDocsPerCollection) {
		return nil, fmt.Errorf("cannot generate %d documents, would exceed limit of %d", count, s.maxDocsPerCollection)
	}

	// Устанавливаем seed если указан
	if req.Seed != nil {
		s.generator = NewDataGenerator(*req.Seed)
	}

	// Генерируем данные
	generatedData := s.generator.GenerateDocuments(collection.Fields, count)

	// Сохраняем в БД
	var documents []*model.MockDocument
	for _, data := range generatedData {
		doc, err := s.docRepo.CreateDocument(projectID, collection.Name, data)
		if err != nil {
			return nil, err
		}
		documents = append(documents, doc)
	}

	return documents, nil
}
