package middleware

import (
	"net/http"

	"go.uber.org/zap"
)

func DefaultContentType(logger *zap.Logger, contentType string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {
			if w.Header().Get("Content-Type") == "" {
				logger.Info("Setting default content type", zap.String("contentType", contentType))
				w.Header().Set("Content-Type", contentType)
			}
			next.ServeHTTP(w, r)
		})
	}
}
