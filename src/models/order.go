package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type Order struct {
	bun.BaseModel `bun:"table:order"`

	ID             int       `bun:"id,pk,nullzero"`
	Key            uuid.UUID `bun:"key,type:uuid,notnull,unique"`
	StoreID        int       `bun:"store_id,notnull"`
	StoreName      string    `bun:"store_name,type:varchar(255),notnull"`
	CustomerID     int       `bun:"customer_id"`
	CustomerName   string    `bun:"customer_name,type:varchar(255)"`
	SalesID        int       `bun:"sales_id"`
	SalesName      string    `bun:"sales_name,type:varchar(255)"`
	InvoiceNumber  string    `bun:"invoice_number,type:varchar(255),notnull"`
	OrderDate      time.Time `bun:"order_date,type:timestamptz,notnull"`
	TotalPrice     float64   `bun:"total_price,type:double precision,notnull"`
	TotalSellPrice float64   `bun:"total_sell_price,type:double precision,notnull"`
	TotalDiscount  float64   `bun:"total_discount,type:double precision,notnull"`
	Total          float64   `bun:"total,type:double precision,notnull"`
	TotalNett      float64   `bun:"total_nett,type:double precision,notnull"`
	Status         int       `bun:"status,notnull"`
	PaymentStatus  int       `bun:"payment_status,notnull"`
	CreatedAt      time.Time `bun:"created_at,type:timestamptz,notnull"`
	UpdatedAt      time.Time `bun:"updated_at,type:timestamptz,notnull"`
	TotalPaid      float64   `bun:"total_paid,type:double precision,notnull"`

	Store *Store `bun:"rel:belongs-to"`
}
