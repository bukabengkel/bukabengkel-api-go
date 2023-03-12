package entity

import "time"

type LocationEntity struct {
	ID        int
	Province  string
	City      string
	District  string
	Urban     string
	CreatedAt time.Time
	UpdatedAt time.Time
}
