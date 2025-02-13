package logic

import (
	"context"
	"fmt"

	"github.com/Mark-Grigorev/auth/internal/db"
	jwtmanager "github.com/Mark-Grigorev/auth/internal/jwt_manager"
	"github.com/Mark-Grigorev/auth/internal/model"
	"github.com/Mark-Grigorev/auth/internal/utils"
)

type Logic struct {
	cfg        *model.Config
	db         *db.DBClient
	jwtManager *jwtmanager.Manager
}

func New(cfg *model.Config, db *db.DBClient, jwtManager *jwtmanager.Manager) *Logic {
	return &Logic{
		cfg:        cfg,
		db:         db,
		jwtManager: jwtManager,
	}
}

func (l *Logic) Register(ctx context.Context, userData *model.UserRegistrationData) (int64, error) {
	logPrefix := "[Register]"
	hashPass, err := utils.HashPass(userData.Password)
	if err != nil {
		return 0, fmt.Errorf("%s hashpass error -  %w", logPrefix, err)
	}
	userData.Password = string(hashPass)
	return l.db.CreateUser(ctx, userData)
}

// Authorization является первичной авторизацией.
// Возвращает токен, ошибку.
// В успешных кейсах отдает JWT на фронт, фронт проставляет JWT в хэддер.
func (l *Logic) Authorization(ctx context.Context, login string, pass string) (string, error) {
	logPrefix := "[Authorization]"
	var err error

	hashPass, err := utils.HashPass(pass)
	if err != nil {
		return "", fmt.Errorf("%s hash pass err - %w", logPrefix, err)
	}
	userID, err := l.db.Authorisation(ctx, login, string(hashPass))
	if err != nil {
		return "", fmt.Errorf("%s db err - %w", logPrefix, err)
	}
	token, err := l.jwtManager.CreateToken(userID)
	if err != nil {
		return "", fmt.Errorf("%s create token err - %w", logPrefix, err)
	}
	return token, nil
}

func (l *Logic) ValidateToken(ctx context.Context, token string) (bool, error) {
	return l.jwtManager.ValidateToken(token)
}
