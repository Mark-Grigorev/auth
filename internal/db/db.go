package db

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Mark-Grigorev/auth/internal/model"
)

type DBClient struct {
	db *sql.DB
}

func New(dbConnection string) (*DBClient, error) {
	// Открываем соединение с базой данных
	db, err := sql.Open("postgres", dbConnection)
	if err != nil {
		return &DBClient{}, err
	}

	// Проверяем подключение к базе данных
	err = db.Ping()
	if err != nil {
		db.Close()
		return &DBClient{}, err
	}

	return &DBClient{
		db: db,
	}, nil
}

func (c *DBClient) CreateUser(ctx context.Context, userData *model.UserRegistrationData) (int64, error) {
	_, err := c.db.ExecContext(
		ctx,
		insertUser,
		userData.FirstName,
		userData.MiddleName,
		userData.LastName,
		userData.Login,
		userData.Password)
	if err != nil {
		return 0, err
	}
	return 0, nil
}

func (c *DBClient) Authorisation(ctx context.Context, login, password string) (int64, error) {
	var userID int64

	// Выполняем запрос с переданными параметрами логина и пароля
	err := c.db.QueryRow(checkUserCredentials, login, password).Scan(&userID)
	if err != nil {
		return 0, fmt.Errorf("пользователь с таким логином и паролем не найден - %s", err)
	}

	return userID, nil
}
