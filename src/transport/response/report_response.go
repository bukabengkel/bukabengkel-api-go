package response

import (
	usecase "github.com/peang/bukabengkel-api-go/src/usecases"
)

type orderSalesReportResponse struct {
	TotalSales   float32 `json:"total_sales"`
	TotalNett    float32 `json:"total_nett"`
	TotalProduct int     `json:"total_product"`
}

type productSalesReportResponse struct {
	ProductKey      string  `json:"product_key"`
	ProductName     string  `json:"product_name"`
	QtySales        int     `json:"qty_sales"`
	QtyCurrentStock float64 `json:"qty_current_stock"`
}

func OrderSalesReportResponse(reports *usecase.SalesOrderResult) *orderSalesReportResponse {
	return &orderSalesReportResponse{
		TotalSales:   reports.TotalSales,
		TotalNett:    reports.TotalNett,
		TotalProduct: reports.TotalProduct,
	}
}

func ProductSalesReportResponse(reports *[]usecase.ProductOrderResult) *[]productSalesReportResponse {
	var result []productSalesReportResponse
	for _, item := range *reports {
		result = append(result, productSalesReportResponse{
			ProductKey:      item.ProductKey,
			ProductName:     item.ProductName,
			QtySales:        item.QtySales,
			QtyCurrentStock: item.QtyStock,
		})
	}
	return &result
}
