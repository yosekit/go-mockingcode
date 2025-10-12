package repository

import (
	"database/sql"
	"fmt"

	"github.com/go-mockingcode/project/internal/model"
)

type ProjectRepository struct {
	db *sql.DB
}

func NewProjectRepository(db *sql.DB) *ProjectRepository {
	return &ProjectRepository{db: db}
}

func (r *ProjectRepository) InitSchema() error {
	query := `
        CREATE TABLE IF NOT EXISTS projects (
            id SERIAL PRIMARY KEY,
            user_id INTEGER NOT NULL,
            name VARCHAR(100) NOT NULL,
            description TEXT,
            api_key VARCHAR(16) UNIQUE NOT NULL,
            base_url VARCHAR(255),
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
            updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
        );

        CREATE INDEX IF NOT EXISTS idx_projects_user_id ON projects(user_id);
        CREATE INDEX IF NOT EXISTS idx_projects_api_key ON projects(api_key);
    `

	_, err := r.db.Exec(query)
	return err
}

func (r *ProjectRepository) CreateProject(project *model.Project) error {
	query := `
        INSERT INTO projects (user_id, name, description, api_key, base_url, created_at, updated_at) 
        VALUES ($1, $2, $3, $4, $5, $6, $7) 
        RETURNING id`

	err := r.db.QueryRow(
		query,
		project.UserID,
		project.Name,
		project.Description,
		project.APIKey,
		project.BaseURL,
		project.CreatedAt,
		project.UpdatedAt,
	).Scan(&project.ID)

	if err != nil {
		return fmt.Errorf("failed to create project: %v", err)
	}

	return nil
}

func (r *ProjectRepository) GetUserProjects(userID int64) ([]*model.Project, error) {
	query := `
		SELECT id, user_id, name, description, api_key, base_url, created_at, updated_at 
        FROM projects 
		WHERE user_id = $1 
		ORDER BY created_at DESC`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []*model.Project
	for rows.Next() {
		project := &model.Project{}
		err := rows.Scan(
			&project.ID,
			&project.UserID,
			&project.Name,
			&project.Description,
			&project.APIKey,
			&project.BaseURL,
			&project.CreatedAt,
			&project.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		projects = append(projects, project)
	}

	return projects, nil
}

func (r *ProjectRepository) GetProjectByID(projectID int64, userID int64) (*model.Project, error) {
	query := `
		SELECT id, user_id, name, description, api_key, base_url, created_at, updated_at 
		FROM projects 
		WHERE id = $1 
			AND user_id = $2`

	project := &model.Project{}
	err := r.db.QueryRow(query, projectID, userID).Scan(
		&project.ID,
		&project.UserID,
		&project.Name,
		&project.Description,
		&project.APIKey,
		&project.BaseURL,
		&project.CreatedAt,
		&project.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return project, nil
}

func (r *ProjectRepository) GetProjectByAPIKey(apiKey string) (*model.Project, error) {
	query := `
		SELECT id, user_id, name, description, api_key, base_url, created_at, updated_at 
        FROM projects 
		WHERE api_key = $1`

	project := &model.Project{}
	err := r.db.QueryRow(query, apiKey).Scan(
		&project.ID,
		&project.UserID,
		&project.Name,
		&project.Description,
		&project.APIKey,
		&project.BaseURL,
		&project.CreatedAt,
		&project.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return project, nil
}

// UpdateProject обновляет проект
func (r *ProjectRepository) UpdateProject(project *model.Project) error {
	query := `UPDATE projects SET name = $1, description = $2, updated_at = $3 WHERE id = $4`
	_, err := r.db.Exec(query, project.Name, project.Description, project.UpdatedAt, project.ID)
	return err
}

// DeleteProject удаляет проект
func (r *ProjectRepository) DeleteProject(projectID int64) error {
	query := `DELETE FROM projects WHERE id = $1`
	_, err := r.db.Exec(query, projectID)
	return err
}
