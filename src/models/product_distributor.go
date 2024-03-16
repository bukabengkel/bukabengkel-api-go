package models

import (

	// Adjust for your database dialect
	"time"

	"github.com/uptrace/bun"
)

type ProductDistributor struct {
	bun.BaseModel `bun:"table:product_distributor"`

	ID               *uint64            `bun:"id,pk"`
	ExternalID       string             `bun:"external_id"`
	Key              string             `bun:"key"`
	DistributorID    uint64             `bun:"distributor_id"`
	CategoryID       uint64             `bun:"category_id"`
	Name             string             `bun:"name"`
	Code             string             `bun:"code"`
	Description      string             `bun:"description"`
	Unit             string             `bun:"unit"`
	Thumbnail        string             `bun:"thumbnail"`
	Images           []string           `bun:"images,array"`
	Price            float64            `bun:"price"`
	BulkPrice        []ProductBulkPrice `bun:"type:jsonb,column:bulk_price"`
	Weight           float64            `bun:"weight"`
	Volume           float64            `bun:"volume"`
	Stock            float64            `bun:"stock"`
	IsStockUnlimited bool               `bun:"is_stock_unlimited"`
	Status           ProductStatus      `bun:"status"`
	CreatedAt        time.Time          `bun:"created_at"`
	UpdatedAt        time.Time          `bun:"updated_at"`
	RemoteUpdate     bool               `bun:"remote_update"`

	Category    *ProductCategoryDistributor `bun:"rel:belongs-to"`
	Distributor *Distributor                `bun:"rel:belongs-to"`
}

type ProductBulkPrice struct {
	Qty   int64
	Price float32
}
