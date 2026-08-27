package user

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

func NewService(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) CreateUser(
	ctx context.Context,
	input CreateInput,
) (User, error) {
	if input.Username == "" ||
		input.Email == "" ||
		input.Password == "" {
		return User{}, ErrInvalidInput
	}
	passwordHash, err := hashPassword(input.Password)
	if err != nil {
		return User{}, err
	}
	return s.repo.Create(ctx, CreateParams{
		Username:     input.Username,
		Email:        input.Email,
		PasswordHash: passwordHash,
		Role:         RoleUser,
	})
}

func (s *Service) GetByID(
	ctx context.Context,
	id int64,
) (User, error) {
	if id <= 0 {
		return User{}, ErrInvalidInput
	}
	return s.repo.GetByID(ctx, id)
}

func hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}
