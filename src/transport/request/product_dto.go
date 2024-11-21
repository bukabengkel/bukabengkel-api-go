package request

type ProductListDTO struct {
	StoreID       uint
	Page          string
	PerPage       string
	Sort          string
	Keyword       string
	Name          string
	CategoryId    string
	StockMoreThan string
	Status        string
	Cursor        string
}

type ProductListDTOV2 struct {
	StoreID       uint
	Limit         string
	Sort          string
	Keyword       string
	Name          string
	CategoryId    string
	StockMoreThan string
	Status        string
	Cursor        string
}
