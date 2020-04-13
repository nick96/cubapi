package security

import (
	"golang.org/x/crypto/bcrypt"
)

const (
	PasswordCost = bcrypt.DefaultCost
)

func HashPassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), PasswordCost)
	return string(hashed), err
}

func HashNewPassword(password string) (string, error) {
	hashed, err := HashPassword(password)
	return hashed, err
}
