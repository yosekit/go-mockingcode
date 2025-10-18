package model

import "github.com/go-mockingcode/models"

type Project = models.Project

// CreateProjectRequest represents project creation data
type CreateProjectRequest struct {
	Name        string `json:"name" validate:"required,min=1,max=100"`
	Description string `json:"description" validate:"max=500"`
}

// UpdateProjectRequest represents project update data
type UpdateProjectRequest struct {
	Name        string `json:"name" validate:"required,min=1,max=100"`
	Description string `json:"description" validate:"max=500"`
}

// ProjectResponse represents a project along with its collections
type ProjectResponse struct {
	Project     *Project      `json:"project"`
	Collections []*Collection `json:"collections,omitempty"` // Опционально: коллекции проекта
}
