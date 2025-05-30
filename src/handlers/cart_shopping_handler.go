package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/peang/bukabengkel-api-go/src/middleware"
	"github.com/peang/bukabengkel-api-go/src/transport/request"
	"github.com/peang/bukabengkel-api-go/src/transport/response"
	usecase "github.com/peang/bukabengkel-api-go/src/usecases"
	"github.com/peang/bukabengkel-api-go/src/utils"
)

type CartShoppingHandler struct {
	cartUsecase usecase.CartShoppingUsecase
}

func NewCartShoppingHandler(e *echo.Echo, middleware *middleware.Middleware, cartUsecase usecase.CartShoppingUsecase) {
	handler := &CartShoppingHandler{
		cartUsecase: cartUsecase,
	}

	apiV1 := e.Group("/v1/cart-shopping")
	apiV1.GET("/:distributor_id/shipping-rate", handler.GetShippingRate, middleware.RBAC())
	apiV1.POST("/:distributor_id/checkout", handler.Checkout, middleware.RBAC())
}

func (h *CartShoppingHandler) GetShippingRate(ctx echo.Context) error {
	distributorId := ctx.Param("distributor_id")
	if distributorId == "" {
		return utils.NewBadRequestError("distributor_id is required")
	}

	storeId, err := strconv.Atoi(ctx.Get("store_id").(string))
	if err != nil {
		return utils.NewAuthenticationFailedError("store_id is required")
	}

	userId, err := strconv.Atoi(ctx.Get("user_id").(string))
	if err != nil {
		return utils.NewAuthenticationFailedError("user_id is required")
	}

	dto := request.CartShoppingGetShippingRateDTO{
		StoreID:       uint64(storeId),
		UserID:        uint64(userId),
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

func (h *CartShoppingHandler) Checkout(c echo.Context) error {
	var dto request.CartShoppingCheckoutDTO
	if err := c.Bind(&dto); err != nil {
		return c.JSON(utils.ParseHttpError(err))
	}

	if err := c.Validate(dto); err != nil {
		return utils.NewHTTPValidationError(c, err.(validator.ValidationErrors))
	}

	distributorId := c.Param("distributor_id")
	if distributorId == "" {
		return c.JSON(utils.ParseHttpError(errors.New("distributor_id is required")))
	}

	storeId, err := strconv.Atoi(c.Get("store_id").(string))
	if err != nil {
		return c.JSON(utils.ParseHttpError(errors.New("store_id is required")))
	}

	userId, err := strconv.Atoi(c.Get("user_id").(string))
	if err != nil {
		return c.JSON(utils.ParseHttpError(errors.New("user_id is required")))
	}

	dto.StoreID = uint64(storeId)
	dto.UserID = uint64(userId)
	dto.DistributorID = distributorId

	checkoutResponse, err := h.cartUsecase.CartCheckout(c.Request().Context(), &dto)
	if err != nil {
		return c.JSON(utils.ParseHttpError(err))
	}

	return utils.ResponseJSON(
		c,
		http.StatusOK,
		"Checkout",
		response.CartShoppingCheckoutResponse(
			checkoutResponse.OrderID,
			checkoutResponse.OrderAmount,
			checkoutResponse.PaymentURL,
		),
		nil,
	)
}
