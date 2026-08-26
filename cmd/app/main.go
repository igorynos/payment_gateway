package main

import (
	"context"
	"io"
	"log"
	"log/slog"
	"net/http"

	"payment_gateway/internal/config"
	"payment_gateway/internal/logger"
	"payment_gateway/internal/storage/postgres"

	"github.com/joho/godotenv"
)

func say_hello(w http.ResponseWriter, _ *http.Request) {
	io.WriteString(w, "hello")
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Can't to load .env")
	}

	config := config.LoadConfig()
	logger := logger.SettupLogger(config.Env)

	ctx := context.Background()

	db, err := postgres.SetupDatabase(ctx, config.Storage, 10)
	if err != nil {
		logger.Error("failed to connect to database", slog.Any("error", err))
		return
	}
	defer db.Close()

	logger.Info("Start app with config", slog.String("env", config.Env))
	logger.Debug("Level logs this app is Debug")

}
