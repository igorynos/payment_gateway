package consumer

import (
	"context"
	"encoding/json"
	"fmt"

	"payment_gateway/internal/payment"
	"payment_gateway/internal/worker"
)

type PaymentCreator interface {
	CreatePayment(
		ctx context.Context,
		input payment.CreateInput,
	) (payment.Payment, error)
}

type PaymentCreateHandler struct {
	service PaymentCreator
}

func NewPaymentCreateHandler(
	service PaymentCreator,
) *PaymentCreateHandler {
	return &PaymentCreateHandler{
		service: service,
	}
}

func (h *PaymentCreateHandler) Handle(
	ctx context.Context,
	message worker.Message,
) error {
	var command payment.CreateCommand

	if err := json.Unmarshal(message.Value, &command); err != nil {
		return fmt.Errorf(
			"decode payment create command: %w",
			err,
		)
	}

	if command.Version != 1 {
		return fmt.Errorf(
			"unsupported payment create command version: %d",
			command.Version,
		)
	}

	if command.RequestID == "" {
		return fmt.Errorf("payment create request ID is empty")
	}

	_, err := h.service.CreatePayment(
		ctx,
		payment.CreateInput{
			Invoice:           command.Invoice,
			Status:            payment.StatusNew,
			Amount:            command.Amount,
			Currency:          command.Currency,
			Provider:          command.Provider,
			ProviderPaymentID: command.ProviderPaymentID,
		},
	)
	if err != nil {
		return fmt.Errorf("create payment: %w", err)
	}

	return nil
}
