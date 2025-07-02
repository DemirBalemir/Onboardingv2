package db

import (
	"context"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPostgresConnection(dsn string) *pgxpool.Pool {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}

	err = pool.Ping(ctx)
	if err != nil {
		log.Fatalf("Unable to ping database: %v", err)
	}

	log.Println("Successfully connected to the database")
	return pool
}
