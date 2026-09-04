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

type Carrency string

const (
	Rub Carrency = "RUB"
	Usd Carrency = "USD"
	Eur Carrency = "EUR"
)

type Payment struct {
	ID                string
	Invoice           string
	Status            Status
	Amount            int64
	Currency          Carrency
	Provider          string
	ProviderPaymentID string
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

type CreateParams struct {
	Invoice           string
	Status            Status
	Amount            int64
	Currency          Carrency
	Provider          string
	ProviderPaymentID string
}

type CreateInput struct {
	Invoice           string
	Status            Status
	Amount            int64
	Currency          Carrency
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
