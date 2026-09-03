package paymenthandler

import (
	"errors"
	"net/http"

	"payment_gateway/internal/lib/api/response"
	"payment_gateway/internal/payment"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

type GetByIDResponse struct {
	response.Response
	Payment PaymentResponse `json:"payment"`
}

func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "paymentID")
	if id == "" {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, response.Error("invalid payment ID"))
		return
	}

	foundPayment, err := h.service.GetPaymentByID(r.Context(), id)
	if errors.Is(err, payment.ErrNotFound) {
		render.Status(r, http.StatusNotFound)
		render.JSON(w, r, response.Error("payment not found"))
		return
	}
	if errors.Is(err, payment.ErrInvalidInput) {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, response.Error("invalid payment ID"))
		return
	}
	if err != nil {
		h.log.Error("failed to get payment", "error", err, "payment_id", id)
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, response.Error("failed to get payment"))
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, GetByIDResponse{
		Response: response.OK(),
		Payment:  newPaymentResponse(foundPayment),
	})
}
