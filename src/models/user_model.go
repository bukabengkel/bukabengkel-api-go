package models

import (
	"time"

	"github.com/uptrace/bun"
)

type Email struct {
	Email      string
	IsVerified bool
}

type Mobile struct {
	Number     string
	IsVerified bool
}

type Password struct {
	String string
}

type RefreshToken struct {
	Token      string
	ValidUntil time.Time
}

type User struct {
	bun.BaseModel `bun:"table:user"`
}
