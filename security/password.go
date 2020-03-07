package security

import (
	"golang.org/x/crypto/bcrypt"
)

const (
	PasswordCost = bcrypt.DefaultCost
)
