package user

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/nick96/cubapi/security"
	"go.uber.org/zap"
	"gopkg.in/go-playground/validator.v9"
)

type ErrorResponse struct {
	Status  int      `json:"-"`
	Message string   `json:"message"`
	Error   string   `json:"error,omitempty"`
	Errors  []string `json:"errors,omitempty"`
}

func (e *ErrorResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.Status)
	return nil
}

func FieldError(fe validator.FieldError) error {
	// TODO: Make this more informative
	var sb strings.Builder
	fmt.Fprintf(&sb, "'%s' is not valid a valid value for '%s'", fe.Value(), fe.Field())
	return fmt.Errorf(sb.String())
}

func ErrMalformedRequest(message string, err error) render.Renderer {
	return &ErrorResponse{
		Message: message,
		Status:  http.StatusBadRequest,
		Error:   err.Error(),
	}
}

func ErrInvalidRequest(message string, errs []error) render.Renderer {
	var errors []string
	for _, err := range errs {
		errors = append(errors, err.Error())
	}

	return &ErrorResponse{
		Message: message,
		Status:  http.StatusBadRequest,
		Errors:  errors,
	}
}

func ErrInternal(err error) render.Renderer {
	return ErrInternalWithMessage("Internal error", err)
}

func ErrInternalWithMessage(message string, err error) render.Renderer {
	return &ErrorResponse{
		Message: message,
		Status:  http.StatusInternalServerError,
	}
}

func ErrForbidden(message string, err security.ClientError) render.Renderer {
	var safeError string
	if err != nil {
		safeError = err.SafeError()
	}
	return &ErrorResponse{
		Message: message,
		Status:  http.StatusForbidden,
		Error:   safeError,
	}
}

// AuthResponse is a response to a successful authentication request. It
// contains the `token` field which is the JWT token used on other endpoints
// that require authentication.
type AuthResponse struct {
	Token string `json:"token"`
}

func (e *AuthResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, http.StatusOK)
	return nil
}

// NewAuthRouter creates a router for the authentication endpoints.
//
// POST /: Authenticate a user using email and password, and return
//     a JWT if correct.
func NewAuthRouter(logger *zap.Logger, store UserStorer) func(chi.Router) {
	validate := validator.New()
	authService := AuthService{store}
	return func(r chi.Router) {
		r.Post("/", signIn(logger, validate, authService))
	}
}

func NewValidator() *validator.Validate {
	return validator.New()
}

type AuthnRequest struct {
	Email    string `json:"email" validate:"email,required"`
	Password string `json:"password" validate:"min=6,required"`
}

func writeError(status int, resp ErrorResponse, w http.ResponseWriter) error {
	body, err := json.Marshal(resp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return err
	}
	_, err = w.Write(body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return err
	}
	w.WriteHeader(status)
	return nil
}

func signIn(logger *zap.Logger, validate *validator.Validate, service AuthService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			logger.Info("Failed to read request body", zap.Error(err))
			render.Render(w, r, ErrInternal(err))
			return
		}
		defer r.Body.Close()

		var request AuthnRequest
		err = json.Unmarshal(body, &request)
		if err != nil {
			logger.Info("Failed to unmarshal request", zap.Error(err), zap.ByteString("body", body))
			resp := ErrMalformedRequest("Request body is invalid JSON", err)
			render.Render(w, r, resp)
			return
		}

		err = validate.Struct(request)
		if err != nil {
			logger.Info("Received invalid request body", zap.Error(err), zap.Any("body", request))
			var errs []error
			for _, err := range err.(validator.ValidationErrors) {
				errs = append(errs, FieldError(err))
			}
			resp := ErrInvalidRequest("Sign in request is not valid", errs)
			render.Render(w, r, resp)
			return
		}
		logger.Info("Received sign in request", zap.String("email", request.Email))

		user, authErr := service.AuthenticateUser(request.Email, request.Password)
		if authErr != nil {
			logger.Info("Authentication failed", zap.String("email", request.Email), zap.Error(authErr))
			resp := ErrForbidden(
				fmt.Sprintf("Could not authenticate user %s", request.Email),
				authErr,
			)
			render.Render(w, r, resp)
			return
		}

		logger.Info("Successfully authenticated user",
			zap.String("email", user.Email),
		)

		token, tokenErr := service.GetToken(user)
		if err != nil {
			logger.Info("Failed to get auth token for user",
				zap.String("email", user.Email), zap.Error(tokenErr),
			)
			resp := ErrInternalWithMessage(
				"Could not get authentication token for user",
				fmt.Errorf("%s", tokenErr.SafeError()),
			)
			render.Render(w, r, resp)
			return
		}

		logger.Info("Successfully retrieved auth token for user", zap.String("email", user.Email))
		// Set the cookie header for use in web app
		cookie := &http.Cookie{
			Name:     "jwt",
			Value:    token,
			HttpOnly: true,
			Path:     "/",
			SameSite: http.SameSiteNoneMode,
			Expires:  time.Now().Add(time.Hour * 24),
			Secure:   false, // TODO: Set this to secure for prod
			Domain:   "localhost.com",
		}
		http.SetCookie(w, cookie)
		logger.Debug("Set cookie", zap.String("cookie", cookie.String()))

		resp := &AuthResponse{Token: token}
		render.Render(w, r, resp)
	}
}
