package admin_response

type AdminDashboardResponse struct {
	TotalOrderThisMonth int `json:"total_orders_this_month"`
	TotalOrderLastMonth int `json:"total_orders_last_month"`
	TotalProductThisMonth int `json:"total_products_this_month"`
	TotalProductLastMonth int `json:"total_products_last_month"`
	TotalStoreThisMonth int `json:"total_stores_this_month"`
	TotalStoreLastMonth int `json:"total_stores_last_month"`
	TotalOrderAmountThisMonth float32 `json:"total_order_amount_this_month"`
	TotalOrderAmountLastMonth float32 `json:"total_order_amount_last_month"`
}