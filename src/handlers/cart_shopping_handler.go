package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/peang/bukabengkel-api-go/src/middleware"
	"github.com/peang/bukabengkel-api-go/src/transport/request"
	"github.com/peang/bukabengkel-api-go/src/transport/response"
	usecase "github.com/peang/bukabengkel-api-go/src/usecases"
	"github.com/peang/bukabengkel-api-go/src/utils"
)

type CartHandler struct {
	cartUsecase usecase.CartShoppingUsecase
}

func NewCartShoppingHandler(e *echo.Echo, middleware *middleware.Middleware, cartUsecase usecase.CartShoppingUsecase) {
	handler := &CartHandler{
		cartUsecase: cartUsecase,
	}

	apiV1 := e.Group("/v1/cart-shopping")
	apiV1.GET("/:distributor_id/shipping-rate", handler.GetShippingRate)
	apiV1.POST("/:distributor_id/checkout", handler.Checkout)
}

func (h *CartHandler) GetShippingRate(ctx echo.Context) error {
	distributorId := ctx.Param("distributor_id")
	if distributorId == "" {
		return ctx.JSON(utils.ParseHttpError(errors.New("distributor_id is required")))
	}

	storeId, err := strconv.Atoi(ctx.Get("store_id").(string))
	if err != nil {
		return ctx.JSON(utils.ParseHttpError(err))
	}

	userId, err := strconv.Atoi(ctx.Get("user_id").(string))
	if err != nil {
		return ctx.JSON(utils.ParseHttpError(err))
	}

	dto := request.CartGetShippingRateDTO{
		StoreID: uint64(storeId),
		UserID:  uint64(userId),
		DistributorID: distributorId,
	}

	distributor, distributorLocation, cartItems, shippingCost, err := h.cartUsecase.CartShipping(ctx.Request().Context(), &dto)
	if err != nil {
		return ctx.JSON(utils.ParseHttpError(err))
	}

	return utils.ResponseJSON(
		ctx,
		http.StatusOK,
		"Shipping Rate List",
		response.CartShippingRateResponse(
			distributor,
			distributorLocation,
			cartItems,
			shippingCost,
		),
		nil,
	)	
}

func (h *CartHandler) Checkout(ctx echo.Context) error {
	distributorId := ctx.Param("distributor_id")
	if distributorId == "" {
		return ctx.JSON(utils.ParseHttpError(errors.New("distributor_id is required")))
	}

	storeId, err := strconv.Atoi(ctx.Get("store_id").(string))
	if err != nil {
		return ctx.JSON(utils.ParseHttpError(err))
	}

	userId, err := strconv.Atoi(ctx.Get("user_id").(string))
	if err != nil {
		return ctx.JSON(utils.ParseHttpError(err))
	}

	dto := request.CartCheckoutDTO{
		StoreID: uint64(storeId),
		UserID:  uint64(userId),
		DistributorID: distributorId,
	}

	fmt.Println(dto)
	
	return utils.ResponseJSON(
		ctx,
		http.StatusOK,
		"Checkout",
		nil,
		nil,
	)	
}
