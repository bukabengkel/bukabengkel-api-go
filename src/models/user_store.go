package models

import (
	"github.com/uptrace/bun"
)

type UserStore struct {
	bun.BaseModel `bun:"table:user_store"`

	ID                     uint64    `bun:"id,pk"`
	UserID                 uint64    `bun:"user_id"`
	StoreID                uint64    `bun:"store_id"`

	User  *User  `bun:"rel:belongs-to"`
	Store *Store `bun:"rel:belongs-to"`
}
