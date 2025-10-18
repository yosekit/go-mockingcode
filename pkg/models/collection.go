package models

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
