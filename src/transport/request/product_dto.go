package request

type ProductListDTO struct {
	StoreID       uint64
	Page          string
	PerPage       string
	Sort          string
	Keyword       string
	Name          string
	CategoryId    string
	StockMoreThan string
}

type ProductExportDTO struct {
	StoreID    uint64
	UserID     uint64
	CategoryId *uint64
}
