package request

type CartGetShippingRateDTO struct {
	StoreID      uint64
	UserID       uint64
	DistributorID string	
}

type CartCheckoutDTO struct {
	StoreID      uint64
	UserID       uint64
	DistributorID string	
}