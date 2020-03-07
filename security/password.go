package security

import (
	"crypto/rand"
	"fmt"
	"io"

	"golang.org/x/crypto/bcrypt"
)

const (
	PasswordCost = bcrypt.DefaultCost
)

func HashPassword(password, salt string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password + salt), PasswordCost)
	return string(hashed), err
}

func newSalt() (string, error) {
	salt := make([]byte, 10)
	_, err := io.ReadFull(rand.Reader, salt)
	if err != nil {
		return "", fmt.Errorf("failed to generate salt: %w", err)
	}
	return string(salt), nil
}

func HashNewPassword(password string) (string, string, error) {
	salt, err := newSalt()
	if err != nil {
		return "", "", fmt.Errorf("failed to hash password: %w", err)
	}
	hashed, err := HashPassword(password, salt)
	return hashed, salt, err
}


