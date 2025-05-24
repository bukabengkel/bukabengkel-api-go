package request

type OrderDistributorListDTO struct {
	Page    string
	PerPage string
	Sort    string
	StoreID uint
	OrderID string
}

type OrderDistributorDetailDTO struct {
	StoreID uint
	OrderID string
}
