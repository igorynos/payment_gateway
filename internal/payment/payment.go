package payment

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
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
	ID                pgtype.UUID
	Invoice           string
	Status            Status
	Amount            pgtype.Numeric
	Currency          Carrency
	Provider          string
	ProviderPaymentID string
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

type CreateParams struct {
	Invoice           string
	Status            string
	Amount            pgtype.Numeric
	Currency          string
	Provider          string
	ProviderPaymentID pgtype.Text
}

type CreateInput struct {
	Invoice           string
	Status            Status
	Amount            pgtype.Numeric
	Currency          Carrency
	Provider          string
	ProviderPaymentID string
}
