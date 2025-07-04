package request

type CartShoppingGetShippingRateDTO struct {
	StoreID       uint
	UserID        uint
	DistributorID string
}

type CartShoppingCheckoutDTO struct {
	StoreID                 uint
	UserID                  uint
	DistributorID           string
	ShippingProviderService string `json:"shipping_provider_service" validate:"required"`
	ShippingProviderCode    string `json:"shipping_provider_code" validate:"required"`
	ShippingProviderRemarks string `json:"shipping_provider_remarks" validate:""`
	StoreRemarks            string `json:"store_remarks" validate:""`
}
