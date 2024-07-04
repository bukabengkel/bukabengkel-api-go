package response

import "github.com/peang/bukabengkel-api-go/src/models"

type orderSalesReportResponse struct {
	TotalSales   float32 `json:"total_sales"`
	TotalNett    float32 `json:"total_nett"`
	TotalProduct int     `json:"total_product"`
}

type productSalesReportResponse struct {
	ProductName     string `json:"product_name"`
	QtySales        int    `json:"qty_sales"`
	QtyCurrentStock int    `json:"qty_current_stock"`
}

func OrderSalesReportResponse(totalSales float32, totalNett float32, totalProduct int) *orderSalesReportResponse {
	return &orderSalesReportResponse{
		TotalSales:   totalSales,
		TotalNett:    totalNett,
		TotalProduct: totalProduct,
	}
}

func ProductSalesReportResponse(*[]models.Product) *[]productSalesReportResponse {
	return nil
}
