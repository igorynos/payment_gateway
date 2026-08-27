package postgres

import (
	"context"
	"errors"
	"fmt"
	"payment_gateway/internal/storage/postgres/sqlc"

	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrEmptyDatabaseURL = errors.New("empty database URL")

type Storage struct {
	pool    *pgxpool.Pool
	queries *sqlc.Queries
}

func SetupDatabase(ctx context.Context, url string, maxConns int32) (*Storage, error) {

	if url == "" {
		return nil, ErrEmptyDatabaseURL
	}

	cfg, err := pgxpool.ParseConfig(url)
	if err != nil {
		return nil, fmt.Errorf("parse database url: %w", err)
	}
	if maxConns > 0 {
		cfg.MaxConns = maxConns
	}

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("open pool: %w", err)
	}
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("ping: %w", err)
	}

	return &Storage{pool: pool, queries: sqlc.New(pool)}, nil
}

func (s *Storage) Close() {
	if s == nil || s.pool == nil {
		return
	}
	s.pool.Close()
}
