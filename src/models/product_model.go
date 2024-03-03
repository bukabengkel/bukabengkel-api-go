package models

import (
	"time"

	"github.com/peang/bukabengkel-api-go/src/domain/entity"
	"github.com/uptrace/bun"
)

type Product struct {
	bun.BaseModel `bun:"table:product"`

	ID               *int64    `bun:"id,pk"`
	Key              string    `bun:"key,notnull,unique"`
	StoreID          int64     `bun:"store_id,notnull"`
	CategoryID       int64     `bun:"category_id,notnull"`
	Name             string    `bun:"name,notnull"`
	Slug             string    `bun:"slug,notnull,unique"`
	Description      string    `bun:"description,notnull"`
	Unit             string    `bun:"unit,notnull"`
	Price            float64   `bun:"price,notnull"`
	SellPrice        float64   `bun:"sell_price,notnull"`
	Stock            float64   `bun:"stock,notnull"`
	StockReserved    float64   `bun:"stock_reserved,notnull"`
	StockValue       float64   `bun:"stock_value,notnull"`
	StockMinimum     float64   `bun:"stock_minimum,notnull"`
	IsStockUnlimited bool      `bun:"is_stock_unlimited,notnull"`
	Status           int       `bun:"status,notnull"`
	CreatedAt        time.Time `bun:"created_at"`
	CreatedBy        string    `bun:"created_by"`
	UpdatedAt        time.Time `bun:"updated_at"`
	UpdatedBy        string    `bun:"updated_by"`
	DeletedAt        time.Time `bun:"deleted_at"`
	DeletedBy        string    `bun:"deleted_by"`

	Images   []Image          `bun:"rel:has-many,join:id=entity_id,join:type=entity_type,polymorphic"`
	Store    *Store           `bun:"rel:belongs-to"`
	Brand    *ProductBrand    `bun:"rel:belongs-to"`
	Category *ProductCategory `bun:"rel:belongs-to"`
}

func LoadProductModel(p Product) *entity.Product {
	return &entity.Product{
		ID:               *p.ID,
		Key:              p.Key,
		Store:            LoadStoreModel(p.Store),
		Brand:            LoadProductBrandModel(p.Brand),
		Category:         LoadProductCategoryModel(p.Category),
		Name:             p.Name,
		Slug:             p.Slug,
		Description:      p.Description,
		Unit:             p.Unit,
		Thumbnail:        &entity.Image{},
		Price:            p.Price,
		SellPrice:        p.SellPrice,
		Stock:            p.Stock,
		StockReserved:    p.StockReserved,
		StockValue:       p.StockValue,
		StockMinimum:     p.StockMinimum,
		IsStockUnlimited: p.IsStockUnlimited,
		Status:           entity.ProductStatus(p.Status),
		CreatedAt:        p.CreatedAt,
		CreatedBy:        p.CreatedBy,
		UpdatedAt:        p.UpdatedAt,
		UpdatedBy:        p.UpdatedBy,
		DeletedAt:        p.DeletedAt,
		DeletedBy:        p.DeletedBy,
	}
}
