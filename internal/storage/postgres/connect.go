package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

func Connect(dbURL string) (*pgxpool.Pool, error) {
	return pgxpool.New(context.Background(), dbURL)
}
