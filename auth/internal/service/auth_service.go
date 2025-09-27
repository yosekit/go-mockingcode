package service

import (
	"errors"
	"time"

	"github.com/go-mockingcode/auth/internal/config"
	"github.com/go-mockingcode/auth/internal/model"
	"github.com/go-mockingcode/auth/internal/repository"
	"github.com/go-mockingcode/auth/internal/utils"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
    userRepo *repository.UserRepository
    tokenRepo *repository.TokenRepository
    jwtSecret string
    accessTokenExpiry  time.Duration
    refreshTokenExpiry time.Duration
}

func NewAuthService(userRepo *repository.UserRepository, tokenRepo *repository.TokenRepository, 
    cfg *config.Config) *AuthService {
    return &AuthService{
        userRepo: userRepo,
        tokenRepo: tokenRepo,
        jwtSecret: cfg.JWTSecret,
        accessTokenExpiry: cfg.AccessTokenExpiry,
        refreshTokenExpiry: cfg.RefreshTokenExpiry,
    }
}

// Register - только создание пользователя
func (s *AuthService) Register(req *model.RegisterRequest) error {
    // Проверяем, существует ли пользователь
    existingUser, err := s.userRepo.FindUserByEmail(req.Email)
    if err != nil {
        return err
    }
    if existingUser != nil {
        return errors.New("user already exists")
    }

    // Хешируем пароль
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
    if err != nil {
        return err
    }

    // Создаем пользователя
    user := &model.User{
        Email:     req.Email,
        Password:  string(hashedPassword),
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
    }

    return s.userRepo.CreateUser(user)
}

// Login - аутентификация + выдача токенов
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

    // Генерируем новые токены
    return s.generateTokens(user.ID, user.Email)
}

// AutoLogin - после успешной регистрации
func (s *AuthService) AutoLogin(email string) (*model.AuthResponse, error) {
    user, err := s.userRepo.FindUserByEmail(email)
    if err != nil {
        return nil, err
    }
    if user == nil {
        return nil, errors.New("user not found")
    }

    // Генерируем новые токены
    return s.generateTokens(user.ID, user.Email)
} 

func (s *AuthService) RefreshTokens(refreshToken string) (*model.AuthResponse, error) {
    // Находим refresh token в БД
    token, err := s.tokenRepo.FindRefreshToken(refreshToken)
    if err != nil {
        return nil, err
    }
    if token == nil {
        return nil, errors.New("invalid refresh token")
    }

    // Находим пользователя
    user, err := s.userRepo.FindUserByID(token.UserID)
    if err != nil {
        return nil, err
    }
    if user == nil {
        return nil, errors.New("user not found")
    }

    // Генерируем новые токены
    accessToken, expiresAt, err := s.generateAccessToken(user.ID, user.Email)
    if err != nil {
        return nil, err
    }

    // Удаляем старый refresh token (one-time use)
    if err := s.tokenRepo.DeleteRefreshToken(refreshToken); err != nil {
        return nil, err
    }

    // Генерируем новый refresh token
    newRefreshToken, err := s.generateAndStoreRefreshToken(user.ID)
    if err != nil {
        return nil, err
    }

    return &model.AuthResponse{
        UserID:       user.ID,
        Email:        user.Email,
        AccessToken:  accessToken,
        RefreshToken: newRefreshToken,
        ExpiresAt:    expiresAt,
    }, nil
}

func (s *AuthService) ValidateAccessToken(tokenString string) (*model.User, error) {
    claims, err := utils.ValidateAccessToken(tokenString, s.jwtSecret)
    if err != nil {
        return nil, err
    }

    user, err := s.userRepo.FindUserByEmail(claims.Email)
    if err != nil {
        return nil, err
    }
    if user == nil {
        return nil, errors.New("user not found")
    }

    return user, nil
}

func (s *AuthService) Logout(refreshToken string) error {
    return s.tokenRepo.DeleteRefreshToken(refreshToken)
}

func (s *AuthService) generateTokens(userID int64, email string) (*model.AuthResponse, error) {
    accessToken, expiresAt, err := s.generateAccessToken(userID, email)
    if err != nil {
        return nil, err
    }

    refreshToken, err := s.generateAndStoreRefreshToken(userID)
    if err != nil {
        return nil, err
    }

    return &model.AuthResponse{
        UserID: userID,
        Email:  email,
        AccessToken:  accessToken,
        RefreshToken: refreshToken,
        ExpiresAt:    expiresAt,
    }, nil
}

func (s *AuthService) generateAccessToken(userID int64, email string) (string, int64, error) {
    return utils.GenerateAccessToken(userID, email, s.jwtSecret, s.accessTokenExpiry);
}

func (s *AuthService) generateAndStoreRefreshToken(userID int64) (string, error) {
    refreshToken, err := utils.GenerateRefreshToken()
    if err != nil {
        return "", err
    }

    token := &model.RefreshToken{
        UserID:    userID,
        Token:     refreshToken,
        ExpiresAt: time.Now().Add(s.refreshTokenExpiry),
    }

    if err := s.tokenRepo.CreateRefreshToken(token); err != nil {
        return "", err
    }

    return refreshToken, nil
}

