package response

import (
	"time"

	"github.com/peang/bukabengkel-api-go/src/models"
)

type ProductDistributorResponse struct {
	ID            string                    `json:"id"`
	Name          string                    `json:"name"`
	Description   string                    `json:"description"`
	Distributor   string                    `json:"distributor"`
	DistributorID string                    `json:"distributor_id"`
	Code          string                    `json:"code"`
	Category      string                    `json:"category"`
	Unit          string                    `json:"unit"`
	Thumbnail     string                    `json:"thumbnail"`
	Price         float64                   `json:"price"`
	BulkPrice     []models.ProductBulkPrice `json:"bulk_price"`
	Weight        float64                   `json:"weight"`
	Volume        float64                   `json:"volume"`
	Stock         float64                   `json:"stock"`
	UpdatedAt     time.Time                 `json:"updated_at"`
}

func ProductDistributorDetailResponse(product *models.ProductDistributor) *ProductDistributorResponse {
	response := &ProductDistributorResponse{
		ID:            product.Key,
		Name:          product.Name,
		Description:   product.Description,
		Distributor:   product.Distributor.Name,
		DistributorID: product.Distributor.Key,
		Code:          product.Code,
		Category:      product.Category.Name,
		Unit:          product.Unit,
		Thumbnail:     product.ThumbnailCDN,
		Price:         product.Price,
		BulkPrice:     product.BulkPrice,
		Weight:        product.Weight,
		Volume:        product.Volume,
		Stock:         product.Stock,
		UpdatedAt:     product.UpdatedAt,
	}

	return response
}

func ProductDistributorListResponse(products *[]models.ProductDistributor) []ProductDistributorResponse {
	var responses = make([]ProductDistributorResponse, 0)
	for _, product := range *products {
		response := &ProductDistributorResponse{
			ID:            product.Key,
			Name:          product.Name,
			Distributor:   product.Distributor.Name,
			DistributorID: product.Distributor.Key,
			Code:          product.Code,
			Category:      product.Category.Name,
			Unit:          product.Unit,
			Thumbnail:     product.ThumbnailCDN,
			Price:         product.Price,
			BulkPrice:     product.BulkPrice,
			UpdatedAt:     product.UpdatedAt,
		}
		responses = append(responses, *response)
	}
	return responses
}
