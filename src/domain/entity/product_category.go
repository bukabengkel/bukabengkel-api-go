package entity

import "time"

type ProductCategory struct {
	ID          int64
	StoreID     int64
	Name        string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
