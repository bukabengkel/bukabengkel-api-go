package models

type CartType uint8

const (
	CartTypePOS CartType = iota
	CartTypeShop
)

type CartStatus uint8

const (
	CartStatusDraft CartStatus = iota
	CartStatusInProgress
)

type CartPaymentMethod uint8

const (
	CartPaymentMethodCash CartPaymentMethod = iota
	CartPaymentMethodInProgress
)

type CartShopping struct {
	ID              *uint64
	DistributorKey  string
	DistributorName string
	ProductKey      string
	ProductName     string
	ProductUnit     string
	ProductImage    string
	Qty             uint64
	BasePrice       float64
	BulkPrice       []struct {
		Qty   uint64
		Price float64
	}
	Price    float64
	Discount float64
	Weight   float64
	Volume   float64
}

type Cart struct {
	ID                 *uint64
	StoreID            uint
	Items              []CartShopping
	TotalCartSellPrice float64
	TotalCartDiscount  float64
	TotalCart          float64
	Status             CartStatus
}
