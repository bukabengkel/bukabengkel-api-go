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

type ReportHandler struct {
	usecase usecase.ReportUsecase
}

func NewReportHandler(
	e *echo.Echo,
	middleware *middleware.Middleware,
	usecase usecase.ReportUsecase,
) {
	handler := &ReportHandler{
		usecase: usecase,
	}

	apiV1 := e.Group("/v1/reports")
	apiV1.GET("/orders", handler.OrderSalesReport, middleware.RBAC())
	apiV1.GET("/products", handler.ProductSalesReport, middleware.RBAC())
}

func (h *ReportHandler) OrderSalesReport(ctx echo.Context) (err error) {
	storeId, err := strconv.Atoi(ctx.Get("store_id").(string))
	if err != nil {
		return ctx.JSON(utils.ParseHttpError(err))
	}

	startDate := ctx.QueryParam("start")
	endDate := ctx.QueryParam("end")

	dto := request.OrderSalesReportDTO{
		StoreID:   uint(storeId),
		StartDate: startDate,
		EndDate:   endDate,
	}

	salesReport, err := h.usecase.OrderSalesReport(ctx.Request().Context(), &dto)
	if err != nil {
		return ctx.JSON(utils.ParseHttpError(err))
	}

	return utils.ResponseJSON(
		ctx,
		http.StatusOK,
		"Order Sales Report",
		response.OrderSalesReportResponse(salesReport),
		nil,
	)
}

func (h *ReportHandler) ProductSalesReport(ctx echo.Context) (err error) {
	storeId, err := strconv.Atoi(ctx.Get("store_id").(string))
	if err != nil {
		return ctx.JSON(utils.ParseHttpError(err))
	}

	startDate := ctx.QueryParam("start")
	endDate := ctx.QueryParam("end")
	page := ctx.QueryParam("page")
	perPage := ctx.QueryParam("perPage")

	dto := request.ProductSalesRxeportDTO{
		StoreID:   uint(storeId),
		StartDate: startDate,
		EndDate:   endDate,
		Page:      page,
		PerPage:   perPage,
	}

	productReport, count, err := h.usecase.ProductSalesReport(ctx.Request().Context(), &dto)
	if err != nil {
		return ctx.JSON(utils.ParseHttpError(err))
	}

	return utils.ResponseJSON(
		ctx,
		http.StatusOK,
		"Product Sales Report",
		response.ProductSalesReportResponse(productReport),
		utils.BuildMeta(dto.Page, dto.PerPage, *count),
	)
}
