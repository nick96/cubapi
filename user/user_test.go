package user

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"go.uber.org/zap"
)

func TestNewUserOk(t *testing.T) {
	store := newMockUserStore()
	service := UserService{store}
	logger := zap.NewNop()
	handler := newUser(logger, service)

	requestBody := UserRequest{
		Email:     "test@test.com",
		FirstName: "firstName",
		LastName:  "lastName",
		Password:  "password",
	}
	content, err := json.Marshal(requestBody)
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/", bytes.NewReader(content))
	r.Header.Add("Content-Type", "application/json")

	handler(w, r)

	if w.Result().StatusCode != http.StatusCreated {
		t.Fatalf("Expected status code %d, got %d: %s", http.StatusCreated, w.Result().StatusCode, string(content))
	}

	user, exists, _ := store.FindByEmail(requestBody.Email)
	if !exists {
		t.Fatalf("Expected user with email %s to have been created", requestBody.Email)
	}

	if user.Email != requestBody.Email {
		t.Fatalf("Expected email %s, got %s", requestBody.Email, user.Email)
	}

	if user.FirstName != requestBody.FirstName {
		t.Fatalf("Expected first name %s, got %s", requestBody.FirstName, user.FirstName)
	}

	if user.LastName != requestBody.LastName {
		t.Fatalf("Expected last name %s, got %s", requestBody.LastName, user.LastName)
	}
}

func TestNewUserInvalidRequest(t *testing.T) {
	testCases := []struct{
		name string
		requestBody UserRequest
	}{
		{
			name: "no-email",
			requestBody: UserRequest{
				FirstName: "firstName",
				LastName:  "lastName",
				Password:  "password",
			},
		},
		{
			name: "no-first-name",
			requestBody: UserRequest{
				Email:     "test@test.com",
				LastName:  "lastName",
				Password:  "password",
			},
		},
		{
			name: "no-last-name",
			requestBody: UserRequest{
				Email:     "test@test.com",
				FirstName: "firstName",
				Password:  "password",
			},
		},
		{
			name: "no-password",
			requestBody: UserRequest{
				Email:     "test@test.com",
				FirstName: "firstName",
				LastName:  "lastName",
			},
		},
		{
			name: "password-too-short",
			requestBody: UserRequest{
				Email:     "test@test.com",
				FirstName: "firstName",
				LastName:  "lastName",
				Password:  "pass",
			},
		},
		{
			name: "invalid-email",
			requestBody: UserRequest{
				Email:     "test",
				FirstName: "firstName",
				LastName:  "lastName",
				Password:  "password",
			},
		},
	}
	store := newMockUserStore()
	service := UserService{store}
	logger := zap.NewNop()
	handler := newUser(logger, service)

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.requestBody)
			r := httptest.NewRequest("POST", "/", bytes.NewReader(body))
			w := httptest.NewRecorder()
			handler(w, r)
			if w.Result().StatusCode != http.StatusBadRequest {
				t.Fatalf("Expected status code %d, got %d", http.StatusBadRequest, w.Result().StatusCode)
			}
		})
	}
}
