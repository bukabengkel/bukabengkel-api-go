package utils

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"

	"github.com/labstack/echo/v4"
)

// Meta consist of pagination details
type Meta struct {
	Page      int `json:"page,omitempty"`
	PerPage   int `json:"perPage,omitempty"`
	Total     int `json:"total,omitempty"`
	TotalPage int `json:"totalPage,omitempty"`
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

func BuildMeta(page int, perPage int, count int) *Meta {
	return &Meta{
		Page:      page,
		PerPage:   perPage,
		Total:     count,
		TotalPage: int(math.Ceil(float64(count) / float64(perPage))),
	}
}
