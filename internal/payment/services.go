package payment

import (
	"context"
)

type Repository interface {
	Create(
		ctx context.Context,
		params CreateParams,
	) (Payment, error)
	GetPaymentByID(
		ctx context.Context,
		id string,
	) (Payment, error)
	UpdatePaymentStatusByID(
		ctx context.Context,
		params StatusUpdateParams,
	) (Payment, error)
}

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) CreatePayment(
	ctx context.Context,
	input CreateInput,
) (Payment, error) {
	if input.Invoice == "" ||
		input.Provider == "" ||
		input.Amount <= 0 ||
		!validStatus(input.Status) ||
		!validCurrency(input.Currency) {
		return Payment{}, ErrInvalidInput
	}
	return s.repo.Create(
		ctx,
		CreateParams{
			Invoice:           input.Invoice,
			Status:            input.Status,
			Amount:            input.Amount,
			Currency:          input.Currency,
			Provider:          input.Provider,
			ProviderPaymentID: input.ProviderPaymentID,
		},
	)
}

func (s *Service) GetPaymentByID(
	ctx context.Context,
	id string,
) (Payment, error) {
	if id == "" {
		return Payment{}, ErrInvalidInput
	}
	return s.repo.GetPaymentByID(ctx, id)
}

func (s *Service) UpdateStatusPaymentByID(
	ctx context.Context,
	input StatusUpdateInput,
) (Payment, error) {
	if input.ID == "" ||
		!validProvider(input.Provider) ||
		!validStatus(input.Status) {
		return Payment{}, ErrInvalidInput
	}
	return s.repo.UpdatePaymentStatusByID(
		ctx,
		StatusUpdateParams{
			ID:       input.ID,
			Provider: input.Provider,
			Status:   input.Status,
		},
	)
}

func validStatus(status Status) bool {
	switch status {
	case StatusNew, StatusProcessing, StatusCompleted, StatusFailed, StatusCanceled:
		return true
	default:
		return false
	}
}

func validCurrency(currency Carrency) bool {
	switch currency {
	case Rub, Usd, Eur:
		return true
	default:
		return false
	}
}

func validProvider(provider string) bool {
	switch provider {
	case "paypal", "yookassa", "freekassa":
		return true
	default:
		return false
	}
}
