package middleware

import (
	"net/http"
	"os"
	"time"

	"github.com/rs/zerolog"
)

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
		logger := zerolog.New(os.Stdout).With().Timestamp().Logger()

		startTime := time.Now()

		logger.Info().
			Str("method", r.Method).
			Str("uri", r.RequestURI).
			Msg("Received request")

		rw := &responseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
			size:           0,
		}

		next.ServeHTTP(rw, r)

		logger.Info().
			Int("status", rw.statusCode).
			Int64("size", rw.size).
			Dur("duration", time.Since(startTime)).
			Msg("SentResponse")
	})
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
	size       int64
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	size, err := rw.ResponseWriter.Write(b)
	rw.size += int64(size)
	return size, err
}
