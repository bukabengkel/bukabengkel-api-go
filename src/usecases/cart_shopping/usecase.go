package usecase

import (
	"context"

	"github.com/peang/bukabengkel-api-go/src/models"
	repository "github.com/peang/bukabengkel-api-go/src/repositories"
	"github.com/peang/bukabengkel-api-go/src/services/email_services"
	"github.com/peang/bukabengkel-api-go/src/services/payment_services"
	"github.com/peang/bukabengkel-api-go/src/services/shipping_services"
	"github.com/peang/bukabengkel-api-go/src/transport/request"
)

type CheckoutResponse struct {
	OrderID     string
	OrderAmount float64
	PaymentURL  string
}

type cartShoppingUsecase struct {
	cartRepo             *repository.CartRepository
	distributorRepo      *repository.DistributorRepository
	userStoreRepo        *repository.UserStoreAggregateRepository
	locationRepo         *repository.LocationRepository
	orderDistributorRepo *repository.OrderDistributorRepository
	shippingService      shipping_services.ShippingService
	paymentService       payment_services.PaymentService
	emailService         email_services.EmailService
}

type CartShoppingUsecase interface {
	CartShipping(ctx context.Context, dto *request.CartShoppingGetShippingRateDTO) (*models.Distributor, *models.RajaOngkirLocation, *[]models.CartShopping, *any, error)
	CartCheckout(ctx context.Context, dto *request.CartShoppingCheckoutDTO) (*CheckoutResponse, error)
}

func NewCartShoppingUsecase(
	cartRepo *repository.CartRepository,
	distributorRepo *repository.DistributorRepository,
	userStoreRepo *repository.UserStoreAggregateRepository,
	locationRepo *repository.LocationRepository,
	orderDistributorRepo *repository.OrderDistributorRepository,
	shippingService shipping_services.ShippingService,
	paymentService payment_services.PaymentService,
	emailService email_services.EmailService,
) CartShoppingUsecase {
	return &cartShoppingUsecase{
		cartRepo:             cartRepo,
		distributorRepo:      distributorRepo,
		userStoreRepo:        userStoreRepo,
		locationRepo:         locationRepo,
		orderDistributorRepo: orderDistributorRepo,
		shippingService:      shippingService,
		paymentService:       paymentService,
		emailService:         emailService,
	}
}
