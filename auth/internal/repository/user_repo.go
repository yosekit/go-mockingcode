package repository

import (
	"database/sql"
	"fmt"

	"github.com/go-mockingcode/auth/internal/model"
)

type UserRepository struct {
    db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
    return &UserRepository{db: db}
}

func (r *UserRepository) CreateUser(user *model.User) error {
    query := `
        INSERT INTO users (email, password, created_at, updated_at) 
        VALUES ($1, $2, $3, $4) 
        RETURNING id`

    err := r.db.QueryRow(
        query,
        user.Email,
        user.Password,
        user.CreatedAt,
        user.UpdatedAt,
    ).Scan(&user.ID)

    if err != nil {
        return fmt.Errorf("failed to create user: %v", err)
    }

    return nil
}

func (r *UserRepository) FindUserByEmail(email string) (*model.User, error) {
    query := `SELECT id, email, password, created_at, updated_at FROM users WHERE email = $1`
    
    user := &model.User{}
    err := r.db.QueryRow(query, email).Scan(
        &user.ID,
        &user.Email,
        &user.Password,
        &user.CreatedAt,
        &user.UpdatedAt,
    )

    if err == sql.ErrNoRows {
        return nil, nil
    }
    if err != nil {
        return nil, fmt.Errorf("failed to find user: %v", err)
    }

    return user, nil
}

func (r *UserRepository) InitSchema() error {
    query := `
        CREATE TABLE IF NOT EXISTS users (
            id SERIAL PRIMARY KEY,
            email VARCHAR(255) UNIQUE NOT NULL,
            password VARCHAR(255) NOT NULL,
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
            updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
        );
        
        CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
    `

    _, err := r.db.Exec(query)
    if err != nil {
        return fmt.Errorf("failed to init schema: %v", err)
    }

    return nil
}