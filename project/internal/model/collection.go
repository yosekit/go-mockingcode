package model

import "time"

// Collection - коллекция/схема в проекте
type Collection struct {
	ID          int64     `json:"id"`
	ProjectID   int64     `json:"project_id"`
	Name        string    `json:"name"`        // Название коллекции (например, "users")
	Description string    `json:"description"` // Описание коллекции
	Fields      string    `json:"fields"`      // JSON с описанием полей коллекции
	IsActive    bool      `json:"is_active"`   // Активна ли коллекция
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type CreateCollectionRequest struct {
	Name        string         `json:"name" validate:"required,min=1,max=50"`
	Description string         `json:"description" validate:"max=500"`
	Fields      map[string]any `json:"fields" validate:"required"` // Структура полей
}
