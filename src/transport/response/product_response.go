package response

import (
	"github.com/peang/bukabengkel-api-go/src/models"
)

type productDetailResponse struct {
	ID               string  `json:"id"`
	Store            string  `json:"store"`
	BrandID          *uint64 `json:"brandId"`
	BrandName        *string `json:"brandName"`
	CategoryID       uint64  `json:"categoryId"`
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

func ProductDetailResponse(product *models.Product) *productDetailResponse {
	response := &productDetailResponse{
		ID:               product.Key,
		Store:            product.Store.Name,
		Name:             product.Name,
		Slug:             product.Slug,
		Description:      product.Description,
		Unit:             product.Unit,
		Thumbnail:        "",
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
		response.BrandID = product.Brand.ID
		response.BrandName = &product.Brand.Name
	}

	if product.Category == nil {
		response.CategoryID = uint64(0)
		response.CategoryName = ""
	} else {
		response.CategoryID = product.CategoryID
		response.CategoryName = product.Category.Name
	}

	if len(product.Images) > 0 {
		response.Thumbnail = product.Images[0].Path
	}

	return response
}

func ProductListResponse(products *[]models.Product) []productDetailResponse {
	if products == nil {
		return []productDetailResponse{}
	}

	var responses = make([]productDetailResponse, 0)
	for _, product := range *products {
		response := ProductDetailResponse(&product)
		responses = append(responses, *response)
	}
	return responses
}
