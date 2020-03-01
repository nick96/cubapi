package middleware

import (
	"net/http"
	"time"

	"go.uber.org/zap"
)

type responseWriter struct {
	writer       http.ResponseWriter
	responseSize int
	responseStatus int
}

func (w *responseWriter) Header() http.Header {
	return w.writer.Header()
}

func (w *responseWriter) Write(contents []byte) (int, error) {
	n, err := w.writer.Write(contents)
	if err != nil {
		w.responseSize += n
	}
	return n, err
}

func (w *responseWriter) WriteHeader(status int) {
	w.responseStatus = status
	w.writer.WriteHeader(status)
}

// Logger is a middlware http handler that logs in the common log format.
func Logger(logger *zap.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedAt := time.Now()

		writer := responseWriter{writer: w}
		next.ServeHTTP(&writer, r)

		respondedAt := time.Now()

		dateReceived := receivedAt.Format("02-10-2006T15:04:05-0700")
		duration := respondedAt.Sub(receivedAt)

		logger.Info(
			"Received request",
			zap.String("remoteHost", r.RemoteAddr),
			zap.String("method", r.Method),
			zap.String("uri", r.RequestURI),
			zap.String("protocol", r.Proto),
			zap.Int("status", writer.responseStatus),
			zap.Int("responseSize", writer.responseSize),
			zap.String("dateReceived", dateReceived),
			zap.Duration("duration", duration),
		)
	})
}
