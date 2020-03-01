package user

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/go-chi/chi"
	"go.uber.org/zap"
	"gopkg.in/go-playground/validator.v9"
)

type ErrorResponse struct {
	Message string `json:"message"`
	Error   error  `json:"error,omitempty"`
}

func NewAuthRouter(logger *zap.Logger, store UserStorer) func(chi.Router) {
	validate := validator.New()
	return func(r chi.Router) {
		r.Post("/", signIn(logger, validate, store))
	}
}

type authnRequest struct {
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

func signIn(logger *zap.Logger, validate *validator.Validate, store UserStorer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			logger.Error("Failed to read request body", zap.Error(err))
			errResp := ErrorResponse{Message: "Could not read request body"}
			err = writeError(http.StatusInternalServerError, errResp, w)
			if err != nil {
				logger.Error("Failed to write error response", zap.Error(err), zap.Any("response", errResp))
			}
			return
		}

		var request authnRequest
		err = json.Unmarshal(body, &request)
		if err != nil {
			logger.Error("Failed to unmarshal request", zap.Error(err), zap.ByteString("body", body))
			errResp := ErrorResponse{Message: "Request body is invalid JSON", Error: err}
			err = writeError(http.StatusBadRequest, errResp, w)
			if err != nil {
				logger.Error("Failed to write error response", zap.Error(err), zap.Any("response", errResp))
			}
			return
		}

		err = validate.Struct(request)
		if err != nil {
			logger.Error("Received invalid request body", zap.Error(err), zap.Any("body", request))
			errResp := ErrorResponse{Message: "Request is not valid", Error: err}
			err = writeError(http.StatusBadRequest, errResp, w)
			if err != nil {
				logger.Error("Failed to write error response", zap.Error(err), zap.Any("response", errResp))
			}
			return
		}
		logger.Info("Received sign in request", zap.String("email", request.Email))
	}
}

