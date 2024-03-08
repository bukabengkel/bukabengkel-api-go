package models

import (
	"encoding/json"

	// Adjust for your database dialect
	"time"

	"github.com/uptrace/bun"
)

type ProductDistributor struct {
	bun.BaseModel `bun:"table:product_distributor,timestamps"`

	ID               uint64           `bun:"id,pk,autoIncrement"`
	ExternalID       string           `bun:"external_id"`
	Key              string           `bun:"key"`
	DistributorID    uint64           `bun:"distributor_id,fk:distributor_distributor_id"`
	CategoryID       uint64           `bun:"category_id,fk:product_category_distributor_category_id"`
	Name             string           `bun:"name"`
	Code             string           `bun:"code"`
	Description      string           `bun:"description"`
	Unit             string           `bun:"unit"`
	Thumbnail        string           `bun:"thumbnail"`
	Images           []string         `bun:"images"`
	Price            float64          `bun:"price"`
	BulkPrice        ProductBulkPrice `bun:"type:jsonb,column:bulk_price"`
	Weight           float64          `bun:"weight"`
	Volume           float64          `bun:"volume"`
	Stock            uint64           `bun:"stock"`
	IsStockUnlimited bool             `bun:"is_stock_unlimited"`
	Status           uint8            `bun:"status"`
	CreatedAt        time.Time        `bun:"created_at"`
	UpdatedAt        time.Time        `bun:"updated_at"`
	RemoteUpdate     bool             `bun:"remote_update"`

	Category    *ProductCategoryDistributor `bun:"rel:hasOne,to:product_category_distributor,foreignKey:category_id"`
	Distributor *Distributor                `bun:"rel:hasOne,to:distributor,foreignKey:distributor_id"`
}

type ProductBulkPrice map[string]interface{}

func (p *ProductDistributor) UnmarshalBulkPrice(data []byte) error {
	return json.Unmarshal(data, &p.BulkPrice)
}

func (p *ProductDistributor) MarshalBulkPrice() ([]byte, error) {
	return json.Marshal(p.BulkPrice)
}
