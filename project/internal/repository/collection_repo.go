package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-mockingcode/project/internal/model"
)

type CollectionRepository struct {
	db *sql.DB
}

func NewCollectionRepository(db *sql.DB) *CollectionRepository {
	return &CollectionRepository{db: db}
}

func (r *CollectionRepository) InitSchema() error {
	query := `
 		CREATE TABLE IF NOT EXISTS collections (
            id SERIAL PRIMARY KEY,
            project_id INTEGER NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
            name VARCHAR(50) NOT NULL,
            description TEXT,
            fields JSONB,
            config JSONB DEFAULT '{"count": 10}',
            is_active BOOLEAN DEFAULT true,
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
            updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
            UNIQUE(project_id, name)
        );

		CREATE INDEX IF NOT EXISTS idx_collections_project_id ON collections(project_id);`

	_, err := r.db.Exec(query)
	return err
}

func (r *CollectionRepository) CreateCollection(collection *model.Collection) error {
	query := `
        INSERT INTO collections (project_id, name, description, fields, config, is_active, created_at, updated_at) 
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8) 
        RETURNING id`

	// Конвертируем в JSON
	fieldsJSON, _ := json.Marshal(collection.Fields)
	configJSON, _ := json.Marshal(collection.Config)

	err := r.db.QueryRow(
		query,
		collection.ProjectID,
		collection.Name,
		collection.Description,
		fieldsJSON,
		configJSON,
		collection.IsActive,
		collection.CreatedAt,
		collection.UpdatedAt,
	).Scan(&collection.ID)

	if err != nil {
		return fmt.Errorf("failed to create collection: %v", err)
	}

	return nil
}

func (r *CollectionRepository) GetProjectCollections(projectID int64) ([]*model.Collection, error) {
	query := `
        SELECT id, project_id, name, description, fields, config, is_active, created_at, updated_at 
        FROM collections 
		WHERE project_id = $1 
		ORDER BY created_at DESC`

	rows, err := r.db.Query(query, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var collections []*model.Collection
	for rows.Next() {
		collection := &model.Collection{}
		var fieldsJSON, configJSON []byte

		err := rows.Scan(
			&collection.ID,
			&collection.ProjectID,
			&collection.Name,
			&collection.Description,
			&fieldsJSON,
			&configJSON,
			&collection.IsActive,
			&collection.CreatedAt,
			&collection.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		// Парсим JSON поля
		json.Unmarshal(fieldsJSON, &collection.Fields)
		json.Unmarshal(configJSON, &collection.Config)

		collections = append(collections, collection)
	}

	return collections, nil
}

// GetCollectionByName возвращает коллекцию по имени и project_id
func (r *CollectionRepository) GetCollectionByName(projectID int64, collectionName string) (*model.Collection, error) {
	query := `
        SELECT id, project_id, name, description, fields, config, is_active, created_at, updated_at 
        FROM collections 
		WHERE project_id = $1 AND name = $2`

	collection := &model.Collection{}
	var fieldsJSON, configJSON []byte

	err := r.db.QueryRow(query, projectID, collectionName).Scan(
		&collection.ID,
		&collection.ProjectID,
		&collection.Name,
		&collection.Description,
		&fieldsJSON,
		&configJSON,
		&collection.IsActive,
		&collection.CreatedAt,
		&collection.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil // Коллекция не найдена (это нормально для гибридного подхода)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get collection by name: %w", err)
	}

	// Десериализуем JSON
	if err := json.Unmarshal(fieldsJSON, &collection.Fields); err != nil {
		return nil, fmt.Errorf("failed to unmarshal fields: %w", err)
	}
	if err := json.Unmarshal(configJSON, &collection.Config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return collection, nil
}

// GetCollectionByID возвращает коллекцию по ID
func (r *CollectionRepository) GetCollectionByID(collectionID int64, projectID int64) (*model.Collection, error) {
	query := `
        SELECT id, project_id, name, description, fields, config, is_active, created_at, updated_at 
        FROM collections 
		WHERE id = $1 AND project_id = $2`

	collection := &model.Collection{}
	var fieldsJSON, configJSON []byte

	err := r.db.QueryRow(query, collectionID, projectID).Scan(
		&collection.ID,
		&collection.ProjectID,
		&collection.Name,
		&collection.Description,
		&fieldsJSON,
		&configJSON,
		&collection.IsActive,
		&collection.CreatedAt,
		&collection.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find collection: %v", err)
	}

	// Парсим JSON поля
	json.Unmarshal(fieldsJSON, &collection.Fields)
	json.Unmarshal(configJSON, &collection.Config)

	return collection, nil
}

// UpdateCollection обновляет коллекцию
func (r *CollectionRepository) UpdateCollection(collection *model.Collection) error {
	query := `
        UPDATE collections 
        SET name = $1, description = $2, fields = $3, config = $4, is_active = $5, updated_at = $6 
        WHERE id = $7 AND project_id = $8`

	fieldsJSON, _ := json.Marshal(collection.Fields)
	configJSON, _ := json.Marshal(collection.Config)

	_, err := r.db.Exec(
		query,
		collection.Name,
		collection.Description,
		fieldsJSON,
		configJSON,
		collection.IsActive,
		time.Now(),
		collection.ID,
		collection.ProjectID,
	)
	return err
}

// DeleteCollection удаляет коллекцию
func (r *CollectionRepository) DeleteCollection(collectionID int64, projectID int64) error {
	query := `DELETE FROM collections WHERE id = $1 AND project_id = $2`
	result, err := r.db.Exec(query, collectionID, projectID)
	if err != nil {
		return fmt.Errorf("failed to delete collection: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("collection not found")
	}

	return nil
}
