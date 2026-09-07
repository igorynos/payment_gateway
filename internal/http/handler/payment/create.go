package paymenthandler

import (
	"net/http"
	"time"

	"payment_gateway/internal/lib/api/response"
	"payment_gateway/internal/payment"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

type CreateRequest struct {
	Invoice           string `json:"invoice" validate:"required"`
	Amount            int64  `json:"amount" validate:"required,gt=0"`
	Currency          string `json:"currency" validate:"required"`
	Provider          string `json:"provider" validate:"required"`
	ProviderPaymentID string `json:"provider_payment_id"`
}

type CreateResponse struct {
	response.Response
	RequestID string `json:"request_id"`
}

func (h *Handler) Create(
	w http.ResponseWriter,
	r *http.Request,
) {
	var req CreateRequest

	if err := render.DecodeJSON(r.Body, &req); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(
			w,
			r,
			response.Error("invalid request"),
		)
		return
	}

	if err := validator.New().Struct(req); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(
			w,
			r,
			response.Error("validation failed"),
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

	command := payment.CreateCommand{
		RequestID:         requestID,
		Invoice:           req.Invoice,
		Amount:            req.Amount,
		Currency:          payment.Currency(req.Currency),
		Provider:          req.Provider,
		ProviderPaymentID: req.ProviderPaymentID,
		OccurredAt:        time.Now().UTC(),
		Version:           1,
	}

	if err := h.createPublisher.Publish(
		r.Context(),
		command.RequestID,
		command,
	); err != nil {
		h.log.Error(
			"failed to publish payment create command",
			"request_id", command.RequestID,
			"error", err,
		)

		render.Status(r, http.StatusServiceUnavailable)
		render.JSON(
			w,
			r,
			response.Error("failed to enqueue payment creation"),
		)
		return
	}

	render.Status(r, http.StatusAccepted)
	render.JSON(w, r, CreateResponse{
		Response:  response.OK(),
		RequestID: command.RequestID,
	})
}
