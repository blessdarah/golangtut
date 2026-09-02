package config

import (
	chi "github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Registra interface {
	Register(r chi.Router)
}

func NewRouter(registrars ...Registra) chi.Router {
	r := chi.NewRouter()

	logger := NewLogger(LoadConfig())

	r.Use(middleware.RequestID)
	r.Use(slogMiddlware(logger))
	r.Use(middleware.Recoverer)
	r.Use(middleware.SetHeader("Content-Type", "application/json"))

	for _, registrar := range registrars {
		registrar.Register(r)
	}

	return r
}
