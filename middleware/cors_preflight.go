package middleware

import (
	"net/http"

	"go.uber.org/zap"
)

func CORSPreflight(logger *zap.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "*")
			w.Header().Set("Access-Control-Allow-Headers", "*")
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			if r.Method == http.MethodOptions {
				logger.Info("Received preflight request", zap.String("path", r.URL.EscapedPath()))
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
