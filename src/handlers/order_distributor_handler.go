package handlers

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/peang/bukabengkel-api-go/src/middleware"
	"github.com/peang/bukabengkel-api-go/src/transport/request"
	"github.com/peang/bukabengkel-api-go/src/transport/response"
	usecase "github.com/peang/bukabengkel-api-go/src/usecases/order_distributor"
	"github.com/peang/bukabengkel-api-go/src/utils"
)

type OrderDistributorHandler struct {
	usecase usecase.OrderDistributorUsecase
}

func NewOrderDistributorHandler(
	e *echo.Echo,
	middleware *middleware.Middleware,
	usecase usecase.OrderDistributorUsecase,
) {
	handler := &OrderDistributorHandler{
		usecase: usecase,
	}

	apiV1 := e.Group("/v1/order-distributors")
	apiV1.GET("", handler.List, middleware.RBAC())
	apiV1.GET("/:id", handler.Detail, middleware.RBAC())
}

func (h *OrderDistributorHandler) List(ctx echo.Context) (err error) {
	storeId, err := strconv.Atoi(ctx.Get("store_id").(string))
	if err != nil {
		return ctx.JSON(utils.ParseHttpError(err))
	}

	dto := request.OrderDistributorListDTO{
		Page:    ctx.QueryParam("page"),
		PerPage: ctx.QueryParam("perPage"),
		Sort:    ctx.QueryParam("sort"),
		StoreID: uint(storeId),
		OrderID: ctx.QueryParam("order_id"),
	}

	orderDistributors, next, err := h.usecase.List(ctx.Request().Context(), dto)
	if err != nil {
		return ctx.JSON(utils.ParseHttpError(err))
	}

	return utils.ResponseJSON(
		ctx,
		http.StatusOK,
		"Order Distributor List",
		response.OrderDistributorListResponse(orderDistributors),
		utils.BuildMetaV2(dto.Page, dto.PerPage, next),
	)
}

func (h *OrderDistributorHandler) Detail(ctx echo.Context) (err error) {
	storeId, err := strconv.Atoi(ctx.Get("store_id").(string))
	if err != nil {
		return ctx.JSON(utils.ParseHttpError(err))
	}

	dto := request.OrderDistributorDetailDTO{
		StoreID: uint(storeId),
		OrderID: ctx.Param("id"),
	}

	orderDistributor, err := h.usecase.Detail(ctx.Request().Context(), dto)
	if err != nil {
		return ctx.JSON(utils.ParseHttpError(err))
	}

	return utils.ResponseJSON(
		ctx,
		http.StatusOK,
		"Order Distributor Detail",
		response.OrderDistributorDetailResponse(orderDistributor),
		nil,
	)
}
