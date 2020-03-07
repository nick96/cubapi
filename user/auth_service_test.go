package user

import (
	"flag"
	"os"
	"strings"
	"testing"

	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

func TestMain(m *testing.M) {
	flag.Parse()
	if testing.Short() {
		withMockUserStorer()
	} else {
		withUserStore()
	}

	os.Exit(m.Run())
}

func TestAuthenticateUser(t *testing.T) {
	defer cleanStore()
	email := "test@test.com"
	salt := "salt"
	password := "password"
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(password+salt), passwordHashCost)
	if err != nil {
		t.Fatal(err)
	}
	user := User{
		Email:     email,
		FirstName: "firstName",
		LastName:  "lastName",
		Password:  string(hashPassword),
		Salt:      salt,
	}
	id, err := store.AddUser(user)
	if err != nil {
		t.Fatal(err)
	}
	user.Id = id

	testCases := []struct {
		name           string
		email          string
		password       string
		expectedErrMsg string
		expectedUser   User
	}{
		{
			name:           "email + password ok",
			email:          email,
			password:       password,
			expectedErrMsg: "",
			expectedUser:   user,
		},
		{
			name:           "email ok",
			email:          email,
			password:       "wrong",
			expectedErrMsg: "username or password is incorrect",
			expectedUser:   User{},
		},
		{
			name:           "password ok",
			email:          "wrong@test.com",
			password:       password,
			expectedErrMsg: "username or password is incorrect",
			expectedUser:   User{},
		},
		{
			name:           "both wrong",
			email:          "wrong@test.com",
			password:       "wrong",
			expectedErrMsg: "username or password is incorrect",
			expectedUser:   User{},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			retUser, err := authService.AuthenticateUser(tt.email, tt.password)
			if tt.expectedErrMsg != "" {
				if tt.expectedErrMsg != err.SafeError() {
					t.Errorf("Expected error message %s, got %s: %v", tt.expectedErrMsg, err.SafeError(), err)
				}
			} else if err != nil {
				t.Fatalf("Expected no error but got: %v", err)
			}

			if tt.expectedUser != retUser {
				t.Errorf("Expected user %v, got %v", tt.expectedUser, retUser)
			}
		})
	}

}

func TestGetToken(t *testing.T) {
	oldJwtSecret := os.Getenv("JWT_SECRET")
	os.Setenv("JWT_SECRET", "secret")
	defer func() {
		os.Setenv("JWT_SECRET", oldJwtSecret)
	}()

	user := User{
		Email:     "test@test.com",
		FirstName: "test",
		LastName:  "test",
		Password:  "testpassword",
		Salt:      "password",
	}

	token, err := authService.GetToken(user)
	if err != nil {
		t.Fatal(err)
	}
	if token == "" {
		t.Fatal("Expected token not to be empty")
	}

	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		t.Fatalf("Expected token to have 3 parts, it has %d", len(parts))
	}
}
