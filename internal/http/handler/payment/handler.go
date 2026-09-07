package paymenthandler

import (
	"context"
	"log/slog"
)

type Publisher interface {
	Publish(
		ctx context.Context,
		key string,
		message any,
	) error
}

type Handler struct {
	log             *slog.Logger
	createPublisher Publisher
	getPublisher    Publisher
}

func New(
	log *slog.Logger,
	createPublisher Publisher,
	getPublisher Publisher,
) *Handler {
	return &Handler{
		log:             log,
		createPublisher: createPublisher,
		getPublisher:    getPublisher,
	}
}
