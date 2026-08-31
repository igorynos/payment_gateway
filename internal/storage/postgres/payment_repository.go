package postgres

import (
	"context"
	"errors"
	"fmt"

	"payment_gateway/internal/payment"
	"payment_gateway/internal/storage/postgres/sqlc"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

var (
	ErrNotFound      = errors.New("payment not found")
	ErrAlreadyExists = errors.New("payment already exists")
	ErrInvalidInput  = errors.New("invalid payment input")
)

type PaymentRepository struct {
	queries *sqlc.Queries
}

func NewPaymentRepository(
	storage *Storage,
) *PaymentRepository {
	return &PaymentRepository{
		queries: storage.queries,
	}
}

func (r *PaymentRepository) Create(
	ctx context.Context,
	params payment.CreateParams,
) (payment.Payment, error) {
	dbPayment, err := r.queries.CreatePayment(
		ctx,
		sqlc.CreatePaymentParams{
			Invoice:           params.Invoice,
			Status:            params.Status,
			Amount:            params.Amount,
			Currency:          params.Currency,
			Provider:          params.Provider,
			ProviderPaymentID: params.ProviderPaymentID,
		},
	)
	if err != nil {
		return payment.Payment{}, ErrInvalidInput
	}
	return toDomainPayment(dbPayment), nil
}

func (r *PaymentRepository) GetPaymentByID(
	ctx context.Context,
	ID pgtype.UUID,
) (payment.Payment, error) {
	dbPayment, err := r.queries.GetPaymentByID(
		ctx,
		ID,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return payment.Payment{}, ErrNotFound
	}
	if err != nil {
		return payment.Payment{}, fmt.Errorf("get payment by ID: %w", err)
	}
	return toDomainPayment(dbPayment), nil
}

func toDomainPayment(dbPayment sqlc.Payment) payment.Payment {
	result := payment.Payment{
		ID:                dbPayment.ID,
		Invoice:           dbPayment.Invoice,
		Status:            payment.Status(dbPayment.Status),
		Amount:            dbPayment.Amount,
		Currency:          payment.Carrency(dbPayment.Currency),
		Provider:          dbPayment.Provider,
		ProviderPaymentID: dbPayment.ProviderPaymentID.String,
		CreatedAt:         dbPayment.CreatedAt.Time,
		UpdatedAt:         dbPayment.UpdatedAt.Time,
	}
	return result
}
