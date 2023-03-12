package entity

import "time"

type ProductCategoryEntity struct {
	ID          int
	StoreID     int
	Name        string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
