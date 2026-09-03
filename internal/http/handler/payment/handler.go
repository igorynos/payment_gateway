package paymenthandler

import (
	"context"
	"log/slog"

	"payment_gateway/internal/payment"
)

type PaymentService interface {
	CreatePayment(
		ctx context.Context,
		input payment.CreateInput,
	) (payment.Payment, error)
	GetPaymentByID(
		ctx context.Context,
		id string,
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
