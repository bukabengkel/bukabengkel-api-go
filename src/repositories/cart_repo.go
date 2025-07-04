package repository

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/gotidy/ptr"
	"github.com/peang/bukabengkel-api-go/src/models"
)

type CartRepository struct {
	memoryDB *redis.Client
}

type CartRepositoryFilter struct {
	StoreID       *uint
	UserID        *uint
	DistributorID *uint
}

func NewCartRepository(memoryDB *redis.Client) *CartRepository {
	return &CartRepository{
		memoryDB: memoryDB,
	}
}

func (r *CartRepository) getCartKey(storeID uint, userID uint) string {
	return fmt.Sprintf("cart_store_%d_%d_%d", storeID, userID, models.CartTypeShop)
}

func (r *CartRepository) GetCartShopping(ctx context.Context, storeID uint, userID uint) (*models.Cart, error) {
	cartKey := r.getCartKey(storeID, userID)

	len, err := r.memoryDB.LLen(ctx, cartKey).Result()
	if err != nil {
		return nil, err
	}

	items, err := r.memoryDB.LRange(ctx, cartKey, 0, len).Result()
	if err != nil {
		return nil, err
	}

	cart := &models.Cart{
		ID:                 ptr.Of(uint64(storeID)),
		StoreID:            storeID,
		Items:              make([]models.CartShopping, 0),
		TotalCartSellPrice: 0,
		TotalCartDiscount:  0,
		TotalCart:          0,
	}
	for index, item := range items {
		var cartItem models.CartShopping

		if err := json.Unmarshal([]byte(item), &cartItem); err != nil {
			return nil, err
		}

		cartItem.ID = ptr.Of(uint64(index))
		cart.Items = append(cart.Items, cartItem)
	}

	return cart, nil
}

func (r *CartRepository) EmptyCartShopping(ctx context.Context, storeID uint, userID uint) error {
	cartKey := r.getCartKey(storeID, userID)

	_, err := r.memoryDB.Del(ctx, cartKey).Result()
	if err != nil {
		return err
	}

	return nil
}
