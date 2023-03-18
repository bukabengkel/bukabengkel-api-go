package value_object

import "time"

type RefreshToken struct {
	Token      string
	ValidUntil time.Time
}
