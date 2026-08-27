package postgres

import (
	"context"
	"errors"
	"fmt"
	"payment_gateway/internal/storage/postgres/sqlc"
	"payment_gateway/internal/user"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type UserRepository struct {
	queries *sqlc.Queries
}

func NewUserRepository(
	storage *Storage,
) *UserRepository {
	return &UserRepository{
		queries: storage.queries,
	}
}

func (r *UserRepository) Create(
	ctx context.Context,
	params user.CreateParams,
) (user.User, error) {
	dbUser, err := r.queries.CreateUser(
		ctx,
		sqlc.CreateUserParams{
			Username:     params.Username,
			PasswordHash: params.PasswordHash,
			Email:        params.Email,
			Role:         string(params.Role),
		},
	)
	if err != nil {
		return user.User{}, mapUserError(err)
	}
	return toDomainUser(dbUser), nil
}

func (r *UserRepository) GetByID(
	ctx context.Context,
	id int64,
) (user.User, error) {
	dbUser, err := r.queries.GetUserByID(ctx, id)
	if errors.Is(err, pgx.ErrNoRows) {
		return user.User{}, user.ErrNotFound
	}
	if err != nil {
		return user.User{}, fmt.Errorf("get user by ID: %w", err)
	}
	return toDomainUser(dbUser), nil
}

func toDomainUser(dbUser sqlc.User) user.User {
	result := user.User{
		ID:       dbUser.ID,
		Username: dbUser.Username,
		Email:    dbUser.Email,
		Role:     user.Role(dbUser.Role),
	}
	if dbUser.CreatedAt.Valid {
		result.CreatedAt = dbUser.CreatedAt.Time
	}
	if dbUser.UpdatedAt.Valid {
		result.UpdatedAt = dbUser.UpdatedAt.Time
	}
	return result
}

func mapUserError(err error) error {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		if pgErr.Code == "23505" {
			return user.ErrAlreadyExists
		}
	}
	return fmt.Errorf("postgres user repository: %w", err)
}
