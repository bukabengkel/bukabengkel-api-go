package shipping_services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/peang/bukabengkel-api-go/src/config"
	"github.com/peang/bukabengkel-api-go/src/transport/response"
)

type RajaOngkirService struct {
	APIKey string
}

type ShippingCostRequest struct {
	Origin      uint64 `json:"origin"`
	Destination uint64 `json:"destination"`
	Weight      int    `json:"weight"`
}

func NewRajaOngkirService(config *config.Config) *RajaOngkirService {
	return &RajaOngkirService{
		APIKey: config.ShippingProvider.ShippingProviderAPIKey,
	}
}

func (s *RajaOngkirService) GetLocation(search string) (any, error) {
	client := http.Client{
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequest("GET", fmt.Sprintf("https://rajaongkir.komerce.id/api/v1/destination/domestic-destination?search=%s", search), nil)
	req.Header.Set("key", s.APIKey)

	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Parse the JSON response into our struct
	var rajaOngkirResponse response.RajaOngkirGetLocationResponse
	if err := json.Unmarshal(body, &rajaOngkirResponse); err != nil {
		return nil, err
	}

	return rajaOngkirResponse, nil
}

func (s *RajaOngkirService) CalculateShipping(data any) (any, error) {
	requestData, ok := data.(ShippingCostRequest)
	if !ok {
		return nil, fmt.Errorf("invalid request data format")
	}

	client := http.Client{
		Timeout: 10 * time.Second,
	}

	// Create form data
	formData := url.Values{}
	formData.Set("origin", strconv.FormatUint(requestData.Origin, 10))
	formData.Set("destination", strconv.FormatUint(requestData.Destination, 10))
	formData.Set("weight", strconv.Itoa(requestData.Weight))
	formData.Set("courier", "jne:sicepat")

	req, err := http.NewRequest("POST", "https://rajaongkir.komerce.id/api/v1/calculate/domestic-cost", strings.NewReader(formData.Encode()))
	if err != nil {
		return nil, err
	}

	req.Header.Set("key", s.APIKey)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var shippingCostResponse response.RajaOngkirGetShippingRateResponse
	if err := json.Unmarshal(body, &shippingCostResponse); err != nil {
		return nil, err
	}

	return shippingCostResponse.Data, nil
}
