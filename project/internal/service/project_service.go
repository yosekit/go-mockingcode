package service

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/go-mockingcode/project/internal/model"
	"github.com/go-mockingcode/project/internal/repository"
)

type ProjectService struct {
	projectRepo        *repository.ProjectRepository
	maxProjectsPerUser int
	baseURLFormat      string
}

func NewProjectService(projectRepo *repository.ProjectRepository, maxProjectsPerUser int, baseURLFormat string) *ProjectService {
	return &ProjectService{
		projectRepo:        projectRepo,
		maxProjectsPerUser: maxProjectsPerUser,
		baseURLFormat:      baseURLFormat,
	}
}

// CreateProject создает новый проект для пользователя
func (s *ProjectService) CreateProject(userID int64, req *model.CreateProjectRequest) (*model.Project, error) {
	// Проверяем лимит проектов
	userProjects, err := s.projectRepo.GetUserProjects(userID)
	if err != nil {
		return nil, err
	}

	if len(userProjects) >= s.maxProjectsPerUser {
		return nil, fmt.Errorf("maximum projects limit reached: %d", s.maxProjectsPerUser)
	}

	// Генерируем уникальный API Key
	apiKey, err := s.generateAPIKey()
	if err != nil {
		return nil, err
	}

	// Создаем проект
	project := &model.Project{
		UserID:      userID,
		Name:        req.Name,
		Description: req.Description,
		APIKey:      apiKey,
		BaseURL:     s.generateBaseURL(apiKey),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := s.projectRepo.CreateProject(project); err != nil {
		return nil, err
	}

	return project, nil
}

// GetUserProjects возвращает все проекты пользователя
func (s *ProjectService) GetUserProjects(userID int64) ([]*model.Project, error) {
	return s.projectRepo.GetUserProjects(userID)
}

// GetProject возвращает проект по ID (с проверкой владельца)
func (s *ProjectService) GetProject(projectID int64, userID int64) (*model.Project, error) {
	project, err := s.projectRepo.GetProjectByID(projectID, userID)
	if err != nil {
		return nil, err
	}
	if project == nil {
		return nil, errors.New("project not found")
	}
	return project, nil
}

// UpdateProject обновляет проект
func (s *ProjectService) UpdateProject(projectID int64, userID int64, req *model.UpdateProjectRequest) (*model.Project, error) {
	// Получаем проект (проверяем владельца)
	project, err := s.projectRepo.GetProjectByID(projectID, userID)
	if err != nil {
		return nil, err
	}
	if project == nil {
		return nil, errors.New("project not found")
	}

	// Обновляем поля
	project.Name = req.Name
	project.Description = req.Description
	project.UpdatedAt = time.Now()

	// Сохраняем в БД
	if err := s.projectRepo.UpdateProject(project); err != nil {
		return nil, err
	}

	return project, nil
}

// DeleteProject удаляет проект (каскадно удаляет коллекции)
func (s *ProjectService) DeleteProject(projectID int64, userID int64) error {
	// Проверяем что проект существует и принадлежит пользователю
	project, err := s.projectRepo.GetProjectByID(projectID, userID)
	if err != nil {
		return err
	}
	if project == nil {
		return errors.New("project not found")
	}

	if err := s.projectRepo.DeleteProject(projectID); err != nil {
		return err
	}

	return nil
}

// GenerateAPIKey создает случайный API Key
func (s *ProjectService) generateAPIKey() (string, error) {
	// TODO more safety, avoid ambiguity, keep shorter
	bytes := make([]byte, 8) // 64 бит
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// GenerateBaseURL создает URL для проекта
func (s *ProjectService) generateBaseURL(apiKey string) string {
	return strings.ReplaceAll(s.baseURLFormat, "{api_key}", apiKey)
}
