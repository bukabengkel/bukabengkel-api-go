package utils

import (
	"fmt"
	"strconv"
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

func GenerateOffsetLimitV2(page, perPage int) (offset, limit int) {
	if page < 1 {
		page = 1
	}

	if perPage < 1 {
		perPage = 10
	}

	offset = (page - 1) * (perPage + 1)
	limit = perPage + 1

	return offset, limit
}

func ParsePageAndPerPage(pageString string, perPageString string) (page int, perPage int, err error) {
	if pageString == "" || pageString == "0" {
		page = 1
	} else {
		page, err = strconv.Atoi(pageString)
		if err != nil {
			return 0, 0, err
		}
	}

	if perPageString == "" || perPageString == "0" {
		perPage = 10
	} else {
		perPage, err = strconv.Atoi(perPageString)
		if err != nil {
			return 0, 0, err
		}
	}
	return
}
