package response

type cartShoppingCheckoutResponse struct {
	OrderID     string  `json:"order_id"`
	OrderAmount float64 `json:"order_amount"`
	PaymentURL  string  `json:"payment_url"`
}

func CartShoppingCheckoutResponse(orderID string, orderAmount float64, paymentURL string) *cartShoppingCheckoutResponse {
	response := &cartShoppingCheckoutResponse{
		OrderID:     orderID,
		OrderAmount: orderAmount,
		PaymentURL:  paymentURL,
	}

	return response
}