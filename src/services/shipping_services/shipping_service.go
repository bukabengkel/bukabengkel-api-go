package shipping_services

import (
	"fmt"

	"github.com/peang/bukabengkel-api-go/src/config"
)

type ShippingService interface {
	GetLocation(search string) (interface{}, error)
	CalculateShipping(data interface{}) (interface{}, error)
}

type FileUploadResponse struct {
	Filename  string
	Size      int64
	Extension string
	Etag      string
	Bucket    string
	Key       string
}

const (
	SHIPPING_SERVICE_RAJAONGKIR = "rajaongkir"
)

func NewShippingService(config *config.Config) (ShippingService, error) {
	switch config.ShippingProvider.ShippingProviderName {
	case SHIPPING_SERVICE_RAJAONGKIR:
		return NewRajaOngkirService(config), nil
	default:
		return nil, fmt.Errorf("shipping provider %s not supported", config.ShippingProvider.ShippingProviderName)
	}
}
