package jwtmanager

import (
	"errors"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
)

type Manager struct {
	secretKey     string
	tokenDuration time.Duration
}

func New(secretKey string, tokenDuration time.Duration) *Manager {
	return &Manager{
		secretKey:     secretKey,
		tokenDuration: tokenDuration,
	}
}

func (j *Manager) CreateToken(userID int64) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(j.tokenDuration).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.secretKey))
}

func (j *Manager) ValidateToken(token string) (bool, error) {
	parceToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		return []byte(j.secretKey), nil
	})
	if err != nil || !parceToken.Valid {
		return false, errors.New("invalid token")
	}
	return true, nil
}
