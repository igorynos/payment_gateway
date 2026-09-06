package webhookhandler

import (
	"net/http"
	"time"

	"payment_gateway/internal/lib/api/response"
	"payment_gateway/internal/payment"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

type UpdateStatusRequest struct {
	EventID    string    `json:"event_id" validate:"required"`
	PaymentID  string    `json:"payment_id" validate:"required"`
	Status     string    `json:"status" validate:"required"`
	OccurredAt time.Time `json:"occurred_at" validate:"required"`
}

func (h *Handler) UpdateStatusByID(
	w http.ResponseWriter,
	r *http.Request,
) {
	provider := chi.URLParam(r, "provider")

	var req UpdateStatusRequest
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, response.Error("invalid request"))
		return
	}
	if err := validator.New().Struct(req); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, response.Error("invalid request"))
		return
	}

	command := payment.StatusChangeCommand{
		EventID:    req.EventID,
		PaymentID:  req.PaymentID,
		Provider:   provider,
		Status:     payment.Status(req.Status),
		OccurredAt: req.OccurredAt,
		Version:    1,
	}

	err := h.publisher.Publish(
		r.Context(),
		command.PaymentID,
		command,
	)

	if err != nil {
		h.log.Error(
			"failed to publish payment status change",
			"payment_id", command.PaymentID,
			"event_id", command.EventID,
			"error", err,
		)

		render.Status(r, http.StatusServiceUnavailable)
		render.JSON(
			w,
			r,
			response.Error("failed to enqueue webhook"),
		)
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, response.OK())
}
