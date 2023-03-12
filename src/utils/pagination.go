package utils

type Sort struct {
	Field  string
	Method string
}

func GenerateSort(str string) Sort {
	order := Sort{}

	len := len(str)
	flag := string(str[0])
	if flag == "-" {
		order.Method = "DESC"
		order.Field = str[1:len]
	} else {
		order.Method = "ASC"
		order.Field = str[0:len]
	}

	return order
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
