package db

const (

	// Запрос для создания нового пользователя.
	insertUser = "INSERT INTO users (first_name ,middle_name, last_name, login, password) VALUES ($1, $2, $3, $4, $5)"
	// Запрос для проверки логина и пароля (авторизация).
	checkUserCredentials = "SELECT id, FROM users WHERE login = $1 AND password = $2"
)
