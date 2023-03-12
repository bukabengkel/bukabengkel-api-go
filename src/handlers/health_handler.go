package handlers

import (
	"github.com/labstack/echo/v4"
	"github.com/peang/bukabengkel-api-go/src/middleware"
	"github.com/peang/bukabengkel-api-go/src/utils"
)

type HealthCheckHandler struct {
}

func NewHealthHandler(e *echo.Echo, middleware *middleware.Middleware) {
	handler := &HealthCheckHandler{}

	e.GET("/health-check", handler.Check)

	return
}

func (h *HealthCheckHandler) Check(ctx echo.Context) (err error) {
	return utils.ResponseJSON(ctx, 200, "OK", nil, nil)
}
