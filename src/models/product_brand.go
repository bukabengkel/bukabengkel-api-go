package models

import (
	"time"

	"github.com/uptrace/bun"
)

type ProductBrand struct {
	bun.BaseModel `bun:"table:product_brand"`

	ID          *uint64   `bun:"id,pk"`
	StoreID     uint64    `bun:"store_id,notnull"`
	Name        string    `bun:"name,notnull"`
	Description string    `bun:"description"`
	CreatedAt   time.Time `bun:"created_at,notnull"`
	UpdatedAt   time.Time `bun:"updated_at"`
}
