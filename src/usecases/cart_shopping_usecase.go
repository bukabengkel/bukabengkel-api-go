package usecase

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/google/uuid"
	"github.com/gotidy/ptr"
	"github.com/peang/bukabengkel-api-go/src/models"
	repository "github.com/peang/bukabengkel-api-go/src/repositories"
	"github.com/peang/bukabengkel-api-go/src/services/payment_services"
	"github.com/peang/bukabengkel-api-go/src/services/shipping_services"
	"github.com/peang/bukabengkel-api-go/src/transport/request"
	"github.com/peang/bukabengkel-api-go/src/utils"
)

type CheckoutResponse struct {
	OrderID          string
  OrderAmount      float64
	PaymentURL       string
}

type CartShoppingUsecase interface {
	CartShipping(ctx context.Context, dto *request.CartShoppingGetShippingRateDTO) (*models.Distributor, *models.RajaOngkirLocation, *[]models.CartShopping, *any, error)
	CartCheckout(ctx context.Context, dto *request.CartShoppingCheckoutDTO) (*CheckoutResponse, error)
}

type cartShoppingUsecase struct {
	cartRepo             *repository.CartRepository
	distributorRepo      *repository.DistributorRepository
	userStoreRepo        *repository.UserStoreAggregateRepository
	locationRepo         *repository.LocationRepository
	orderDistributorRepo *repository.OrderDistributorRepository
	shippingService      shipping_services.ShippingService
	paymentService       payment_services.PaymentService
}

func NewCartShoppingUsecase(
	cartRepo *repository.CartRepository,
	distributorRepo *repository.DistributorRepository,
	userStoreRepo *repository.UserStoreAggregateRepository,
	locationRepo *repository.LocationRepository,
	orderDistributorRepo *repository.OrderDistributorRepository,
	shippingService shipping_services.ShippingService,
	paymentService payment_services.PaymentService,
) CartShoppingUsecase {
	return &cartShoppingUsecase{
		cartRepo:             cartRepo,
		distributorRepo:      distributorRepo,
		userStoreRepo:        userStoreRepo,
		locationRepo:         locationRepo,
		orderDistributorRepo: orderDistributorRepo,
		shippingService:      shippingService,
		paymentService:       paymentService,
	}
}

func (u *cartShoppingUsecase) CartShipping(ctx context.Context, dto *request.CartShoppingGetShippingRateDTO) (*models.Distributor, *models.RajaOngkirLocation, *[]models.CartShopping, *any, error) {
	cart, err := u.cartRepo.GetCartShopping(ctx, dto.StoreID, dto.UserID)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	distributor, err := u.distributorRepo.FindOne(repository.DistributorRepositoryFilter{
		Key: &dto.DistributorID,
	})
	if err != nil {
		return nil, nil, nil, nil, err
	}

  fmt.Println(dto.UserID)
  fmt.Println(dto.StoreID)
	userStore, err := u.userStoreRepo.FindOne(ctx, repository.UserStoreAggregateRepositoryFilter{
		UserID:  &dto.UserID,
		StoreID: &dto.StoreID,
	})

  fmt.Println(err)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	distributorLocation, err := u.locationRepo.FindOne(ctx, repository.LocationRepositoryFilter{
		EntityID:   &distributor.ID,
		EntityType: ptr.Of("distributor"),
	})
	if err != nil {
		return nil, nil, nil, nil, err
	}

	storeLocation, err := u.locationRepo.FindOne(ctx, repository.LocationRepositoryFilter{
		EntityID:   userStore.Store.ID,
		EntityType: ptr.Of("store"),
	})
	if err != nil {
		return nil, nil, nil, nil, err
	}

	distributorCartItems := make([]models.CartShopping, 0)
	for _, item := range cart.Items {
		if item.DistributorKey == dto.DistributorID {
			distributorCartItems = append(distributorCartItems, item)
		}
	}
	var totalWeight float64
	for _, item := range distributorCartItems {
		totalWeight += item.Weight
	}

	shippingCost, err := u.shippingService.CalculateShipping(shipping_services.ShippingCostRequest{
		Origin:      distributorLocation.LocationID,
		Destination: storeLocation.LocationID,
		Weight:      int(math.Ceil(totalWeight) * 1000), // in database weight is in KG, but in rajaongkir weight is in gram
	})
	if err != nil {
		return nil, nil, nil, nil, err
	}

	return distributor, distributorLocation, &distributorCartItems, ptr.Of(shippingCost), nil
}

func (u *cartShoppingUsecase) CartCheckout(ctx context.Context, dto *request.CartShoppingCheckoutDTO) (*CheckoutResponse, error) {
	cart, err := u.cartRepo.GetCartShopping(ctx, dto.StoreID, dto.UserID)
	if err != nil {
		return nil, err
	}

	distributor, err := u.distributorRepo.FindOne(repository.DistributorRepositoryFilter{
		Key: &dto.DistributorID,
	})
	if err != nil {
		return nil, err
	}

	userStore, err := u.userStoreRepo.FindOne(ctx, repository.UserStoreAggregateRepositoryFilter{
		UserID:  &dto.UserID,
		StoreID: &dto.StoreID,
	})
	if err != nil {
		return nil, err
	}

	distributorCartItems := make([]models.OrderDistributorItem, 0)
	for _, item := range cart.Items {
		if item.DistributorKey == dto.DistributorID {
			distributorCartItems = append(distributorCartItems, models.OrderDistributorItem{
				ProductKey:   item.ProductKey,
				ProductName:  item.ProductName,
				ProductUnit:  item.ProductUnit,
				ProductImage: item.ProductImage,
				Qty:          item.Qty,
				BasePrice:    item.BasePrice,
				BulkPrice:    item.BulkPrice,
				Price:        item.Price,
				Discount:     item.Discount,
				Weight:       item.Weight,
				Volume:       item.Volume,
			})
		}
	}
	var totalWeight float64
	var totalPrice float64
	var totalDiscount float64
	for _, item := range distributorCartItems {
		totalWeight += item.Weight
		totalPrice += item.Price
		totalDiscount += item.Discount
	}

	distributorLocation, err := u.locationRepo.FindOne(ctx, repository.LocationRepositoryFilter{
		EntityID:   &distributor.ID,
		EntityType: ptr.Of("distributor"),
	})
	if err != nil {
		return nil, err
	}

	storeLocation, err := u.locationRepo.FindOne(ctx, repository.LocationRepositoryFilter{
		EntityID:   userStore.Store.ID,
		EntityType: ptr.Of("store"),
	})
	if err != nil {
		return nil, err
	}

	rajaOngkirShippingCost, err := u.shippingService.CalculateShipping(shipping_services.ShippingCostRequest{
		Origin:      distributorLocation.LocationID,
		Destination: storeLocation.LocationID,
		Weight:      int(math.Ceil(totalWeight) * 1000), // in database weight is in KG, but in rajaongkir weight is in gram
	})

	var shippingCost int
	shippingCosts := rajaOngkirShippingCost.(shipping_services.ShippingResponse)
	for _, cost := range shippingCosts.Data {
		if cost.Service == dto.ShippingProviderService && cost.Code == dto.ShippingProviderCode {
			shippingCost = cost.Cost
			break
		}
	}

	if shippingCost == 0 {
		return nil, utils.NewUnprocessableEntityError("Shipping cost not found")
	}

	orderDistributor := &models.OrderDistributor{
		Key:             uuid.New(),
		DistributorID:   uint(distributor.ID),
		DistributorName: distributor.Name,
		CustomerID:      uint(*userStore.Store.ID),
		CustomerName:    userStore.Store.Name,
		ShippingWeight:  math.Round(totalWeight*100) / 100,
		TotalPrice:      math.Round(totalPrice*100) / 100,
		TotalDiscount:   math.Round(totalDiscount*100) / 100,
		TotalShipping:   math.Round(float64(shippingCost)*100) / 100,
		Total:           math.Round((totalPrice+float64(shippingCost))*100) / 100,
		Items:           distributorCartItems,
		ExpiredAt:       nil,
		PaidAt:          nil,
		TransactionLogs: []models.OrderDistributorTransactionLog{
			{
				Status:    string(models.OrderDistributorStatusCreated),
				Timestamp: time.Now(),
				Remarks:   "Order created",
			},
		},

		// From UI
		ShippingProvider:        dto.ShippingProviderCode,
		ShippingProviderService: dto.ShippingProviderService,
		ShippingProviderRemarks: dto.ShippingProviderRemarks,
		StoreRemarks:            dto.StoreRemarks,

		// Status
		Status:    models.OrderDistributorStatusCreated,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	paymentUrl, err := u.paymentService.CreatePayment(payment_services.PaymentCustomerDetail{
		CustomerName:  userStore.Store.Name,
		CustomerEmail: userStore.User.Email,
		CustomerPhone: userStore.User.Mobile,
	}, payment_services.PaymentOrderDetail{
		OrderID: orderDistributor.Key.String(),
    OrderAmount: orderDistributor.Total,
	})
	if err != nil {
		return nil, err
	}

	orderDistributor.TransactionLogs = append(orderDistributor.TransactionLogs, models.OrderDistributorTransactionLog{
		Status:    string(models.OrderDistributorStatusWaitingForPayment),
		Timestamp: time.Now(),
		Remarks:   "Payment URL: " + paymentUrl,
	})
  orderDistributor.Status = models.OrderDistributorStatusWaitingForPayment
	orderDistributor.TransactionRemarks = paymentUrl

	err = u.orderDistributorRepo.CreateOrderDistributor(ctx, orderDistributor)
	if err != nil {
		return nil, err
	}

  // TODO: Uncomment this when we have a way to empty the cart
	// err = u.cartRepo.EmptyCartShopping(ctx, dto.StoreID, dto.UserID)
	// if err != nil {
	// 	return nil, err
	// }

	return &CheckoutResponse{
		OrderID:          orderDistributor.Key.String(),
		OrderAmount:      orderDistributor.Total,
		PaymentURL:       paymentUrl,
	}, nil
}
