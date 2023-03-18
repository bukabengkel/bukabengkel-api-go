package config

import (
	"context"
	"log"

	"github.com/jackc/pgx/v4/pgxpool"
)

// NewPgxPool returns a new connection pool to the provided PostgreSQL database
func NewPgxDatabase(databaseURL string) (db *pgxpool.Pool, err error) {
	config, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		panic(err)
	}

	pool, err := pgxpool.ConnectConfig(context.Background(), config)
	if err != nil {
		log.Fatalf("error connecting to database: %v", err)
	}

	return pool, err
}
