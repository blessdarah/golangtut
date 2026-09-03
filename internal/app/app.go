package app

import (
	"blessdarah/tuts/internal/auth"
	"blessdarah/tuts/internal/config"
	"blessdarah/tuts/internal/event"
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
	authMiddleware func(http.Handler) http.Handler,
	eventHandler *event.Handler,
) {

	// auth routes
	a.router.Post("/auth/signup", authHandler.Signup)
	a.router.Post("/auth/login", authHandler.Token) // login

	a.router.Get("/events/me", eventHandler.GetByUserID)
	a.router.Get("/events", eventHandler.GetAll)

	a.router.Group(func(r chi.Router) {
		r.Use(authMiddleware)
		r.Get("/auth/me", authHandler.Me)
		r.Get("/users", userHandler.GetUsers)

		// event routes
		r.Post("/events", eventHandler.Create)

	})

}

func (a *App) Run() {
	userRepo := user.NewRepository(a.db)
	userSvc := user.NewService(userRepo)
	userHandler := user.NewHandler(a.config, userSvc)
	eventRepo := event.NewRepository(a.db)

	authService := auth.NewService(userRepo)
	oauthServer, err := auth.NewOAuthServer(
		a.config.OAuthClientID,
		a.config.OAuthClientSecret,
		a.config.OAuthAccessTokenTTLMinutes,
		a.config.OAuthRefreshTokenTTLHours,
		authService,
	)
	if err != nil {
		log.Fatal(err)
	}

	eventService := event.NewService(eventRepo)

	authHandler := auth.NewAuthHandler(
		authService,
		oauthServer,
		a.logger.With(slog.String("module", "auth")),
	)
	authMiddleware := auth.RequireBearer(
		oauthServer,
		a.logger.With(slog.String("module", "auth-middleware")),
	)

	eventHandler := event.NewHandler(
		eventService,
		a.logger.With(slog.String("module", "event")),
	)

	a.RegisterRoutes(userHandler, authHandler, authMiddleware, eventHandler)

	a.logger.Info(fmt.Sprintf("Server is running on port %s", a.config.AppPort))

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", a.config.AppPort), a.router))
}
