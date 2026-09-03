package paymenthandler

import "payment_gateway/internal/payment"

type PaymentResponse struct {
	ID                string `json:"id"`
	Invoice           string `json:"invoice"`
	Status            string `json:"status"`
	Amount            int64  `json:"amount"`
	Currency          string `json:"currency"`
	Provider          string `json:"provider"`
	ProviderPaymentID string `json:"provider_payment_id"`
}

func newPaymentResponse(value payment.Payment) PaymentResponse {
	return PaymentResponse{
		ID:                value.ID,
		Invoice:           value.Invoice,
		Status:            string(value.Status),
		Amount:            value.Amount,
		Currency:          string(value.Currency),
		Provider:          value.Provider,
		ProviderPaymentID: value.ProviderPaymentID,
	}
}
