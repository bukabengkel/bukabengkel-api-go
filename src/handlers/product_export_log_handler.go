package handlers

import (
	"github.com/labstack/echo/v4"
	"github.com/peang/bukabengkel-api-go/src/middleware"
	usecase "github.com/peang/bukabengkel-api-go/src/usecases"
)

type ProductExportLogHandler struct {
	// usecase usecase.ProductExportLogUsecase
}

func NewProductExportLogHandler(
	e *echo.Echo,
	middleware *middleware.Middleware,
	usecase usecase.ProductExportLogUsecase,
) {
	// handler := &ProductExportLogHandler{
	// 	usecase: usecase,
	// }

	// apiV1 := e.Group("/v1/products/export")
	// apiV1.GET("", handler.List, middleware.RBAC())
}
