package distributor_services

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/peang/bukabengkel-api-go/src/config"
	"github.com/peang/bukabengkel-api-go/src/models"
	"github.com/peang/bukabengkel-api-go/src/utils"
)

type CartItem struct {
	ProductCode uint64 `json:"prodcode"`
	Qty         uint64 `json:"qty"`
}

type CheckoutPayload struct {
	Address           string     `json:"address"`
	Notes             string     `json:"notes"`
	ReceiveName       string     `json:"receive_name"`
	ReceivePhone      string     `json:"receive_phone"`
	ReceiveAddress    string     `json:"receive_address"`
	ReceiveAddressZip string     `json:"receive_address_zip"`
	ReceiveProvince   string     `json:"receive_province"`
	ReceiveCity       string     `json:"receive_city"`
	ReceiveDistrict   string     `json:"receive_district"`
	SenderName        string     `json:"sender_name"`
	SenderPhone       string     `json:"sender_phone"`
	SenderAddress     string     `json:"sender_address"`
	SenderAddressZip  string     `json:"sender_address_zip"`
	Dropshipper       string     `json:"dropshipper"`
	ShippingMethod    string     `json:"shipping_method"`
	ShippingPrice     float64    `json:"shipping_price"`
	ShippingETD       string     `json:"shipping_etd"`
	TotalWeight       float64    `json:"total_weight"`
	BookingCode       int        `json:"booking_code"`
	CartItems         []CartItem `json:"cart_items"`
}

type AsianAccessoriesService struct {
	APIKey      string
	CheckoutURL string
}

func NewAsianAccessoriesService(config *config.Config) *AsianAccessoriesService {
	return &AsianAccessoriesService{
		APIKey:      config.AsianAccessoriesAPIKey,
		CheckoutURL: config.AsianAccessoriesCheckoutURL,
	}
}

func (s *AsianAccessoriesService) Checkout(ctx context.Context, orderDistributor *models.OrderDistributor) error {
	items := []CartItem{}
	for _, item := range orderDistributor.Items {
		productKey, err := strconv.ParseUint(item.ProductKey, 10, 64)
		if err != nil {
			return err
		}

		items = append(items, CartItem{
			ProductCode: productKey,
			Qty:         item.Qty,
		})
	}
	payload := CheckoutPayload{
		Address:           orderDistributor.Store.LocationDetail,
		Notes:             orderDistributor.StoreRemarks,
		ReceiveName:       orderDistributor.ShippingAddressName,
		ReceivePhone:      orderDistributor.ShippingAddressPhone,
		ReceiveAddress:    orderDistributor.ShippingAddress,
		ReceiveProvince:   orderDistributor.ShippingAddressProvince,
		ReceiveCity:       orderDistributor.ShippingAddressCity,
		ReceiveDistrict:   orderDistributor.ShippingAddressDistrict,
		ReceiveAddressZip: orderDistributor.ShippingAddressZipCode,
		SenderName:        "Buka Bengkel",
		SenderPhone:       "087770577101",
		SenderAddress:     "Bogor, Indonesia",
		SenderAddressZip:  "16610",
		Dropshipper:       "Yes",
		ShippingMethod:    orderDistributor.ShippingProviderService,
		ShippingPrice:     orderDistributor.TotalShipping,
		ShippingETD:       "1-2",
		TotalWeight:       orderDistributor.ShippingWeight,
		BookingCode:       0,
		CartItems:         items,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	response, err := utils.HttpPostWithRetry(s.CheckoutURL, body, 10)
	if err != nil {
		return err
	}

	fmt.Println(string(response))

	return nil
}
