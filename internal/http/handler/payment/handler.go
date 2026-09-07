package paymenthandler

import (
	"context"
	"log/slog"

	"payment_gateway/internal/payment"
)

type Publisher interface {
	Publish(
		ctx context.Context,
		key string,
		message any,
	) error
}

type PaymentService interface {
	GetPaymentByID(
		ctx context.Context,
		id string,
	) (payment.Payment, error)
}

type Handler struct {
	log             *slog.Logger
	createPublisher Publisher
	service         PaymentService
}

func New(
	log *slog.Logger,
	createPublisher Publisher,
	service PaymentService,
) *Handler {
	return &Handler{
		log:             log,
		createPublisher: createPublisher,
		service:         service,
	}
}
