package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/go-mockingcode/project/internal/model"
	"github.com/go-mockingcode/project/internal/repository"
)

type CollectionService struct {
	projectRepo              *repository.ProjectRepository
	collectionRepo           *repository.CollectionRepository
	maxCollectionsPerProject int
}

func NewCollectionService(collectionRepo *repository.CollectionRepository, maxCollectionsPerProject int) *CollectionService {
	return &CollectionService{
		collectionRepo:           collectionRepo,
		maxCollectionsPerProject: maxCollectionsPerProject,
	}
}

// CreateCollection создает новую коллекцию в проекте
func (s *CollectionService) CreateCollection(projectID int64, userID int64, req *model.CreateCollectionRequest) (*model.Collection, error) {
	// Проверяем что проект существует и принадлежит пользователю
	project, err := s.projectRepo.GetProjectByID(projectID, userID)
	if err != nil {
		return nil, err
	}
	if project == nil {
		return nil, errors.New("project not found")
	}

	// Проверяем лимит коллекций
	collections, err := s.collectionRepo.GetProjectCollections(projectID)
	if err != nil {
		return nil, err
	}

	if len(collections) >= s.maxCollectionsPerProject {
		return nil, fmt.Errorf("maximum collections limit reached: %d", s.maxCollectionsPerProject)
	}

	// Конвертируем fields в JSON
	fieldsJSON, err := json.Marshal(req.Fields)
	if err != nil {
		return nil, fmt.Errorf("invalid fields format: %v", err)
	}

	// Создаем коллекцию
	collection := &model.Collection{
		ProjectID:   projectID,
		Name:        req.Name,
		Description: req.Description,
		Fields:      string(fieldsJSON),
		IsActive:    true,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := s.collectionRepo.CreateCollection(collection); err != nil {
		return nil, err
	}

	return collection, nil
}

// GetProjectCollections возвращает все коллекции проекта
func (s *CollectionService) GetProjectCollections(projectID, userID int64) ([]*model.Collection, error) {
	// Проверяем что проект принадлежит пользователю
	project, err := s.projectRepo.GetProjectByID(projectID, userID)
	if err != nil {
		return nil, err
	}
	if project == nil {
		return nil, errors.New("project not found")
	}

	return s.collectionRepo.GetProjectCollections(projectID)
}
