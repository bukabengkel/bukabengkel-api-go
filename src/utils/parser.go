package utils

func Uint(i uint) *uint {
	return &i
}

func Uint64(i int) *uint64 {
	parsed := uint64(i)
	return &parsed
}

func IntToInt64(i int) *int64 {
	parsed := int64(i)
	return &parsed
}

func String(s string) *string {
	return &s
}

func Boolean(b bool) *bool {
	return &b
}
