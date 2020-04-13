package user

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/nick96/cubapi/security"
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

func ErrUserAlreadyExists(message string) render.Renderer {
	return &ErrorResponse{
		Message: message,
		Status:  http.StatusBadRequest,
	}
}

// NewUserRouter creates a router for the user endpoints.
//
// POST /: Create a new user
// GET /{userID}: Get the user with the given ID (if the requesting user has access).
func NewUserRouter(logger *zap.Logger, store UserStorer) func(chi.Router) {
	service := UserService{store}
	return func(r chi.Router) {
		r.Post("/", newUser(logger, service))
		r.Get("/me", getAuthdUser(logger, service))
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
			return
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
		if IsErrUserAlreadyExists(err) {
			logger.Error(
				"Failed to create new user as they already exist",
				zap.Error(err),
				zap.String("requestID", middleware.GetReqID(r.Context())),
				zap.String("email", user.Email),
			)
			render.Render(w, r, ErrUserAlreadyExists(fmt.Sprintf("User with email %s already exists", user.Email)))
			return
		} else if err != nil {
			logger.Error("Failed to create new user",
				zap.Error(err),
				zap.String("requestID", middleware.GetReqID(r.Context())),
				zap.String("email", user.Email),
			)
			render.Render(w, r, ErrInternalWithMessage(fmt.Sprintf("Failed to create user %s", user.Email), err))
			return
		}
		logger.Debug("Created new user", zap.String("email", createdUser.Email), zap.Any("userID", createdUser.Id))
		render.Render(w, r, UserResponse(createdUser))
		w.WriteHeader(http.StatusCreated)
	}
}

func getAuthdUser(logger *zap.Logger, service UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var jwt string
		jwtCookie, err := r.Cookie("jwt")
		if err == nil {
			// Get the JWT from the cookie
			jwt = jwtCookie.String()
			logger.Debug("Retrieved JWT from cookie")
		} else {
			// Get the cookie from the Authorization header
			jwt = r.Header.Get("Authorization")
			jwt = strings.Replace(jwt, "Bearer ", "", 1)
			logger.Debug("Retrieved JWT from authorization header")
		}

		if jwt == "" {
			logger.Info("No JWT was provided")
			render.Render(w, r, ErrForbidden("'jwt' cookie or Authorization header with bearer token is required", nil))
			return
		}

		secret := os.Getenv("JWT_SECRET")
		token, err := security.ValidateToken(jwt, secret)
		if err != nil {
			logger.Info("JWT validation failed", zap.Error(err), zap.String("invalidJWT", jwt))
			render.Render(w, r, ErrForbidden("JWT validation failed", security.NewClientError(err.Error(), err)))
			return
		}

		user, isFound, err := service.store.FindByEmail(token.Email)
		if err != nil {
			logger.Error("Failed to retrieve user from database", zap.String("email", token.Email), zap.Error(err))
			render.Render(w, r, ErrInternal(err))
			return
		}

		if !isFound {
			logger.Info("Could not find user", zap.String("email", token.Email))
			render.Render(w, r, ErrForbidden(fmt.Sprintf("Could not find user with email %s", token.Email), nil))
			return
		}
		response := UserResponse(user)
		logger.Debug("Successfully validated JWT, responding with user details", zap.Any("user", response))
		render.Render(w, r, response)
	}
}
