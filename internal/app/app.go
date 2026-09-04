package app

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"payment_gateway/internal/config"
	paymenthandler "payment_gateway/internal/http/handler/payment"
	userhandler "payment_gateway/internal/http/handler/user"
	webhookhandler "payment_gateway/internal/http/handler/webhook"
	httprouter "payment_gateway/internal/http/router"
	"payment_gateway/internal/payment"
	"payment_gateway/internal/storage/postgres"
	"payment_gateway/internal/user"
)

type App struct {
	apiServer     *http.Server
	webhookServer *http.Server
	storage       *postgres.Storage
}

func New(
	ctx context.Context,
	cfg *config.Config,
	log *slog.Logger,
) (*App, error) {
	db, err := postgres.SetupDatabase(ctx, cfg.Storage.URL(), 10)
	if err != nil {
		log.Error("failed to connect to database", slog.Any("error", err))
		return nil, err
	}

	userRepository := postgres.NewUserRepository(db)
	userService := user.NewService(userRepository)
	userHandler := userhandler.New(log, userService)

	paymentRepository := postgres.NewPaymentRepository(db)
	paymentService := payment.NewService(paymentRepository)
	paymentHandler := paymenthandler.New(log, paymentService)
	webhookHandler := webhookhandler.New(log, paymentService)

	apiRouter := httprouter.New(
		log,
		userHandler,
		paymentHandler,
	)

	webhookRouter := httprouter.New(
		log,
		webhookHandler,
	)

	apiServer := &http.Server{
		Addr: fmt.Sprintf(
			"%s:%d",
			cfg.HTTPServer.Address,
			cfg.HTTPServer.Port,
		),
		Handler:      apiRouter,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	webhookServer := &http.Server{
		Addr: fmt.Sprintf(
			"%s:%d",
			cfg.WebhookServer.Address,
			cfg.WebhookServer.Port,
		),
		Handler:      webhookRouter,
		ReadTimeout:  cfg.WebhookServer.Timeout,
		WriteTimeout: cfg.WebhookServer.Timeout,
		IdleTimeout:  cfg.WebhookServer.IdleTimeout,
	}
	return &App{
		apiServer:     apiServer,
		webhookServer: webhookServer,
		storage:       db,
	}, nil

}

func (a *App) Run() error {
	errCh := make(chan error, 2)

	go func() {
		errCh <- a.apiServer.ListenAndServe()
	}()

	go func() {
		errCh <- a.webhookServer.ListenAndServe()
	}()

	err := <-errCh
	if errors.Is(err, http.ErrServerClosed) {
		return nil
	}
	return err
}

func (a *App) Close() {
	if a.storage != nil {
		a.storage.Close()
	}
}
