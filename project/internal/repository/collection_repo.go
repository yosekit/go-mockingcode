package repository

import (
	"database/sql"
	"fmt"

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
            fields JSONB NOT NULL,
            is_active BOOLEAN DEFAULT true,
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
            updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
            UNIQUE(project_id, name)
        );

		CREATE INDEX IF NOT EXISTS idx_collections_project_id ON collections(project_id);
	`

	_, err := r.db.Exec(query)
	return err
}

func (r *CollectionRepository) CreateCollection(collection *model.Collection) error {
	query := `
        INSERT INTO collections (project_id, name, description, fields, is_active, created_at, updated_at) 
        VALUES ($1, $2, $3, $4, $5, $6, $7) 
        RETURNING id`

	err := r.db.QueryRow(
		query,
		collection.ProjectID,
		collection.Name,
		collection.Description,
		collection.Fields,
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
		SELECT id, project_id, name, description, fields, is_active, created_at, updated_at 
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
		err := rows.Scan(
			&collection.ID,
			&collection.ProjectID,
			&collection.Name,
			&collection.Description,
			&collection.Fields,
			&collection.IsActive,
			&collection.CreatedAt,
			&collection.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		collections = append(collections, collection)
	}

	return collections, nil
}
