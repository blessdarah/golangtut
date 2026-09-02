package config

import (
	"log/slog"
	"net/http"
)

// slog middleware
func slogMiddlware(logger *slog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger.Info(
				"Request",
				slog.String("RequestID", r.Header.Get("X-Request-ID")),
				slog.String("method", r.Method),
				slog.String("url", r.URL.String()),
			)
			next.ServeHTTP(w, r)
		})
	}
}
