package usecase

import (
	"context"
	"math"
	"time"

	"github.com/google/uuid"
	"github.com/gotidy/ptr"
	"github.com/peang/bukabengkel-api-go/src/models"
	repository "github.com/peang/bukabengkel-api-go/src/repositories"
	"github.com/peang/bukabengkel-api-go/src/services/email_services"
	"github.com/peang/bukabengkel-api-go/src/services/payment_services"
	"github.com/peang/bukabengkel-api-go/src/services/shipping_services"
	"github.com/peang/bukabengkel-api-go/src/transport/request"
	"github.com/peang/bukabengkel-api-go/src/utils"
)

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

	if !userStore.User.IsMobileVerified {
		return nil, utils.NewUnprocessableEntityError("Mobile number must be verified first for checkout")
	}

	if !userStore.User.IsEmailVerified {
		return nil, utils.NewUnprocessableEntityError("Email must be verified first for checkout")
	}

	if userStore.Store.LocationDetail == "" {
		return nil, utils.NewUnprocessableEntityError("Store location must be set first for checkout")
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
	if err != nil {
		return nil, err
	}

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
		Key:            uuid.New(),
		DistributorID:  distributor.ID,
		StoreID:        *userStore.Store.ID,
		UserID:         userStore.User.ID,
		ShippingWeight: math.Round(totalWeight*100) / 100,
		TotalPrice:     math.Round(totalPrice*100) / 100,
		TotalDiscount:  math.Round(totalDiscount*100) / 100,
		TotalShipping:  math.Round(float64(shippingCost)*100) / 100,
		Total:          math.Round((totalPrice+float64(shippingCost))*100) / 100,
		Items:          distributorCartItems,
		ExpiredAt:      nil,
		PaidAt:         nil,
		TransactionLogs: []models.OrderDistributorTransactionLog{
			{
				Status:    string(models.OrderDistributorStatusCreated),
				Timestamp: time.Now(),
				Remarks:   "Order created",
			},
		},

		// From UI
		ShippingProvider:               dto.ShippingProviderCode,
		ShippingProviderService:        dto.ShippingProviderService,
		ShippingProviderTrackingNumber: "",
		ShippingProviderRemarks:        dto.ShippingProviderRemarks,
		ShippingAddress:                userStore.Store.LocationDetail,
		ShippingAddressName:            userStore.User.Firstname + " " + userStore.User.Lastname,
		ShippingAddressPhone:           userStore.User.Mobile,
		ShippingAddressProvince:        storeLocation.ProvinceName,
		ShippingAddressCity:            storeLocation.CityName,
		ShippingAddressDistrict:        storeLocation.DistrictName,
		ShippingAddressSubdistrict:     storeLocation.SubdistrictName,
		ShippingAddressZipCode:         storeLocation.ZipCode,
		StoreRemarks:                   dto.StoreRemarks,

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
		OrderID:     orderDistributor.Key.String(),
		OrderAmount: orderDistributor.Total,
	})
	if err != nil {
		return nil, err
	}

	expiredAt := time.Now().Add(time.Hour * 24)
	orderDistributor.TransactionLogs = append(orderDistributor.TransactionLogs, models.OrderDistributorTransactionLog{
		Status:    string(models.OrderDistributorStatusWaitingForPayment),
		Timestamp: time.Now(),
		Remarks:   "Payment URL: " + paymentUrl,
	})
	orderDistributor.Status = models.OrderDistributorStatusWaitingForPayment
	orderDistributor.ExpiredAt = ptr.Of(expiredAt)
	orderDistributor.TransactionRemarks = paymentUrl

	err = u.orderDistributorRepo.CreateOrderDistributor(ctx, orderDistributor)
	if err != nil {
		return nil, err
	}

	go u.emailService.SendWaitingForPaymentEmail(ctx, userStore.User.Email, email_services.WaitingForPaymentData{
		StoreName:  userStore.Store.Name,
		ExpiredAt:  expiredAt.Format("02 Jan 2006 15:04"),
		PaymentURL: paymentUrl,
	})

	// TODO: Uncomment this
	// err = u.cartRepo.EmptyCartShopping(ctx, dto.StoreID, dto.UserID)
	// if err != nil {
	// 	return nil, err
	// }

	return &CheckoutResponse{
		OrderID:     orderDistributor.Key.String(),
		OrderAmount: orderDistributor.Total,
		PaymentURL:  paymentUrl,
	}, nil
}
