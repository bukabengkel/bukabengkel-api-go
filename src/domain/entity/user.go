package entity

import (
	"time"

	vo "github.com/peang/bukabengkel-api-go/src/domain/value_objects"
)

type UserStatus int

const (
	UserInactive UserStatus = iota
	UserActive
)

func (s UserStatus) String() string {
	switch s {
	case UserInactive:
		return "inactive"
	case UserActive:
		return "active"
	default:
		return "unknown"
	}
}

type User struct {
	ID        int
	Key       string
	FirstName string
	LastName  string
	Email     vo.Email
	Mobile    vo.Mobile
	Status    struct {
		ID   int
		Name string
	}
	RefreshToken vo.RefreshToken
	Password     vo.Password
	ProfileImage Image
	CreatedAt    time.Time
	CreatedBy    string
	UpdatedAt    time.Time
	UpdatedBy    string
	DeletedAt    *time.Time
	DeletedBy    string
}
