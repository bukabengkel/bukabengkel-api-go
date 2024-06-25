package utils

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

// Meta consist of pagination details
type Meta struct {
	Page      int `json:"page,omitempty"`
	PerPage   int `json:"perPage,omitempty"`
	Total     int `json:"total"`
	TotalPage int `json:"totalPage"`
}

type ResponseMessage struct {
	Message interface{} `json:"message"`
	Data    interface{} `json:"data"`
	Meta    *Meta       `json:"meta,omitempty"`
}

type ResponseErrorMessage struct {
	Message interface{} `json:"message,omitempty"`
	Code    int         `json:"code,omitempty"`
	Detail  interface{} `json:"detail,omitempty"`
	Trace   interface{} `json:"trace,omitempty"`
}

type Error struct {
	Message string      `json:"message"`
	Code    int         `json:"code"`
	Error   interface{} `json:"error"`
}

func ResponseJSON(
	ctx echo.Context,
	httpCode int,
	message string,
	data interface{},
	meta *Meta,
) error {
	w := ctx.Response().Writer

	response := ResponseMessage{
		Message: message,
		Data:    data,
		Meta:    meta,
	}
	resp, _ := json.Marshal(response)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpCode)
	w.Write(resp)
	return nil
}

func ResponseError(
	w http.ResponseWriter,
	err error,
) error {
	code, httpErr := ParseHttpError(err)

	response := ResponseErrorMessage{
		Message: "error",
		Code:    code,
		Detail:  fmt.Sprint(httpErr.Details()),
	}

	resp, _ := json.Marshal(response)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpErr.Status())
	w.Write(resp)
	return nil
}

func BuildMeta(pageString string, perPageString string, count int) *Meta {
	page, err := strconv.Atoi(pageString)
	if err != nil || page < 1 {
		page = 1
	}

	perPage, err := strconv.Atoi(perPageString)
	if err != nil || perPage < 1 || perPage > 100 {
		perPage = 10
	}

	return &Meta{
		Page:      page,
		PerPage:   perPage,
		Total:     count,
		TotalPage: int(math.Ceil(float64(count) / float64(perPage))),
	}
}
