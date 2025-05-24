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

	ID                     uint64    `bun:"id,pk"`
	Key                    string    `bun:"key,unique"`
	Firstname              string    `bun:"firstname"`
	Lastname               string    `bun:"lastname"`
	Email                  string    `bun:"email"`
	Username               string    `bun:"username"`
	IsEmailVerified        bool      `bun:"is_email_verified"`
	Mobile                 string    `bun:"mobile"`
	IsMobileVerified       bool      `bun:"is_mobile_verified"`
	Status                 uint64    `bun:"status"`
	RefreshToken           string    `bun:"refresh_token"`
	RefreshTokenValidUntil time.Time `bun:"refresh_token_valid_until"`
	ResetToken             string    `bun:"reset_token"`
	ResetTokenValidUntil   time.Time `bun:"reset_token_valid_until"`
}
