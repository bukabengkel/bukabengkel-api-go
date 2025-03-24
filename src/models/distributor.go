package models

import (
	"time"

	"github.com/uptrace/bun"
)

type Distributor struct {
	bun.BaseModel `bun:"table:distributor"`

	ID             uint64    `bun:"id,pk"`
	Key            string    `bun:"key"`
	Name           string    `bun:"name"`
	LocationDetail string    `bun:"location_detail"`
	CreatedAt      time.Time `bun:"created_at"`
	UpdatedAt      time.Time `bun:"updated_at"`
}
