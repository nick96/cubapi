package user

import (
	"fmt"
	"os"
	"time"

	"github.com/nick96/cubapi/security"
	"golang.org/x/crypto/bcrypt"
)

const (
	passwordHashCost = bcrypt.DefaultCost
)

type AuthService struct {
	store UserStorer
}

func ErrUserNotFound(message string, err error) security.ClientError {
	return security.NewClientError(message, err)
}

// AuthenticateUser authenticates the user by the given email and password. If
// all is well, the User entity is returned. Otherwise an error is returned.
// This error is safe to return to the client.
func (s AuthService) AuthenticateUser(email, password string) (User, security.ClientError) {
	user, found, err := s.store.FindByEmail(email)
	if err != nil {
		return User{}, security.NewClientError(
			"failed to retrieve user by email",
			fmt.Errorf("failed to retrieve user by email: %w", err),
		)
	} else if !found {
		return User{}, ErrUserNotFound(
			"username or password is incorrect",
			fmt.Errorf("could not find user with email '%s'", email),
		)
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(password+user.Salt),
	)
	if err != nil {
		// Specifically return an empty user here to prevent the real user being
		// used if the error is not checked for some reason.
		return User{}, security.NewClientError(
			"username or password is incorrect",
			err,
		)
	}

	return user, nil
}

// GetToken gets a new JWT token for a given user. The token expires in 24 hours.
func (s AuthService) GetToken(user User) (string, security.ClientError) {
	var jwt security.JWT
	token, err := jwt.Subject(user.Email).
		Issuer(os.Getenv("JWT_ISSUER")).
		Audience(user.Email).
		ExpireIn(24 * time.Hour).
		SignedToken(os.Getenv("JWT_SECRET"))
	if err != nil {
		return "", security.NewClientError("Failed to create authentication token", err)
	}
	return token, nil
}
