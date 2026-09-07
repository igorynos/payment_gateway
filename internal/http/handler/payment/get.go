package paymenthandler

import (
	"net/http"
	"time"

	"payment_gateway/internal/lib/api/response"
	"payment_gateway/internal/payment"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type GetByIDResponse struct {
	response.Response
	RequestID string `json:"request_id"`
	PaymentID string `json:"payment_id"`
}

func (h *Handler) GetByID(
	w http.ResponseWriter,
	r *http.Request,
) {
	paymentID := chi.URLParam(r, "paymentID")
	if paymentID == "" {
		render.Status(r, http.StatusBadRequest)
		render.JSON(
			w,
			r,
			response.Error("invalid payment ID"),
		)
		return
	}

	requestID := middleware.GetReqID(r.Context())
	if requestID == "" {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(
			w,
			r,
			response.Error("request ID is missing"),
		)
		return
	}

	command := payment.GetCommand{
		RequestID:  requestID,
		PaymentID:  paymentID,
		OccurredAt: time.Now().UTC(),
		Version:    1,
	}

	if err := h.getPublisher.Publish(
		r.Context(),
		command.PaymentID,
		command,
	); err != nil {
		h.log.Error(
			"failed to publish payment get command",
			"request_id", command.RequestID,
			"payment_id", command.PaymentID,
			"error", err,
		)

		render.Status(r, http.StatusServiceUnavailable)
		render.JSON(
			w,
			r,
			response.Error("failed to enqueue payment retrieval"),
		)
		return
	}

	render.Status(r, http.StatusAccepted)
	render.JSON(w, r, GetByIDResponse{
		Response:  response.OK(),
		RequestID: command.RequestID,
		PaymentID: command.PaymentID,
	})
}
