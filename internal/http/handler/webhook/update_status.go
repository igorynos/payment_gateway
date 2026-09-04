package webhookhandler

import (
	"net/http"

	"payment_gateway/internal/lib/api/response"
	"payment_gateway/internal/payment"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

type UpdateStatusRequest struct {
	PaymentID string `json:"payment_id" validate:"required"`
	Status    string `json:"status" validate:"required"`
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
	_, err := h.service.UpdateStatusPaymentByID(
		r.Context(),
		payment.StatusUpdateInput{
			ID:       req.PaymentID,
			Provider: provider,
			Status:   payment.Status(req.Status),
		},
	)
	if err != nil {
		h.log.Error("failed to process webhook", "error", err)
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, response.Error("failed to process webhook"))
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, response.OK())
}
