package models

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
