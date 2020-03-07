package middleware

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/middleware"
	"go.uber.org/zap"
)

// Logger is a middlware http handler that logs requests.
func Logger(logger *zap.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			receivedAt := time.Now()
			dateReceived := receivedAt.Format("2006-01-02T15:0405-0700")
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			defer func() {
				logger.Info(
					"Received request",
					zap.String("remoteHost", r.RemoteAddr),
					zap.String("method", r.Method),
					zap.String("uri", r.RequestURI),
					zap.String("protocol", r.Proto),
					zap.Int("status", ww.Status()),
					zap.Int("responseSize", ww.BytesWritten()),
					zap.String("dateReceived", dateReceived),
					zap.Duration("duration", time.Since(receivedAt)),
				)
			}()
			next.ServeHTTP(ww, r)
		})
	}
}
