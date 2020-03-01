package user

import (
	"net/http"

	"github.com/go-chi/chi"
	"go.uber.org/zap"
)

func NewUserRouter(logger *zap.Logger, store UserStorer) func(chi.Router) {
	return func(r chi.Router) {
		r.Post("/user", newUser(logger, store))
		r.Get("/user/{userID:[0-9]+}", getUserByID(logger, store))
	}
}

func newUser(logger *zap.Logger, store UserStorer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotImplemented)
		w.Write([]byte("new user"))
	}
}

func getUserByID(logger *zap.Logger, store UserStorer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotImplemented)
	}
}
