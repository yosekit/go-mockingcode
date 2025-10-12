package model

import "time"

// Collection represents a data collection in project
type Collection struct {
	ID          int64            `json:"id" example:"1"`
	ProjectID   int64            `json:"project_id" example:"1"`
	Name        string           `json:"name" example:"users"`
	Description string           `json:"description" example:"User data collection"`
	Fields      []FieldTemplate  `json:"fields"` // Шаблоны полей для генерации
	Config      CollectionConfig `json:"config"` // Настройки генерации
	IsActive    bool             `json:"is_active" example:"true"`
	CreatedAt   time.Time        `json:"created_at"`
	UpdatedAt   time.Time        `json:"updated_at"`
}

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

// FieldTemplate represents field for generation
type FieldTemplate struct {
	Name     string   `json:"name" example:"email"`
	Type     string   `json:"type" example:"string" enums:"string,number,boolean,date"`
	Format   string   `json:"format,omitempty" example:"email"`
	Required bool     `json:"required" example:"true"`
	Unique   bool     `json:"unique" example:"false"`
	Min      *float64 `json:"min,omitempty" example:"0"`                   // для numbers
	Max      *float64 `json:"max,omitempty" example:"100"`                 // для numbers
	Options  []string `json:"options,omitempty" example:"active,inactive"` // для enum
}

// CollectionConfig настройки генерации данных
type CollectionConfig struct {
	Count int    `json:"count" example:"10"` // Количество генерируемых записей
	Seed  *int64 `json:"seed,omitempty"`     // Seed для воспроизводимости
}
