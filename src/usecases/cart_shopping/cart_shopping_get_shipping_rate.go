package usecase

import (
	"context"
	"math"

	"github.com/gotidy/ptr"
	"github.com/peang/bukabengkel-api-go/src/models"
	repository "github.com/peang/bukabengkel-api-go/src/repositories"
	"github.com/peang/bukabengkel-api-go/src/services/shipping_services"
	"github.com/peang/bukabengkel-api-go/src/transport/request"
	"github.com/peang/bukabengkel-api-go/src/utils"
)

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

	userStore, err := u.userStoreRepo.FindOne(ctx, repository.UserStoreAggregateRepositoryFilter{
		UserID:  &dto.UserID,
		StoreID: &dto.StoreID,
	})

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
	if distributorLocation == nil {
		return nil, nil, nil, nil, utils.NewUnprocessableEntityError("Distributor Locaton Not Set")
	}

	storeLocation, err := u.locationRepo.FindOne(ctx, repository.LocationRepositoryFilter{
		EntityID:   userStore.Store.ID,
		EntityType: ptr.Of("store"),
	})
	if err != nil {
		return nil, nil, nil, nil, err
	}
	if storeLocation == nil {
		return nil, nil, nil, nil, utils.NewUnprocessableEntityError("Store Locaton Not Set")
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
