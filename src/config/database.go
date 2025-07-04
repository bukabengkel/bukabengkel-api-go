package config

import (
	"database/sql"
	"runtime"
	"time"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bundebug"
)

func LoadDatabase(c *Config) *bun.DB {
	pgconn := pgdriver.NewConnector(
		pgdriver.WithDSN(c.DatabaseURL),
	)

	sqldb := sql.OpenDB(pgconn)
	if c.Env == "production" {
		maxOpenConns := 4 * runtime.GOMAXPROCS(0)
		sqldb.SetMaxOpenConns(maxOpenConns)
		sqldb.SetMaxIdleConns(maxOpenConns)
		sqldb.SetConnMaxLifetime(time.Minute * 5)
	}

	db := bun.NewDB(sqldb, pgdialect.New())

	if c.Env != "production" {
		db.AddQueryHook(bundebug.NewQueryHook(
			bundebug.WithVerbose(true),
			bundebug.FromEnv("BUNDEBUG"),
		))
	}

	return db
}
