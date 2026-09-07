package payment

import (
	"errors"
	"time"
)

var (
	ErrNotFound     = errors.New("payment not found")
	ErrInvalidInput = errors.New("invalid payment input")
)

type Status string

const (
	StatusNew        Status = "NEW"
	StatusProcessing Status = "PROCESSING"
	StatusCompleted  Status = "COMPLETED"
	StatusFailed     Status = "FAILED"
	StatusCanceled   Status = "CANCELED"
)

type Currency string

const (
	Rub Currency = "RUB"
	Usd Currency = "USD"
	Eur Currency = "EUR"
)

type Payment struct {
	ID                string
	Invoice           string
	Status            Status
	Amount            int64
	Currency          Currency
	Provider          string
	ProviderPaymentID string
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

type CreateParams struct {
	Invoice           string
	Status            Status
	Amount            int64
	Currency          Currency
	Provider          string
	ProviderPaymentID string
}

type CreateInput struct {
	Invoice           string
	Status            Status
	Amount            int64
	Currency          Currency
	Provider          string
	ProviderPaymentID string
}

type StatusUpdateParams struct {
	ID       string
	Provider string
	Status   Status
}

type StatusUpdateInput struct {
	ID       string
	Provider string
	Status   Status
}

type StatusChangeCommand struct {
	EventID    string    `json:"event_id"`
	PaymentID  string    `json:"payment_id"`
	Provider   string    `json:"provider"`
	Status     Status    `json:"status"`
	OccurredAt time.Time `json:"occurred_at"`
	Version    int       `json:"version"`
}

type CreateCommand struct {
	RequestID         string    `json:"request_id"`
	Invoice           string    `json:"invoice"`
	Amount            int64     `json:"amount"`
	Currency          Currency  `json:"currency"`
	Provider          string    `json:"provider"`
	ProviderPaymentID string    `json:"provider_payment_id,omitempty"`
	OccurredAt        time.Time `json:"occurred_at"`
	Version           int       `json:"version"`
}

type GetCommand struct {
	RequestID  string    `json:"request_id"`
	PaymentID  string    `json:"payment_id"`
	OccurredAt time.Time `json:"occurred_at"`
	Version    int       `json:"version"`
}
