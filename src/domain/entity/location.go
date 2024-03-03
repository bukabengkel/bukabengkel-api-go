package entity

import "time"

type Location struct {
	ID        int64
	Province  string
	City      string
	District  string
	Urban     string
	CreatedAt time.Time
	UpdatedAt time.Time
}
