package entity

import "time"

type ProductStatus int

const (
	Inactive ProductStatus = iota
	Active
)

func (s ProductStatus) String() string {
	switch s {
	case Inactive:
		return "inactive"
	case Active:
		return "active"
	default:
		return "unknown"
	}
}

type ProductEntity struct {
	ID               int
	Key              string
	Store            StoreEntity
	Category         ProductCategoryEntity
	Name             string
	Slug             string
	Description      string
	Unit             string
	Thumbnail        ImageEntity
	Images           []ImageEntity
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

type ProductEntityRepositoryFilter struct {
	Name string
}
