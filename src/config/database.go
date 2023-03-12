package config

import (
	"context"

	"github.com/jackc/pgx/v4"
)

// NewPgxPool returns a new connection pool to the provided PostgreSQL database
func NewPgxDatabase(databaseURL string) (db *pgx.Conn, err error) {
	config, err := pgx.ParseConfig(databaseURL)
	if err != nil {
		panic(err)
	}
	conn, err := pgx.ConnectConfig(context.Background(), config)
	if err != nil {
		panic(err)
	}

	return conn, err
}
