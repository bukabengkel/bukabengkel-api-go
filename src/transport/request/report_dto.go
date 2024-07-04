package request

type OrderSalesReportDTO struct {
	StoreID   uint
	StartDate string
	EndDate   string
}

type ProductSalesRxeportDTO struct {
	Page      string
	PerPage   string
	StoreID   uint
	StartDate string
	EndDate   string
}
