package model

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

// MockDocument представляет документ в MongoDB
type MockDocument struct {
	ID             bson.ObjectID  `bson:"_id,omitempty" json:"id"`
	ProjectID      int64          `bson:"project_id" json:"project_id"`
	CollectionName string         `bson:"collection_name" json:"collection_name"`
	Data           map[string]any `bson:"data" json:"data"` // Динамические данные
	CreatedAt      time.Time      `bson:"created_at" json:"created_at"`
	UpdatedAt      time.Time      `bson:"updated_at" json:"updated_at"`
}

// DocumentResponse ответ с документом
type DocumentResponse struct {
	Document *MockDocument `json:"document"`
}

// DocumentsResponse ответ с документами
type DocumentsResponse struct {
	Documents []*MockDocument `json:"documents"`
	Total     int64           `json:"total"`
	Limit     int64           `json:"limit"`
	Offset    int64           `json:"offset"`
}

// GenerateRequest запрос на генерацию данных
type GenerateRequest struct {
	Count *int    `json:"count,omitempty" example:"50"`
	Seed  *uint64 `json:"seed,omitempty" example:"12345"`
}

// QueryOptions опции для запросов
type QueryOptions struct {
	Limit  *int64 `json:"limit,omitempty" example:"10"`
	Offset *int64 `json:"offset,omitempty" example:"0"`
	Sort   string `json:"sort,omitempty" example:"created_at"`
	Order  string `json:"order,omitempty" example:"desc" enums:"asc,desc"`
}
