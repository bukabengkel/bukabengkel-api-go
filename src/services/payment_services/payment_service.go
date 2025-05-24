package payment_services

import (
	"fmt"

	"github.com/peang/bukabengkel-api-go/src/config"
)

type PaymentCustomerDetail struct {
	CustomerName string
	CustomerEmail string
	CustomerPhone string
}

type PaymentOrderDetail struct {
	OrderID string
	OrderAmount float64
	OrderCurrency string
}

type PaymentService interface {
	CreatePayment(customerDetail PaymentCustomerDetail, orderDetail PaymentOrderDetail) (string, error)
}

const (
	PAYMENT_SERVICE_MIDTRANS = "midtrans"
)

func NewPaymentService(config *config.Config) (PaymentService, error) {
	switch config.PaymentProvider.PaymentProviderName {
	case PAYMENT_SERVICE_MIDTRANS:
		return NewMidtransService(config), nil
	default:
		return nil, fmt.Errorf("payment provider %s not supported", config.PaymentProvider.PaymentProviderName)
	}
}
