package models

import (
	"time"

	"github.com/uptrace/bun"
)

type Location struct {
	bun.BaseModel `bun:"table:location"`

	ID        *uint64   `bun:"id,pk"`
	Province  string    `bun:"province,notnull"`
	City      string    `bun:"city,notnull"`
	District  string    `bun:"district,notnull"`
	Urban     string    `bun:"urban,notnull"`
	CreatedAt time.Time `bun:"created_at"`
	UpdatedAt time.Time `bun:"updated_at"`
}
