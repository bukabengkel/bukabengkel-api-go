package models

import (
	"time"

	"github.com/uptrace/bun"
)

type ProductStatus int

const (
	ProductInactive ProductStatus = iota
	ProductActive
)

func (s ProductStatus) String() string {
	switch s {
	case ProductInactive:
		return "inactive"
	case ProductActive:
		return "active"
	default:
		return "unknown"
	}
}

type Product struct {
	bun.BaseModel `bun:"table:product"`

	ID               *uint64       `bun:"id,pk"`
	Key              string        `bun:"key,notnull,unique"`
	StoreID          uint64        `bun:"store_id,notnull"`
	CategoryID       uint64        `bun:"category_id,notnull"`
	Name             string        `bun:"name,notnull"`
	Slug             string        `bun:"slug,notnull,unique"`
	Description      string        `bun:"description,notnull"`
	Unit             string        `bun:"unit,notnull"`
	Price            float64       `bun:"price,notnull"`
	SellPrice        float64       `bun:"sell_price,notnull"`
	Stock            float64       `bun:"stock,notnull"`
	StockReserved    float64       `bun:"stock_reserved,notnull"`
	StockValue       float64       `bun:"stock_value,notnull"`
	StockMinimum     float64       `bun:"stock_minimum,notnull"`
	IsStockUnlimited bool          `bun:"is_stock_unlimited,notnull"`
	Status           ProductStatus `bun:"status,notnull"`
	CreatedAt        time.Time     `bun:"created_at"`
	CreatedBy        string        `bun:"created_by"`
	UpdatedAt        time.Time     `bun:"updated_at"`
	UpdatedBy        string        `bun:"updated_by"`
	DeletedAt        time.Time     `bun:"deleted_at"`
	DeletedBy        string        `bun:"deleted_by"`

	Images   []Image          `bun:"rel:has-many,join:id=entity_id,join:type=entity_type,polymorphic"`
	Store    *Store           `bun:"rel:belongs-to"`
	Brand    *ProductBrand    `bun:"rel:belongs-to"`
	Category *ProductCategory `bun:"rel:belongs-to"`
}
