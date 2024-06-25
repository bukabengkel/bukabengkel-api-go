package handlers

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/peang/bukabengkel-api-go/src/middleware"
	"github.com/peang/bukabengkel-api-go/src/transport/request"
	"github.com/peang/bukabengkel-api-go/src/transport/response"
	usecase "github.com/peang/bukabengkel-api-go/src/usecases"
	"github.com/peang/bukabengkel-api-go/src/utils"
)

type ProductHandler struct {
	usecase usecase.ProductUsecase
}

func NewProductHandler(
	e *echo.Echo,
	middleware *middleware.Middleware,
	usecase usecase.ProductUsecase,
) {
	handler := &ProductHandler{
		usecase: usecase,
	}

	apiV1 := e.Group("/v1/products")
	apiV1.GET("", handler.List, middleware.RBAC())
}

func (h *ProductHandler) List(ctx echo.Context) (err error) {
	storeId, err := strconv.Atoi(ctx.Get("store_id").(string))
	if err != nil {
		return ctx.JSON(utils.ParseHttpError(err))
	}

	dto := request.ProductListDTO{
		StoreID: uint(storeId),
	}
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
		response.ProductListResponse(products),
		utils.BuildMeta(dto.Page, dto.PerPage, count),
	)
}
