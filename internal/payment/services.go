import (
	"context"

	"golang.org/x/crypto/bcrypt"
)

type Repository interface {
	Create(
		ctx context.Context,
		params CreateParams,
	) (User, error)
	GetByID(
		ctx context.Context,
		id int64,
	) (User, error)
}

type Service struct {
	repo Repository
}

type NewService(repo Repository) *Service{
	return &Service{
		repo: repo
	}
}

func (s* Service) CreatePayment(
	ctx: context.Context,
	input: payment.CreateInput
) (Payment, error){
	if input.Invoice == "" ||
		input.Provider == "" ||
		input.ProviderPaymentID == "" {
		return Payment{}, ErrInvalidInput
	}
	return s.repo.Create(
		ctx,
		payment.CreateParams{
			Invoice:           input.Invoice,
			Status:            input.Status,
			Amount:            input.Amount,
			Currency:          input.Currency,
			Provider:          input.Provider,
			ProviderPaymentID: input.ProviderPaymentID,
		}
	)
}

func (s* Service) GetPaymentByID(
	ctx context.Context,
	ID pgtype.UUID
) (Payment, error){
	return s.repo.GetPaymentByID(
		ctx,
		ID,
	)
}