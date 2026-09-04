package postgres

import (
	"context"
	"errors"
	"fmt"
	"math/big"

	"payment_gateway/internal/payment"
	"payment_gateway/internal/storage/postgres/sqlc"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
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
			Invoice: params.Invoice,
			Status:  string(params.Status),
			Amount: pgtype.Numeric{
				Int:   big.NewInt(params.Amount),
				Valid: true,
			},
			Currency: string(params.Currency),
			Provider: params.Provider,
			ProviderPaymentID: pgtype.Text{
				String: params.ProviderPaymentID,
				Valid:  params.ProviderPaymentID != "",
			},
		},
	)
	if err != nil {
		return payment.Payment{}, fmt.Errorf("create payment: %w", err)
	}
	return toDomainPayment(dbPayment)
}

func (r *PaymentRepository) GetPaymentByID(
	ctx context.Context,
	id string,
) (payment.Payment, error) {
	var paymentID pgtype.UUID
	if err := paymentID.Scan(id); err != nil {
		return payment.Payment{}, payment.ErrInvalidInput
	}

	dbPayment, err := r.queries.GetPaymentByID(
		ctx,
		paymentID,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return payment.Payment{}, payment.ErrNotFound
	}
	if err != nil {
		return payment.Payment{}, fmt.Errorf("get payment by ID: %w", err)
	}
	return toDomainPayment(dbPayment)
}

func (r *PaymentRepository) UpdatePaymentStatusByID(
	ctx context.Context,
	params payment.StatusUpdateParams,
) (payment.Payment, error) {
	var paymentID pgtype.UUID
	if err := paymentID.Scan(params.ID); err != nil {
		return payment.Payment{}, payment.ErrInvalidInput
	}
	dbPayment, err := r.queries.UpdatePaymentStatusByID(
		ctx,
		sqlc.UpdatePaymentStatusByIDParams{
			Status:   string(params.Status),
			ID:       paymentID,
			Provider: params.Provider,
		},
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return payment.Payment{}, payment.ErrNotFound
	}
	if err != nil {
		return payment.Payment{}, fmt.Errorf(
			"update payment status by ID: %w",
			err,
		)
	}
	return toDomainPayment(dbPayment)
}

func toDomainPayment(dbPayment sqlc.Payment) (payment.Payment, error) {
	if !dbPayment.ID.Valid {
		return payment.Payment{}, errors.New("payment ID is NULL")
	}

	amount, err := dbPayment.Amount.Int64Value()
	if err != nil {
		return payment.Payment{}, fmt.Errorf("convert payment amount: %w", err)
	}
	if !amount.Valid {
		return payment.Payment{}, errors.New("payment amount is NULL")
	}

	result := payment.Payment{
		ID:                dbPayment.ID.String(),
		Invoice:           dbPayment.Invoice,
		Status:            payment.Status(dbPayment.Status),
		Amount:            amount.Int64,
		Currency:          payment.Carrency(dbPayment.Currency),
		Provider:          dbPayment.Provider,
		ProviderPaymentID: dbPayment.ProviderPaymentID.String,
		CreatedAt:         dbPayment.CreatedAt.Time,
		UpdatedAt:         dbPayment.UpdatedAt.Time,
	}
	return result, nil
}
