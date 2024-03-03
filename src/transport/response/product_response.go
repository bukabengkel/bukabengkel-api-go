package response

import (
	"github.com/peang/bukabengkel-api-go/src/domain/entity"
)

// BuildingResponse represents building response payload
type ProductResponse struct {
	ID               string  `json:"id"`
	Store            string  `json:"store"`
	BrandID          *int64  `json:"brandId"`
	BrandName        *string `json:"brandName"`
	CategoryID       int64   `json:"categoryId"`
	CategoryName     string  `json:"categoryName"`
	Name             string  `json:"name"`
	Slug             string  `json:"slug"`
	Description      string  `json:"description"`
	Unit             string  `json:"unit"`
	Thumbnail        string  `json:"thumbnail"`
	Price            float64 `json:"price"`
	SellPrice        float64 `json:"sellPrice"`
	Stock            float64 `json:"stock"`
	StockMinimum     float64 `json:"stockMinimum"`
	IsStockUnlimited bool    `json:"isStockUnlimited"`
	Status           int     `json:"status"`
	StatusString     string  `json:"statusString"`
}

// BuildingDetailResponse transforms entity.Building to BuildingResponse
func ProductDetailResponse(product *entity.Product) *ProductResponse {
	response := &ProductResponse{
		ID:               product.Key,
		Store:            product.Store.Name,
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

	if product.Brand == nil {
		response.BrandID = nil
		response.BrandName = nil
	} else {
		response.BrandID = &product.Brand.ID
		response.BrandName = &product.Brand.Name
	}

	return response
}

// BuildingListResponse transforms []entity.Building to []BuildingResponse
func ProductListResponse(products *[]entity.Product) []ProductResponse {
	var responses = make([]ProductResponse, 0)
	for _, product := range *products {
		response := ProductDetailResponse(&product)
		responses = append(responses, *response)
	}
	return responses
}
