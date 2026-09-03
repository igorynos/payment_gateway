package paymenthandler

import (
	"errors"
	"net/http"

	"payment_gateway/internal/lib/api/response"
	"payment_gateway/internal/payment"

	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

type CreateRequest struct {
	Invoice           string `json:"invoice" validate:"required"`
	Status            string `json:"status" validate:"required"`
	Amount            int64  `json:"amount" validate:"required,gt=0"`
	Currency          string `json:"currency" validate:"required"`
	Provider          string `json:"provider" validate:"required"`
	ProviderPaymentID string `json:"provider_payment_id"`
}

type CreateResponse struct {
	response.Response
	Payment PaymentResponse `json:"payment"`
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateRequest
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, response.Error("invalid request"))
		return
	}

	if err := validator.New().Struct(req); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, response.Error("validation failed"))
		return
	}

	createdPayment, err := h.service.CreatePayment(
		r.Context(),
		payment.CreateInput{
			Invoice:           req.Invoice,
			Status:            payment.Status(req.Status),
			Amount:            req.Amount,
			Currency:          payment.Carrency(req.Currency),
			Provider:          req.Provider,
			ProviderPaymentID: req.ProviderPaymentID,
		},
	)
	if errors.Is(err, payment.ErrInvalidInput) {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, response.Error("invalid payment data"))
		return
	}
	if err != nil {
		h.log.Error("failed to create payment", "error", err)
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, response.Error("failed to create payment"))
		return
	}

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, CreateResponse{
		Response: response.OK(),
		Payment:  newPaymentResponse(createdPayment),
	})
}
