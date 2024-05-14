package utils

import (
	"strings"
	"unicode"
)

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

func ToSnakeCase(s string) string {
	var sb strings.Builder

	// Iterate over each character in the string
	for i, r := range s {
		if i > 0 && unicode.IsUpper(r) {
			// Insert underscore before uppercase letters (except for the first letter)
			sb.WriteRune('_')
		}
		// Convert character to lowercase and append to the result
		sb.WriteRune(unicode.ToLower(r))
	}

	return sb.String()
}
