package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/go-mockingcode/auth/internal/model"
)

type TokenRepository struct {
    db *sql.DB
}

func NewTokenRepository(db *sql.DB) *TokenRepository {
    return &TokenRepository{db: db}
}


func (r *TokenRepository) CreateRefreshToken(token *model.RefreshToken) error {
    query := `
        INSERT INTO refresh_tokens (user_id, token, expires_at, created_at) 
        VALUES ($1, $2, $3, $4) 
        RETURNING id`

    err := r.db.QueryRow(
        query,
        token.UserID,
        token.Token,
        token.ExpiresAt,
        time.Now(),
    ).Scan(&token.ID)

    if err != nil {
        return fmt.Errorf("failed to create refresh token: %v", err)
    }

    return nil
}

func (r *TokenRepository) FindRefreshToken(tokenString string) (*model.RefreshToken, error) {
    query := `
        SELECT id, user_id, token, expires_at, created_at 
        FROM refresh_tokens 
        WHERE token = $1 AND expires_at > $2`

    token := &model.RefreshToken{}
    err := r.db.QueryRow(query, tokenString, time.Now()).Scan(
        &token.ID,
        &token.UserID,
        &token.Token,
        &token.ExpiresAt,
        &token.CreatedAt,
    )

    if err == sql.ErrNoRows {
        return nil, nil
    }
    if err != nil {
        return nil, fmt.Errorf("failed to find refresh token: %v", err)
    }

    return token, nil
}

func (r *TokenRepository) DeleteRefreshToken(tokenString string) error {
    query := `DELETE FROM refresh_tokens WHERE token = $1`
    _, err := r.db.Exec(query, tokenString)
    return err
}

func (r *TokenRepository) DeleteExpiredTokens() error {
    query := `DELETE FROM refresh_tokens WHERE expires_at <= $1`
    _, err := r.db.Exec(query, time.Now())
    return err
}


func (r *TokenRepository) InitSchema() error {
    query := `
        CREATE TABLE IF NOT EXISTS refresh_tokens (
            id SERIAL PRIMARY KEY,
            user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
            token VARCHAR(512) UNIQUE NOT NULL,
            expires_at TIMESTAMP NOT NULL,
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
        );
        
        CREATE INDEX IF NOT EXISTS idx_refresh_tokens_token ON refresh_tokens(token);
        CREATE INDEX IF NOT EXISTS idx_refresh_tokens_expires ON refresh_tokens(expires_at);
    `

    _, err := r.db.Exec(query)
    return err
}