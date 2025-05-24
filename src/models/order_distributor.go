package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type OrderDistributorStatus string

const (
	OrderDistributorStatusCreated           OrderDistributorStatus = "created"
	OrderDistributorStatusWaitingForPayment OrderDistributorStatus = "waiting_for_payment"
	// Wating Response from Payment Provider
	OrderDistributorStatusWaitingForPaymentResponse OrderDistributorStatus = "waiting_payment_response"
	//... other status
	OrderDistributorStatusWaitingForShipment   OrderDistributorStatus = "waiting_for_shippment"
	OrderDistributorStatusInShipping           OrderDistributorStatus = "in_shipping"
	OrderDistributorStatusCustomerReceived     OrderDistributorStatus = "customer_received"
	OrderDistributorStatusCustomerCancelled    OrderDistributorStatus = "customer_cancelled"
	OrderDistributorStatusDistributorCancelled OrderDistributorStatus = "distributor_cancelled"
	OrderDistributorStatusPaymentFailed        OrderDistributorStatus = "payment_failed"
)

type OrderDistributorItem struct {
	ProductKey   string
	ProductName  string
	ProductUnit  string
	ProductImage string
	Qty          uint64
	BasePrice    float64
	BulkPrice    []struct {
		Qty   uint64
		Price float64
	}
	Price    float64
	Discount float64
	Weight   float64
	Volume   float64
}

type OrderDistributorTransactionLog struct {
	Status    string
	Timestamp time.Time
	Remarks   string
}

type OrderDistributor struct {
	bun.BaseModel `bun:"table:order_distributor"`

	ID                      uint                             `bun:"id,pk,autoincrement"`
	Key                     uuid.UUID                        `bun:"key,type:uuid,notnull,unique"`
	DistributorID           uint                             `bun:"distributor_id,notnull"`
	DistributorName         string                           `bun:"distributor_name,notnull"`
	CustomerID              uint                             `bun:"customer_id,notnull"`
	CustomerName            string                           `bun:"customer_name,notnull"`
	ShippingProvider        string                           `bun:"shipping_provider"`
	ShippingProviderService string                           `bun:"shipping_provider_service"`
	ShippingProviderRemarks string                           `bun:"shipping_provider_remarks"`
	ShippingWeight          float64                          `bun:"shipping_weight"`
	TotalPrice              float64                          `bun:"total_price,notnull"`
	TotalDiscount           float64                          `bun:"total_discount,notnull"`
	TotalShipping           float64                          `bun:"total_shipping,notnull"`
	Total                   float64                          `bun:"total,notnull"`
	StoreRemarks            string                           `bun:"store_remarks"`
	ExpiredAt               *time.Time                       `bun:"expired_at"`
	PaidAt                  *time.Time                       `bun:"paid_at"`
	Items                   []OrderDistributorItem           `bun:"items,type:jsonb"`
	TransactionLogs         []OrderDistributorTransactionLog `bun:"transaction_logs,type:jsonb"`
	TransactionRemarks      string                           `bun:"transaction_remarks"`
	Status                  OrderDistributorStatus           `bun:"status,notnull,default:'created'"`
	CreatedAt               time.Time                        `bun:"created_at,notnull,default:current_timestamp"`
	UpdatedAt               time.Time                        `bun:"updated_at,notnull,default:current_timestamp"`

	Distributor *Distributor `bun:"rel:belongs-to,join:distributor_id=id"`
	Store       *Store       `bun:"rel:belongs-to,join:customer_id=id"`
}
