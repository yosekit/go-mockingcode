package model

import "time"

// Project represents a user project
type Project struct {
	ID          int64     `json:"id"`
	UserID      int64     `json:"user_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	APIKey      string    `json:"api_key"`  // Уникальный ключ для доступа к API проекта
	BaseURL     string    `json:"base_url"` // https://<api_key>.mockingcode.org
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

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
