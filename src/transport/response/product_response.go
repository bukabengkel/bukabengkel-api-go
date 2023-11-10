package response

import (
	"github.com/peang/bukabengkel-api-go/src/domain/entity"
)

// BuildingResponse represents building response payload
type ProductResponse struct {
	ID               string  `json:"id"`
	CategoryID       int     `json:"categoryId"`
	CategoryName     string  `json:"categoryName"`
	Name             string  `json:"name"`
	Slug             string  `json:"slug"`
	Description      string  `json:"description"`
	Unit             string  `json:"unit"`
	Thumbnail        string  `json:"thumbnail"`
	Price            float64 `json:"price"`
	SellPrice        float64 `json:"sellPrice"`
	Stock            int     `json:"stock"`
	StockMinimum     int     `json:"stockMinimum"`
	IsStockUnlimited bool    `json:"isStockUnlimited"`
	Status           int     `json:"status"`
	StatusString     string  `json:"statusString"`
}

// BuildingDetailResponse transforms entity.Building to BuildingResponse
func ProductDetailResponse(product *entity.Product) *ProductResponse {
	return &ProductResponse{
		ID:               product.Key,
		CategoryID:       product.Category.ID,
		CategoryName:     product.Category.Name,
		Name:             product.Name,
		Slug:             product.Slug,
		Description:      product.Description,
		Unit:             product.Unit,
		Thumbnail:        product.Thumbnail.Path,
		Price:            product.Price,
		SellPrice:        product.SellPrice,
		Stock:            product.Stock,
		StockMinimum:     product.StockMinimum,
		IsStockUnlimited: product.IsStockUnlimited,
		Status:           int(product.Status),
		StatusString:     product.Status.String(),
	}
}

// BuildingListResponse transforms []entity.Building to []BuildingResponse
func ProductListResponse(products []*entity.Product) []ProductResponse {
	var responses = make([]ProductResponse, 0)
	for _, building := range products {
		response := ProductDetailResponse(building)
		responses = append(responses, *response)
	}
	return responses
}
