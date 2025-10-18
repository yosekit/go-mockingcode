package model

import "github.com/go-mockingcode/models"

type Collection = models.Collection
type CollectionConfig = models.CollectionConfig
type FieldTemplate = models.FieldTemplate

// CreateCollectionRequest represents collection creation data
type CreateCollectionRequest struct {
	Name        string           `json:"name" binding:"required,min=1,max=50" example:"users"`
	Description string           `json:"description" binding:"max=500" example:"User data"`
	Fields      []FieldTemplate  `json:"fields,omitempty"` // Опционально при создании
	Config      CollectionConfig `json:"config,omitempty"` // Опционально при создании
}

// UpdateCollectionRequest represents collection updating data
type UpdateCollectionRequest struct {
	Name        string            `json:"name,omitempty" binding:"max=50" example:"updated-users"`
	Description string            `json:"description,omitempty" binding:"max=500" example:"Updated user data"`
	Fields      []FieldTemplate   `json:"fields,omitempty"` // Можно обновить поля
	Config      *CollectionConfig `json:"config,omitempty"` // Можно обновить конфиг
	IsActive    *bool             `json:"is_active,omitempty" example:"true"`
}
