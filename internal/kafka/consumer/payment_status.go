package consumer

import (
	"context"
	"encoding/json"
	"fmt"

	"payment_gateway/internal/payment"
	"payment_gateway/internal/worker"
)

type PaymentStatusUpdater interface {
	UpdateStatusPaymentByID(
		ctx context.Context,
		input payment.StatusUpdateInput,
	) (payment.Payment, error)
}

type PaymentStatusHandler struct {
	service PaymentStatusUpdater
}

func NewPaymentStatusHandler(
	service PaymentStatusUpdater,
) *PaymentStatusHandler {
	return &PaymentStatusHandler{
		service: service,
	}
}

func (h *PaymentStatusHandler) Handle(
	ctx context.Context,
	message worker.Message,
) error {
	var command payment.StatusChangeCommand

	if err := json.Unmarshal(message.Value, &command); err != nil {
		return fmt.Errorf(
			"decode payment status command: %w",
			err,
		)
	}

	if command.Version != 1 {
		return fmt.Errorf(
			"unsupported payment status command version: %d",
			command.Version,
		)
	}

	if command.EventID == "" {
		return fmt.Errorf("payment status event ID is empty")
	}

	_, err := h.service.UpdateStatusPaymentByID(
		ctx,
		payment.StatusUpdateInput{
			ID:       command.PaymentID,
			Provider: command.Provider,
			Status:   command.Status,
		},
	)
	if err != nil {
		return fmt.Errorf(
			"update payment status: %w",
			err,
		)
	}

	return nil
}
