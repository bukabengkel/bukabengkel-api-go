package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/peang/bukabengkel-api-go/src/config"
	"github.com/peang/bukabengkel-api-go/src/middleware"
	"github.com/peang/bukabengkel-api-go/src/services/shipping_services"
	"github.com/peang/bukabengkel-api-go/src/transport/response"
	"github.com/peang/bukabengkel-api-go/src/utils"
)

type LocationHandler struct {
	config *config.Config
	shippingService shipping_services.ShippingService
}

func NewLocationHandler(
	e *echo.Echo,
	middleware *middleware.Middleware,
	c *config.Config,
	shippingService shipping_services.ShippingService,
) {
	handler := &LocationHandler{
		config: c,
		shippingService: shippingService,
	}

	apiV1 := e.Group("/v1/locations")
	apiV1.GET("/search", handler.GetLocations)
}

func (h *LocationHandler) GetLocations(ctx echo.Context) error {
	search := ctx.QueryParam("search")
	searchLabel := ctx.QueryParam("search[label]")
  if searchLabel != "" {
    search = searchLabel
  }

	if search == "" {
		return ctx.JSON(http.StatusBadRequest, "Search or Label is required")
	}

	location, err := h.shippingService.GetLocation(search)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err)
	}

	// Return the parsed data instead of raw body
	return utils.ResponseJSON(
		ctx,
		http.StatusOK,
		"Success",
		location.(response.RajaOngkirGetLocationResponse).Data,
		nil,
	)
}
