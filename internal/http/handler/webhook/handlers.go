package webhookhandler

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
	log       *slog.Logger
	publisher Publisher
}

func New(
	log *slog.Logger,
	publisher Publisher,
) *Handler {
	return &Handler{
		log:       log,
		publisher: publisher,
	}
}
