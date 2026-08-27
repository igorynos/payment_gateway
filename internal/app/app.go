package app

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"payment_gateway/internal/config"
	userhandler "payment_gateway/internal/http/handler/user"
	httprouter "payment_gateway/internal/http/router"
	"payment_gateway/internal/storage/postgres"
	"payment_gateway/internal/user"
)

type App struct {
	server  *http.Server
	storage *postgres.Storage
}

func New(
	ctx context.Context,
	cfg *config.Config,
	log *slog.Logger,
) (*App, error) {
	db, err := postgres.SetupDatabase(ctx, cfg.Storage, 10)
	if err != nil {
		log.Error("failed to connect to database", slog.Any("error", err))
		return nil, err
	}

	userRepository := postgres.NewUserRepository(db)
	userService := user.NewService(userRepository)
	userHandler := userhandler.New(log, userService)

	httpHandler := httprouter.New(
		log,
		userHandler,
	)
	address := fmt.Sprintf(
		"%s:%d",
		cfg.Address,
		cfg.Port,
	)

	server := &http.Server{
		Addr:         address,
		Handler:      httpHandler,
		ReadTimeout:  cfg.Timeout,
		WriteTimeout: cfg.Timeout,
		IdleTimeout:  cfg.Idle_timeout,
	}
	return &App{
		server:  server,
		storage: db,
	}, nil

}

func (a *App) Run() error {
	return a.server.ListenAndServe()
}

func (a *App) Close() {
	if a.storage != nil {
		a.storage.Close()
	}
}
