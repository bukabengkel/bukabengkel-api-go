package response

import (
	"time"

	"github.com/peang/bukabengkel-api-go/src/models"
)

type orderDistributorItemResponse struct {
	ProductKey   string  `json:"product_key"`
	ProductName  string  `json:"product_name"`
	ProductUnit  string  `json:"product_unit"`
	ProductImage string  `json:"product_image"`
	Qty          uint64  `json:"qty"`
	Price        float64 `json:"price"`
	Discount     float64 `json:"discount"`
}

type orderDistributorTransactionLogResponse struct {
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
	Remarks   string    `json:"remarks"`
}

type orderDistributorDetailResponse struct {
	ID                      string                                   `json:"id"`
	DistributorName         string                                   `json:"distributor_name"`
	ShippingProvider        string                                   `json:"shipping_provider"`
	ShippingProviderService string                                   `json:"shipping_provider_service"`
	ShippingProviderRemarks string                                   `json:"shipping_provider_remarks"`
	ShippingWeight          float64                                  `json:"shipping_weight"`
	TotalPrice              float64                                  `json:"total_price"`
	TotalDiscount           float64                                  `json:"total_discount"`
	TotalShipping           float64                                  `json:"total_shipping"`
	Total                   float64                                  `json:"total"`
	StoreRemarks            string                                   `json:"store_remarks"`
	ExpiredAt               *time.Time                               `json:"expired_at"`
	PaidAt                  *time.Time                               `json:"paid_at"`
	Items                   []orderDistributorItemResponse           `json:"items"`
	TransactionLogs         []orderDistributorTransactionLogResponse `json:"transaction_logs"`
	TransactionRemarks      string                                   `json:"transaction_remarks"`
	Status                  models.OrderDistributorStatus            `json:"status"`
	CreatedAt               time.Time                                `json:"created_at"`
	UpdatedAt               time.Time                                `json:"updated_at"`
}

type orderDistributorListResponse struct {
	ID              string                        `json:"id"`
	DistributorName string                        `json:"distributor_name"`
	Total           float64                       `json:"total"`
	Status          models.OrderDistributorStatus `json:"status"`
	ExpiredAt       *time.Time                    `json:"expired_at"`
	PaidAt          *time.Time                    `json:"paid_at"`
	CreatedAt       time.Time                     `json:"created_at"`
	UpdatedAt       time.Time                     `json:"updated_at"`
}

func OrderDistributorDetailResponse(orderDistributor *models.OrderDistributor) *orderDistributorDetailResponse {
	items := make([]orderDistributorItemResponse, 0)
	for _, item := range orderDistributor.Items {
		items = append(items, orderDistributorItemResponse{
			ProductKey:   item.ProductKey,
			ProductName:  item.ProductName,
			ProductUnit:  item.ProductUnit,
			ProductImage: item.ProductImage,
			Qty:          item.Qty,
			Price:        item.Price,
			Discount:     item.Discount,
		})
	}

	transactionLogs := make([]orderDistributorTransactionLogResponse, 0)
	for _, log := range orderDistributor.TransactionLogs {
		transactionLogs = append(transactionLogs, orderDistributorTransactionLogResponse{
			Status:    log.Status,
			Timestamp: log.Timestamp,
			Remarks:   log.Remarks,
		})
	}

	response := &orderDistributorDetailResponse{
		ID:                      orderDistributor.Key.String(),
		DistributorName:         orderDistributor.Distributor.Name,
		ShippingProvider:        orderDistributor.ShippingProvider,
		ShippingProviderService: orderDistributor.ShippingProviderService,
		ShippingProviderRemarks: orderDistributor.ShippingProviderRemarks,
		ShippingWeight:          orderDistributor.ShippingWeight,
		TotalPrice:              orderDistributor.TotalPrice,
		TotalDiscount:           orderDistributor.TotalDiscount,
		TotalShipping:           orderDistributor.TotalShipping,
		Total:                   orderDistributor.Total,
		StoreRemarks:            orderDistributor.StoreRemarks,
		ExpiredAt:               orderDistributor.ExpiredAt,
		PaidAt:                  orderDistributor.PaidAt,
		Items:                   items,
		TransactionLogs:         transactionLogs,
		TransactionRemarks:      orderDistributor.TransactionRemarks,
		Status:                  orderDistributor.Status,
		CreatedAt:               orderDistributor.CreatedAt,
		UpdatedAt:               orderDistributor.UpdatedAt,
	}

	return response
}

func OrderDistributorListResponse(orderDistributors *[]models.OrderDistributor) []orderDistributorListResponse {
	var responses = make([]orderDistributorListResponse, 0)
	for _, product := range *orderDistributors {
		response := orderDistributorListResponse{
			ID:              product.Key.String(),
			DistributorName: product.Distributor.Name,
			Total:           product.Total,
			Status:          product.Status,
			ExpiredAt:       product.ExpiredAt,
			PaidAt:          product.PaidAt,
			CreatedAt:       product.CreatedAt,
			UpdatedAt:       product.UpdatedAt,
		}
		responses = append(responses, response)
	}
	return responses
}
