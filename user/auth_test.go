package user

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/nick96/cubapi/security"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

func TestSignInNotExistantUser(t *testing.T) {
	reqContent := AuthnRequest{
		Email:    "test@test.com",
		Password: "password",
	}
	content, _ := json.Marshal(reqContent)
	req := httptest.NewRequest("POST", "/auth", bytes.NewReader(content))
	w := httptest.NewRecorder()

	logger := zap.NewNop()
	store := newMockUserStore()
	service := AuthService{store}
	handler := signIn(logger, NewValidator(), service)

	handler(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusForbidden {
		t.Fatalf("Expected status code %d, got %d", http.StatusForbidden, resp.StatusCode)
	}
}

func TestSignInExistingUser(t *testing.T) {
	reqContent := AuthnRequest{
		Email:    "test@test.com",
		Password: "password",
	}
	content, _ := json.Marshal(reqContent)
	req := httptest.NewRequest("POST", "/auth", bytes.NewReader(content))
	w := httptest.NewRecorder()
	logger := zap.NewNop()
	store := newMockUserStore()
	service := AuthService{store}
	handler := signIn(logger, NewValidator(), service)

	hashedPw, _ := bcrypt.GenerateFromPassword([]byte("password"+"salt"), security.PasswordCost)

	usr := User{
		Email:     reqContent.Email,
		FirstName: "Bobby",
		LastName:  "Tables",
		Password:  string(hashedPw),
		Salt:      "salt",
	}
	store.AddUser(usr)
	handler(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}
	body, _ := ioutil.ReadAll(resp.Body)
	var authResponse AuthResponse
	err := json.Unmarshal(body, &authResponse)
	if err != nil {
		t.Fatal(err)
	}
	if authResponse.Token == "" {
		t.Fatal("Expected token to not be an empty string")
	}
}
