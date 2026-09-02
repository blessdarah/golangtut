package app

import (
	"blessdarah/tuts/internal/config"
	"blessdarah/tuts/internal/user"
	"fmt"
	"log"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"gorm.io/gorm"
)

type App struct {
	config *config.AppEnv
	logger *slog.Logger
	db     *gorm.DB
	router chi.Router
}

func NewApp(cfg *config.AppEnv) (*App, error) {
	logger := config.NewLogger(cfg)

	db, err := config.ConnectDB(cfg)
	if err != nil {
		return nil, err
	}

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(config.SlogMiddlware(logger))
	r.Use(middleware.Recoverer)
	r.Use(middleware.SetHeader("Content-Type", "application/json"))

	return &App{
		config: cfg,
		logger: logger,
		db:     db,
		router: r,
	}, nil
}

// RegisterRoutes registers routes for the app
func (a *App) RegisterRoutes(userHandler *user.Handler) {
	a.router.Route("/users", func(r chi.Router) {
		r.Get("/", userHandler.GetUsers)
		r.Post("/", userHandler.CreateUser)
	})
}

func (a *App) Run() {

	// repositories
	repo := user.NewRepository(a.db)

	// app
	userApp := user.NewService(repo)

	// handlers
	userHandler := user.NewHandler(a.config, userApp)

	// register routes
	a.RegisterRoutes(userHandler)

	// start server
	a.logger.Info(fmt.Sprintf("Server is running on port %s", a.config.AppPort))

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", a.config.AppPort), a.router))
}
