package config

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
)

// slog middleware
func SlogMiddlware(logger *slog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			next.ServeHTTP(w, r)

			logger.Info(
				"Request",
				slog.String("requestId", middleware.GetReqID(r.Context())),
				slog.String("method", r.Method),
				slog.String("url", r.URL.String()),
				slog.Duration("duration", time.Since(start)),
				slog.Int("status", ww.Status()),
			)
		})
	}
}
