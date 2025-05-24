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
	Page      int `json:"page"`
	PerPage   int `json:"perPage"`
	Total     int `json:"total"`
	TotalPage int `json:"totalPage"`
}

type MetaV2 struct {
	Page    int  `json:"page"`
	PerPage int  `json:"perPage"`
	Next    bool `json:"next"`
}

type ResponseMessage struct {
	Message any `json:"message"`
	Data    any `json:"data"`
	Meta    any `json:"meta,omitempty"`
}

type ResponseErrorMessage struct {
	Message any `json:"message,omitempty"`
	Code    int `json:"code,omitempty"`
	Detail  any `json:"detail,omitempty"`
	Trace   any `json:"trace,omitempty"`
}

type Error struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
	Error   any    `json:"error"`
}

func ResponseJSON(
	ctx echo.Context,
	httpCode int,
	message string,
	data any,
	meta any,
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

func BuildMetaV2(pageString string, perPageString string, next bool) *MetaV2 {
	page, err := strconv.Atoi(pageString)
	if err != nil || page < 1 {
		page = 1
	}

	perPage, err := strconv.Atoi(perPageString)
	if err != nil || perPage < 1 || perPage > 100 {
		perPage = 10
	}

	return &MetaV2{
		Page:    page,
		PerPage: perPage,
		Next:    next,
	}
}
