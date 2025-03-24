package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/peang/bukabengkel-api-go/src/config"
	"github.com/peang/bukabengkel-api-go/src/middleware"
	"github.com/peang/bukabengkel-api-go/src/transport/response"
	"github.com/peang/bukabengkel-api-go/src/utils"
)

type RajaOngkirHandler struct {
	config *config.Config
}

func NewRajaOngkirHandler(
	e *echo.Echo,
	middleware *middleware.Middleware,
	c *config.Config,
) {
	handler := &RajaOngkirHandler{
		config: c,
	}

	apiV1 := e.Group("/v1/rajaongkir")
	apiV1.GET("/locations", handler.GetLocations)
}

func (h *RajaOngkirHandler) GetLocations(ctx echo.Context) error {
	search := ctx.QueryParam("search")
	if search == "" {
		return ctx.JSON(http.StatusBadRequest, "Search is required")
	}

	client := http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprintf("https://rajaongkir.komerce.id/api/v1/destination/domestic-destination?search=%s", search), nil)
	req.Header.Set("key", h.config.RajaOngkirAPIKey)
	
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err)
	}

	// Parse the JSON response into our struct
	var rajaOngkirResponse response.RajaOngkirResponse
	if err := json.Unmarshal(body, &rajaOngkirResponse); err != nil {
		return ctx.JSON(http.StatusInternalServerError, err)
	}

	// Return the parsed data instead of raw body
	return utils.ResponseJSON(
		ctx,
		http.StatusOK,
		"Success",
		rajaOngkirResponse.Data,
		nil,
	)
}
