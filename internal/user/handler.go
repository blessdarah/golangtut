package user

import (
	"blessdarah/tuts/internal/config"
	"log/slog"

	"net/http"
)

type Handler struct {
	cfg    *config.AppEnv
	logger *slog.Logger
}

func NewHandler(cfg *config.AppEnv) *Handler {
	return &Handler{
		cfg:    cfg,
		logger: config.NewLogger(cfg).With(slog.String("module", "user")),
	}
}

func (h *Handler) GetUsers(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("GetUsers")
	w.Write([]byte("GetUsers"))
}

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("CreateUser")
	w.Write([]byte("CreateUser"))
}
