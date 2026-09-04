package webhookhandler

import (
	"context"
	"log/slog"

	"payment_gateway/internal/payment"
)

type PaymentService interface {
	UpdateStatusPaymentByID(
		ctx context.Context,
		input payment.StatusUpdateInput,
	) (payment.Payment, error)
}

type Handler struct {
	log     *slog.Logger
	service PaymentService
}

func New(
	log *slog.Logger,
	service PaymentService,
) *Handler {
	return &Handler{
		log:     log,
		service: service,
	}
}
