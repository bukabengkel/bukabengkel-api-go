package request

type MidtransWebhookDTO struct {
	OrderID           string  `json:"order_id"`
	TransactionID     string  `json:"transaction_id"`
	TransactionStatus string  `json:"transaction_status"`
	PaymentType       string  `json:"payment_type"`
	GrossAmount       string  `json:"gross_amount"`
	FraudStatus       string  `json:"fraud_status"`
	ExpiredAt         *string `json:"expired_at"`
}
