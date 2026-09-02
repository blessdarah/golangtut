package app

import (
	"blessdarah/tuts/internal/config"
	"log"
	"log/slog"

	"fmt"
	"net/http"

	"blessdarah/tuts/internal/user"
)

type App struct {
	config *config.AppEnv
	logger *slog.Logger
}

func NewApp(cfg *config.AppEnv) *App {
	return &App{
		config: cfg,
		logger: config.NewLogger(cfg),
	}
}

func (a *App) Run() {

	// handlers
	userHandler := user.NewHandler(a.config)

	// routers
	ur := user.NewUserRouter(userHandler)

	// register routes
	r := config.NewRouter(ur)

	// start server
	a.logger.Info(fmt.Sprintf("Server is running on port %s", a.config.AppPort))

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", a.config.AppPort), r))
}
