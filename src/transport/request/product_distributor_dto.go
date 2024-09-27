package request

type ProductDistributorListDTO struct {
	Page          string
	PerPage       string
	Sort          string
	Keyword       string
	DistributorID string
	Name          string
	Code          string
	RemoteUpdate  string
}

type ProductDistributorDetailDTO struct {
	ID string
}
