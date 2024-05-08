package utils

import (
	"fmt"
)

func GenerateSort(str string) string {
	len := len(str)
	flag := string(str[0])
	if flag == "-" {
		return fmt.Sprintf("%v %v", ToSnakeCase(str[1:len]), "desc")
	} else {
		return fmt.Sprintf("%v %v", ToSnakeCase(str[0:len]), "asc")
	}
}

func GenerateOffsetLimit(page, perPage int) (offset, limit int) {
	if page < 1 {
		page = 1
	}

	if perPage < 1 {
		perPage = 10
	}

	offset = (page - 1) * perPage
	limit = perPage

	return offset, limit
}
