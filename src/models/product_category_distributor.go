package models

import (
	"time"

	"github.com/uptrace/bun"
)

type ProductCategoryDistributor struct {
	bun.BaseModel `bun:"table:product_category_distributor"`

	ID            *uint64   `bun:"id,pk"`
	ExternalID    string    `bun:"external_id"`
	DistributorID uint64    `bun:"distributor_id"`
	Name          string    `bun:"name"`
	Code          string    `bun:"code"`
	Description   string    `bun:"description"`
	CreatedAt     time.Time `bun:"created_at"`
	UpdatedAt     time.Time `bun:"updated_at"`
	RemoteUpdate  bool      `bun:"remote_update"`

	Distributor *Distributor `bun:"rel:belongs-to"`
}
