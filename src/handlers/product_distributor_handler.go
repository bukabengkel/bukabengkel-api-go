package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/peang/bukabengkel-api-go/src/middleware"
	"github.com/peang/bukabengkel-api-go/src/transport/request"
	"github.com/peang/bukabengkel-api-go/src/transport/response"
	usecase "github.com/peang/bukabengkel-api-go/src/usecases"
	"github.com/peang/bukabengkel-api-go/src/utils"
)

type ProductDistributorHandler struct {
	usecase usecase.ProductDistributorUsecase
}

func NewProductDistributorHandler(
	e *echo.Echo,
	middleware *middleware.Middleware,
	usecase usecase.ProductDistributorUsecase,
) {
	handler := &ProductDistributorHandler{
		usecase: usecase,
	}

	apiV1 := e.Group("/v1/product-distributors")
	apiV1.GET("", handler.List, middleware.RBAC())
}

func (h *ProductDistributorHandler) List(ctx echo.Context) (err error) {
	dto := request.ProductDistributorListDTO{}
	if err := ctx.Bind(&dto); err != nil {
		return err
	}

	products, count, err := h.usecase.List(ctx.Request().Context(), dto)
	if err != nil {
		return ctx.JSON(utils.ParseHttpError(err))
	}

	return utils.ResponseJSON(
		ctx,
		http.StatusOK,
		"Product List",
		response.ProductDistributorListResponse(products), utils.BuildMeta(dto.Page, dto.PerPage, count),
	)
}
