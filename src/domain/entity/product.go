package entity

import "time"

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
	ID               int64
	Key              string
	Store            *Store
	Brand            *ProductBrand
	Category         *ProductCategory
	Name             string
	Slug             string
	Description      string
	Unit             string
	Thumbnail        *Image
	Images           []*Image
	Price            float64
	SellPrice        float64
	Stock            float64
	StockReserved    float64
	StockValue       float64
	StockMinimum     float64
	IsStockUnlimited bool
	Status           ProductStatus
	CreatedAt        time.Time
	CreatedBy        string
	UpdatedAt        time.Time
	UpdatedBy        string
	DeletedAt        time.Time
	DeletedBy        string
}
