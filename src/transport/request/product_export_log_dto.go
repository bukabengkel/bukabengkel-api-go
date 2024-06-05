package request

type ProductExportLogListDTO struct {
	StoreID uint
	UserID  uint
	Page    string
	PerPage string
	Sort    string
}
