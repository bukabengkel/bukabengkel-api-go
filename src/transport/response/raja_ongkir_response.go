package response

type RajaOngkirLocation struct {
	ID              int    `json:"id"`
	Label           string `json:"label"`
	ProvinceName    string `json:"province_name"`
	CityName        string `json:"city_name"`
	DistrictName    string `json:"district_name"`
	SubdistrictName string `json:"subdistrict_name"`
	ZipCode         string `json:"zip_code"`
}

type RajaOngkirGetLocationResponse struct {
	Meta struct {
		Message string `json:"message"`
		Code    int    `json:"code"`
		Status  string `json:"status"`
	} `json:"meta"`
	Data []RajaOngkirLocation `json:"data"`
}

type RajaOngkirShippingRate struct {
	Name              string `json:"name"`
	Code              string `json:"code"`
	Service           string `json:"service"`
	Description       string `json:"description"`
	Cost              int    `json:"cost"`
	ETD               string `json:"etd"`
}

type RajaOngkirGetShippingRateResponse struct {
	Meta struct {
		Message string `json:"message"`
		Code    int    `json:"code"`
		Status  string `json:"status"`
	} `json:"meta"`
	Data []RajaOngkirShippingRate `json:"data"`
}
