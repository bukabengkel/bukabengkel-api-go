package response


type orderSalesReportResponse struct {
	TotalSales   float32 `json:"total_sales"`
	TotalNett    float32 `json:"total_nett"`
	TotalProduct int     `json:"total_product"`
}

type productSalesReportResponse struct {
	ProductKey      string  `json:"product_key"`
	ProductName     string  `json:"product_name"`
	ProductCategory string  `json:"product_category"`
	ProductUnit     string  `json:"product_unit"`
	QtySales        int     `json:"qty_sales"`
	QtyCurrentStock float64 `json:"qty_current_stock"`
}

type SalesOrderResult struct {
	TotalSales   float32
	TotalNett    float32
	TotalProduct int
}

type ProductOrderResult struct {
	ProductKey      string
	ProductName     string
	ProductCategory string
	ProductUnit     string
	QtySales        int
	QtyStock        float64
}

func OrderSalesReportResponse(reports *SalesOrderResult) *orderSalesReportResponse {
	return &orderSalesReportResponse{
		TotalSales:   reports.TotalSales,
		TotalNett:    reports.TotalNett,
		TotalProduct: reports.TotalProduct,
	}
}

func ProductSalesReportResponse(reports *[]ProductOrderResult) *[]productSalesReportResponse {
	result := []productSalesReportResponse{}

	for _, item := range *reports {
		result = append(result, productSalesReportResponse{
			ProductKey:      item.ProductKey,
			ProductName:     item.ProductName,
			ProductCategory: item.ProductCategory,
			ProductUnit:     item.ProductUnit,
			QtySales:        item.QtySales,
			QtyCurrentStock: item.QtyStock,
		})
	}

	return &result
}
