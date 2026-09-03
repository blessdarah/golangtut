package app

import (
	"blessdarah/tuts/internal/auth"
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

func (a *App) RegisterRoutes(
	userHandler *user.Handler,
	authHandler *auth.Handler,
) {
	a.router.Route("/users", func(r chi.Router) {
		r.Get("/", userHandler.GetUsers)
	})

	a.router.Route("/auth", func(r chi.Router) {
		r.Post("/signup", authHandler.Signup)
		r.Get("/me", authHandler.Me)
	})

	a.router.Route("/oauth", func(r chi.Router) {
		r.Post("/token", authHandler.Token)
	})
}

func (a *App) Run() {
	userRepo := user.NewRepository(a.db)
	userSvc := user.NewService(userRepo)
	userHandler := user.NewHandler(a.config, userSvc)

	authService := auth.NewService(userRepo)
	oauthServer, err := auth.NewOAuthServer(
		a.config.OAuthClientID,
		a.config.OAuthClientSecret,
		authService,
	)
	if err != nil {
		log.Fatal(err)
	}

	authHandler := auth.NewAuthHandler(
		authService,
		oauthServer,
		a.logger.With(slog.String("module", "auth")),
	)

	a.RegisterRoutes(userHandler, authHandler)

	a.logger.Info(fmt.Sprintf("Server is running on port %s", a.config.AppPort))

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", a.config.AppPort), a.router))
}
