package handlers

import (
	"errors"
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/peang/bukabengkel-api-go/src/middleware"
	"github.com/peang/bukabengkel-api-go/src/transport/request"
	usecase "github.com/peang/bukabengkel-api-go/src/usecases/webhook.go"
	"github.com/peang/bukabengkel-api-go/src/utils"
)

type WebhookHandler struct {
	usecase usecase.WebhookUsecase
}

func NewWebhookHandler(
	e *echo.Echo,
	middleware *middleware.Middleware,
	usecase usecase.WebhookUsecase,
) {
	handler := &WebhookHandler{
		usecase: usecase,
	}

	apiV1 := e.Group("/v1/webhooks")
	apiV1.POST("/midtrans", handler.MidtransWebhook)
}

func (h *WebhookHandler) MidtransWebhook(ctx echo.Context) error {
	var webhookRequest map[string]interface{}
	if err := ctx.Bind(&webhookRequest); err != nil {
		return ctx.JSON(utils.ParseHttpError(err))
	}

	orderID, ok := webhookRequest["transaction_status"].(string)
	if !ok {
		return ctx.JSON(utils.ParseHttpError(errors.New("order_id is required")))
	}

	fmt.Println(orderID)

	return h.usecase.MidtransWebhook(ctx.Request().Context(), request.MidtransWebhookDTO{})
}
