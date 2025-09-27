package utils

import (
	"crypto/rand"
	"encoding/hex"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type AccessClaims struct {
    UserID int64  `json:"user_id"`
    Email  string `json:"email"`
    jwt.RegisteredClaims
}

func GenerateAccessToken(userID int64, email string, secret string, 
    expiry time.Duration) (string, int64, error) {
    expirationTime := time.Now().Add(expiry)
    expiresAt := expirationTime.Unix()

    claims := &AccessClaims{
        UserID: userID,
        Email:  email,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(expirationTime),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
            Issuer:    "mockingcode-auth",
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    tokenString, err := token.SignedString([]byte(secret))
    return tokenString, expiresAt, err
}

// Refresh Token - случайная строка (не JWT)
func GenerateRefreshToken() (string, error) {
    bytes := make([]byte, 32) // 256 бит
    if _, err := rand.Read(bytes); err != nil {
        return "", err
    }
    return hex.EncodeToString(bytes), nil
}

func ValidateAccessToken(tokenString string, secret string) (*AccessClaims, error) {
    claims := &AccessClaims{}

    token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
        return []byte(secret), nil
    })

    if err != nil {
        return nil, err
    }

    if !token.Valid {
        return nil, jwt.ErrSignatureInvalid
    }

    return claims, nil
}