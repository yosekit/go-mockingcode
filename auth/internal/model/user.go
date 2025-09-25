package model

import "time"

type User struct {
    ID        int64     `json:"id"`
    Email     string    `json:"email"`
    Password  string    `json:"-"` // Пароль не отдаем в JSON
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

type RegisterRequest struct {
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required,min=6"`
}

type LoginRequest struct {
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required"`
}

type AuthResponse struct {
    UserID int64  `json:"user_id"`
    Email  string `json:"email"`
    Token  string `json:"token"`
}