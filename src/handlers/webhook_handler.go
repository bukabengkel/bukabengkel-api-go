package handlers

import (
	"context"
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/peang/bukabengkel-api-go/src/middleware"
	"github.com/peang/bukabengkel-api-go/src/transport/request"
	usecase "github.com/peang/bukabengkel-api-go/src/usecases/webhook"
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

	orderID, _ := webhookRequest["order_id"].(string)
	if orderID == "" {
		return ctx.JSON(utils.ParseHttpError(errors.New("order_id is required")))
	}

	dto := request.MidtransWebhookDTO{
		OrderID:           orderID,
		TransactionID:     webhookRequest["transaction_id"].(string),
		TransactionStatus: webhookRequest["transaction_status"].(string),
		PaymentType:       webhookRequest["payment_type"].(string),
		GrossAmount:       webhookRequest["gross_amount"].(string),
		FraudStatus:       webhookRequest["fraud_status"].(string),
	}

	if webhookRequest["transaction_status"] == "pending" {
		expiredAt, _ := webhookRequest["expiry_time"].(string)
		dto.ExpiredAt = &expiredAt
	}

	go h.usecase.MidtransWebhook(context.Background(), dto)

	return utils.ResponseJSON(
		ctx,
		http.StatusOK,
		"Webhook Midtrans",
		nil,
		nil,
	)
}
