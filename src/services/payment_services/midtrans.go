package payment_services

import (
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
	"github.com/peang/bukabengkel-api-go/src/config"
)

type MidtransService struct {
	Config *config.Config
	Snap   *snap.Client
}

func NewMidtransService(config *config.Config) *MidtransService {
	var snapClient = snap.Client{}
	var midtransEnv = midtrans.Sandbox
	if config.Env == "production" {
		midtransEnv = midtrans.Production
	}

	snapClient.New(config.PaymentProvider.PaymentProviderAPIKey, midtransEnv)

	return &MidtransService{Config: config, Snap: &snapClient}
}

func (s *MidtransService) CreatePayment(customerDetail PaymentCustomerDetail, orderDetail PaymentOrderDetail) (string, error) {
	midtrans.ServerKey = s.Config.PaymentProvider.PaymentProviderAPIKey
	midtrans.Environment = midtrans.Sandbox

	snapReq := &snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID: orderDetail.OrderID,
			GrossAmt: int64(orderDetail.OrderAmount),
		}, 
		CreditCard: &snap.CreditCardDetails{
			Secure: true,
		},
		CustomerDetail: &midtrans.CustomerDetails{
			FName: customerDetail.CustomerName,
			Email: customerDetail.CustomerEmail,
			Phone: customerDetail.CustomerPhone,
		},
	}

	snapUrl, err := s.Snap.CreateTransactionUrl(snapReq)
	if err != nil {
		return "", err
	}

	return snapUrl, nil
}
