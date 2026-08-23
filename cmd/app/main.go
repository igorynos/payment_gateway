package main

import (
	"io"
	"log"
	"log/slog"
	"net/http"

	"payment_gateway/internal/config"
	"payment_gateway/internal/logger"

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

	logger.Info("Start app with config", slog.String("env", config.Env))
	logger.Debug("Level logs this app is Debug")
}
