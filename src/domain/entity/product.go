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
	ID               int
	Key              string
	Store            Store
	Category         ProductCategory
	Name             string
	Slug             string
	Description      string
	Unit             string
	Thumbnail        Image
	Images           []Image
	Price            float64
	SellPrice        float64
	Stock            int
	StockReserved    int
	StockValue       int
	StockMinimum     int
	IsStockUnlimited bool
	Status           ProductStatus
	CreatedAt        time.Time
	CreatedBy        string
	UpdatedAt        time.Time
	UpdatedBy        string
	DeletedAt        time.Time
	DeletedBy        string
}
