package config

import (
	"log/slog"
	"os"
)

func NewLogger(cfg *AppEnv) *slog.Logger {

	var logger *slog.Logger
	if cfg.Debug {
		logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}))
	} else {
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		}))
	}

	return logger.With(
		slog.String("env", cfg.Env),
		slog.String("host", cfg.AppHost),
		slog.String("version", "1.0.0"),
	)
}
