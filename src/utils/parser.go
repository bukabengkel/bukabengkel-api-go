package utils

func Uint(i uint) *uint {
	return &i
}

func IntToInt64(i int) *int64 {
	i64 := int64(i)
	return &i64
}

func String(s string) *string {
	return &s
}
