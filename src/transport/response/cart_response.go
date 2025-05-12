package response

import (
	"github.com/peang/bukabengkel-api-go/src/models"
)

type cartShippingRateDistributorLocation struct {
	Province    string `json:"province"`
	City        string `json:"city"`
	District    string `json:"district"`
	Subdistrict string `json:"subdistrict"`
	ZipCode     string `json:"zip_code"`
}

type cartShippingRateDistributor struct {
	Name     string                              `json:"name"`
	Location cartShippingRateDistributorLocation `json:"location"`
}

type cartShippingRateBulkPrice struct {
	Qty   uint64  `json:"qty"`
	Price float64 `json:"price"`
}

type cartShippingRateItems struct {
	ID           *uint64                     `json:"id"`
	ProductCode  string                      `json:"productCode"`
	ProductName  string                      `json:"productName"`
	ProductUnit  string                      `json:"productUnit"`
	ProductImage string                      `json:"productImage"`
	Qty          uint64                      `json:"qty"`
	BasePrice    float64                     `json:"basePrice"`
	BulkPrice    []cartShippingRateBulkPrice `json:"bulkPrice"`
	Price        float64                     `json:"price"`
	Discount     float64                     `json:"discount"`
	Weight       float64                     `json:"weight"`
	Volume       float64                     `json:"volume"`
}

type cartShippingRateResponse struct {
	Distributor  cartShippingRateDistributor `json:"distributor"`
	Items        []cartShippingRateItems     `json:"items"`
	ShippingRate any                         `json:"shippingRate"`
}

func CartShippingRateResponse(
	distributor *models.Distributor,
	distributorLocation *models.RajaOngkirLocation,
	cartItems *[]models.CartShopping,
	shippingCost *any,
) *cartShippingRateResponse {

	items := make([]cartShippingRateItems, 0)
	for _, item := range *cartItems {

		var bulkPrice []cartShippingRateBulkPrice
		for _, bulk := range item.BulkPrice {
			bulkPrice = append(bulkPrice, cartShippingRateBulkPrice{
				Qty:   bulk.Qty,
				Price: bulk.Price,
			})
		}

		items = append(items, cartShippingRateItems{
			ID:           item.ID,
			ProductCode:  item.ProductKey,
			ProductName:  item.ProductName,
			ProductUnit:  item.ProductUnit,
			ProductImage: item.ProductImage,
			Qty:          item.Qty,
			BasePrice:    item.BasePrice,
			BulkPrice:    bulkPrice,
			Price:        item.Price,
			Discount:     item.Discount,
			Weight:       item.Weight,
			Volume:       item.Volume,
		})
	}

	response := &cartShippingRateResponse{
		Distributor: cartShippingRateDistributor{
			Name: distributor.Name,
			Location: cartShippingRateDistributorLocation{
				Province:    distributorLocation.ProvinceName,
				City:        distributorLocation.CityName,
				District:    distributorLocation.DistrictName,
				Subdistrict: distributorLocation.SubdistrictName,
				ZipCode:     distributorLocation.ZipCode,
			},
		},
		Items:        items,
		ShippingRate: shippingCost,
	}

	return response
}
