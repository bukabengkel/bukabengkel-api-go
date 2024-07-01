package response

type reportSalesResponse struct {
	TotalSales   float32 `json:"total_sales"`
	TotalNett    float32 `json:"total_nett"`
	TotalProduct int     `json:"total_product"`
}

func ReportSalesResponse(totalSales float32, totalNett float32, totalProduct int) *reportSalesResponse {
	return &reportSalesResponse{
		TotalSales:   totalSales,
		TotalNett:    totalNett,
		TotalProduct: totalProduct,
	}
}
