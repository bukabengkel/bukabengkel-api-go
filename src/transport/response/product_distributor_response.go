package response

import (
	"github.com/peang/bukabengkel-api-go/src/models"
)

type ProductDistributorResponse struct {
	ID          string                    `json:"id"`
	Name        string                    `json:"name"`
	Distributor string                    `json:"distributor"`
	Code        string                    `json:"code"`
	Category    string                    `json:"category"`
	Unit        string                    `json:"unit"`
	Thumbnail   string                    `json:"thumbnail"`
	Price       float64                   `json:"price"`
	BulkPrice   []models.ProductBulkPrice `json:"bulk_price"`
	Stock       float64                   `json:"stock"`
}

func ProductDistributorDetailResponse(product *models.ProductDistributor) *ProductDistributorResponse {
	response := &ProductDistributorResponse{
		ID:          product.Key,
		Name:        product.Name,
		Distributor: product.Distributor.Name,
		Code:        product.Code,
		Category:    product.Category.Name,
		Unit:        product.Unit,
		Thumbnail:   product.Thumbnail,
		Price:       product.Price,
		BulkPrice:   product.BulkPrice,
		Stock:       product.Stock,
	}

	return response
}

func ProductDistributorListResponse(products *[]models.ProductDistributor) []ProductDistributorResponse {
	var responses = make([]ProductDistributorResponse, 0)
	for _, product := range *products {
		response := ProductDistributorDetailResponse(&product)
		responses = append(responses, *response)
	}
	return responses
}
