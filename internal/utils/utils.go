package utils

import "golang.org/x/crypto/bcrypt"

func HashPass(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}
