package user

import (
	"fmt"

	"github.com/nick96/cubapi/security"
)

type UserService struct {
	store UserStorer
}

func (s UserService) hashPassword(password string) (string, error) {
	return security.HashPassword(password)
}

type errUserAlreadyExists struct {
	email string
}

func (e errUserAlreadyExists) Error() string {
	return fmt.Sprintf("email %s is already in user", e.email)
}

func (e errUserAlreadyExists) SafeError() string {
	return e.Error()
}

func (s UserService) NewUser(user User) (User, security.ClientError) {
	hashedPassword, err := security.HashNewPassword(user.Password)
	if err != nil {
		return User{}, security.NewClientError("failed to create new user", err)
	}
	if isAvailable, err := s.isEmailAvailable(user.Email); err != nil {
		return User{}, err
	} else if isAvailable {
		user.Password = hashedPassword
		id, err := s.store.AddUser(user)
		if err != nil {
			return User{}, security.NewClientError("failed to create new user", err)
		}
		user.Id = id
		return user, nil
	}
	return User{}, errUserAlreadyExists{email: user.Email}
}

func (s UserService) isEmailAvailable(email string) (bool, security.ClientError) {
	_, exists, err := s.store.FindByEmail(email)
	if err != nil {
		return false, security.NewClientError(
			fmt.Sprintf("failed to check if user with email %s exists", email),
			err,
		)
	}
	isAvailable := !exists
	return isAvailable, nil
}

func IsErrUserAlreadyExists(err error) bool {
	_, ok := err.(errUserAlreadyExists)
	return ok
}
