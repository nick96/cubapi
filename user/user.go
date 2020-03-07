package user

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"go.uber.org/zap"
)

type UserRequest struct {
	Email     string `json:"email" validate:"email"`
	FirstName string `json:"firstName" validate:"required"`
	LastName  string `json:"lastName" validate:"required"`
	Password  string `json:"password" validate:"required,min=6"`
}

func (u UserRequest) Validate() error {
	validate := NewValidator()
	return validate.Struct(u)
}

type UserResponse User

func (u UserResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, http.StatusCreated)
	return nil
}

// NewUserRouter creates a router for the user endpoints.
//
// POST /: Create a new user
// GET /{userID}: Get the user with the given ID (if the requesting user has access).
func NewUserRouter(logger *zap.Logger, store UserStorer) func(chi.Router) {
	service := UserService{store}
	return func(r chi.Router) {
		r.Post("/", newUser(logger, service))
		r.Get("/{userID:[0-9]+}", getUserByID(logger, service))
	}
}

func newUser(logger *zap.Logger, service UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			logger.Info("Failed to read request body", zap.Error(err))
			render.Render(w, r, ErrInternal(err))
			return
		}
		defer r.Body.Close()

		var request UserRequest
		if err = json.Unmarshal(body, &request); err != nil {
			logger.Info("Failed to unmarshal user request", zap.Error(err), zap.ByteString("body", body))
			resp := ErrMalformedRequest("User request body is invalid", err)
			render.Render(w, r, resp)
			return
		}

		if err := request.Validate(); err != nil {
			logger.Info("Invalid user request",
				zap.Error(err),
				zap.String("requestID", middleware.GetReqID(r.Context())),
			)
			errs := validationErrors(err)
			render.Render(w, r, ErrInvalidRequest("Invalid request body", errs))
		}

		logger.Debug(
			"Received user creation request",
			zap.String("email", request.Email),
		)

		user := User{
			Email:     request.Email,
			FirstName: request.FirstName,
			LastName:  request.LastName,
			Password:  request.Password,
		}
		createdUser, err := service.NewUser(user)
		if err != nil {
			logger.Error("Failed to create new user",
				zap.Error(err),
				zap.String("requestID", middleware.GetReqID(r.Context())),
				zap.String("email", user.Email),
			)
		}
		logger.Debug("Created new user", zap.String("email", createdUser.Email), zap.Any("userID", createdUser.Id))
		render.Render(w, r, UserResponse(createdUser))
		w.WriteHeader(http.StatusCreated)
	}
}

func getUserByID(logger *zap.Logger, service UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotImplemented)
	}
}
