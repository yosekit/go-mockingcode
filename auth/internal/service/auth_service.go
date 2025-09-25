package service

import (
	"errors"
	"time"

	"github.com/go-mockingcode/auth/internal/model"
	"github.com/go-mockingcode/auth/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
    userRepo *repository.UserRepository
    jwtSecret string
}

func NewAuthService(userRepo *repository.UserRepository, jwtSecret string) *AuthService {
    return &AuthService{
        userRepo: userRepo,
        jwtSecret: jwtSecret,
    }
}


func (s *AuthService) Register(req *model.RegisterRequest) (*model.AuthResponse, error) {
    // Проверяем, существует ли пользователь
    existingUser, err := s.userRepo.FindUserByEmail(req.Email)
    if err != nil {
        return nil, err
    }
    if existingUser != nil {
        return nil, errors.New("user already exists")
    }

    // Хешируем пароль
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
    if err != nil {
        return nil, err
    }

    // Создаем пользователя
    user := &model.User{
        Email:     req.Email,
        Password:  string(hashedPassword),
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
    }

    err = s.userRepo.CreateUser(user)
    if err != nil {
        return nil, err
    }

    // Генерируем JWT токен
    token, err := s.generateJWT(user.ID, user.Email)
    if err != nil {
        return nil, err
    }

    return &model.AuthResponse{
        UserID: user.ID,
        Email:  user.Email,
        Token:  token,
    }, nil
}

func (s *AuthService) Login(req *model.LoginRequest) (*model.AuthResponse, error) {
    user, err := s.userRepo.FindUserByEmail(req.Email)
    if err != nil {
        return nil, err
    }
    if user == nil {
        return nil, errors.New("invalid credentials")
    }

    // Проверяем пароль
    err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
    if err != nil {
        return nil, errors.New("invalid credentials")
    }

    // Генерируем JWT токен
    token, err := s.generateJWT(user.ID, user.Email)
    if err != nil {
        return nil, err
    }

    return &model.AuthResponse{
        UserID: user.ID,
        Email:  user.Email,
        Token:  token,
    }, nil
}

// Заглушка для JWT - реализуем в следующем шаге
func (s *AuthService) generateJWT(userID int64, email string) (string, error) {
    // Временно возвращаем простой токен
    return "jwt-token-placeholder-" + email, nil
}