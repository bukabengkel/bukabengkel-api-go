package models

import (
	"time"

	"github.com/uptrace/bun"
)

type ProductCategory struct {
	bun.BaseModel `bun:"table:product_category"`

	ID          *int64    `bun:"id,pk"`
	StoreID     int64     `bun:"store_id,notnull"`
	Name        string    `bun:"name,notnull"`
	Description string    `bun:"description"`
	CreatedAt   time.Time `bun:"created_at,notnull"`
	UpdatedAt   time.Time `bun:"updated_at"`
}
