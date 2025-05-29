package request

type MidtransWebhookDTO struct {
	OrderID           string `json:"order_id"`
	TransactionID     string `json:"transaction_id"`
	SettlementTime    string `json:"settlement_time"`
	PaymentType       string `json:"payment_type"`
	GrossAmount       string `json:"gross_amount"`
	FraudStatus       string `json:"fraud_status"`
	
}
