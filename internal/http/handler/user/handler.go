package userhandler

import (
	"context"
	"log/slog"

	"payment_gateway/internal/user"
)

type UserService interface {
	CreateUser(
		ctx context.Context,
		input user.CreateInput,
	) (user.User, error)
	GetByID(
		ctx context.Context,
		id int64,
	) (user.User, error)
}

type Handler struct {
	log     *slog.Logger
	service UserService
}

func New(
	log *slog.Logger,
	service UserService,
) *Handler {
	return &Handler{
		log:     log,
		service: service,
	}
}
