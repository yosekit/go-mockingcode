package service

import (
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

	// Создаем коллекцию с дефолтными значениями
	collection := &model.Collection{
		ProjectID:   projectID,
		Name:        req.Name,
		Description: req.Description,
		Fields:      req.Fields, // Может быть пустым
		Config:      req.Config, // Может быть пустым
		IsActive:    true,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Устанавливаем дефолтный config если не указан
	if collection.Config.Count == 0 {
		collection.Config.Count = 10
	}

	if err := s.collectionRepo.CreateCollection(collection); err != nil {
		return nil, err
	}

	return collection, nil
}

// GetProjectCollections возвращает все коллекции проекта
func (s *CollectionService) GetProjectCollections(projectID int64, userID int64) ([]*model.Collection, error) {
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

// GetCollection возвращает коллекцию по ID
func (s *CollectionService) GetCollection(collectionID int64, projectID int64, userID int64) (*model.Collection, error) {
	// Проверяем что проект принадлежит пользователю
	project, err := s.projectRepo.GetProjectByID(projectID, userID)
	if err != nil {
		return nil, err
	}
	if project == nil {
		return nil, errors.New("project not found")
	}

	return s.collectionRepo.GetCollectionByID(collectionID, projectID)
}

// UpdateCollection обновляет коллекцию
func (s *CollectionService) UpdateCollection(collectionID int64, projectID int64, userID int64, req *model.UpdateCollectionRequest) (*model.Collection, error) {
	// Проверяем что проект принадлежит пользователю
	project, err := s.projectRepo.GetProjectByID(projectID, userID)
	if err != nil {
		return nil, err
	}
	if project == nil {
		return nil, errors.New("project not found")
	}

	// Получаем текущую коллекцию
	collection, err := s.collectionRepo.GetCollectionByID(collectionID, projectID)
	if err != nil {
		return nil, err
	}
	if collection == nil {
		return nil, errors.New("collection not found")
	}

	// Обновляем только переданные поля
	if req.Name != "" {
		collection.Name = req.Name
	}
	if req.Description != "" {
		collection.Description = req.Description
	}
	if req.Fields != nil {
		collection.Fields = req.Fields
	}
	if req.Config != nil {
		collection.Config = *req.Config
	}
	if req.IsActive != nil {
		collection.IsActive = *req.IsActive
	}

	collection.UpdatedAt = time.Now()

	if err := s.collectionRepo.UpdateCollection(collection); err != nil {
		return nil, err
	}

	return collection, nil
}
