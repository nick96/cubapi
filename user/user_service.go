package user

import (
	"github.com/nick96/cubapi/security"
)

type UserService struct {
	store UserStorer
}

func (s UserService) hashPassword(password, salt string) (string, error) {
	return security.HashPassword(password, salt)
}

func (s UserService) NewUser(user User) (User, security.ClientError) {
	saltedPassword, salt, err := security.HashNewPassword(user.Password)
	if err != nil {
		return User{}, security.NewClientError("failed to create new user", err)
	}
	user.Password = saltedPassword
	user.Salt = salt
	id, err := s.store.AddUser(user)
	if err != nil {
		return User{}, security.NewClientError("failed to create new user", err)
	}
	user.Id = id
	return user, nil
}
