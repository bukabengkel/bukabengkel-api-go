package entity

import "time"

type ProductCategory struct {
	ID          int
	StoreID     int
	Name        string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
