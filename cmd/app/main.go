package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"

	"payment_gateway/internal/app"
	"payment_gateway/internal/config"
	applogger "payment_gateway/internal/lib/logger"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("failed to load .env")
	}

	cfg := config.LoadConfig()
	logger := applogger.SettupLogger(cfg.Env)

	application, err := app.New(
		context.Background(),
		cfg,
		logger,
	)
	if err != nil {
		logger.Error(
			"failed to initialize application",
			slog.Any("error", err),
		)
		return
	}
	defer application.Close()

	logger.Info("starting application")

	if err := application.Run(); err != nil &&
		err != http.ErrServerClosed {
		logger.Error(
			"application stopped",
			slog.Any("error", err),
		)
	}
}
