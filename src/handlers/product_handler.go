package handlers

import (
	"net/http"
	"strconv"

	"github.com/gotidy/ptr"
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
	apiV1.POST("/export", handler.Export, middleware.RBAC())
}

func (h *ProductHandler) List(ctx echo.Context) (err error) {
	storeId, err := strconv.Atoi(ctx.Get("store_id").(string))
	if err != nil {
		return ctx.JSON(utils.ParseHttpError(err))
	}

	dto := request.ProductListDTO{
		StoreID:       uint64(storeId),
		Page:          ctx.QueryParam("page"),
		PerPage:       ctx.QueryParam("perPage"),
		Sort:          ctx.QueryParam("sort"),
		Keyword:       ctx.QueryParam("keyword"),
		CategoryId:    ctx.QueryParam("categoryId"),
		Name:          ctx.QueryParam("name"),
		StockMoreThan: ctx.QueryParam("stockMoreThan"),
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

func (h *ProductHandler) Export(ctx echo.Context) (err error) {
	storeId, err := strconv.Atoi(ctx.Get("store_id").(string))
	if err != nil {
		return ctx.JSON(utils.ParseHttpError(err))
	}

	userId, err := strconv.Atoi(ctx.Get("user_id").(string))
	if err != nil {
		return ctx.JSON(utils.ParseHttpError(err))
	}

	var categoryId *uint64
	categoryIdString := ctx.QueryParam("categoryId")
	if categoryIdString != "" {
		categoryIdInt, err := strconv.Atoi(categoryIdString)
		if err != nil {
			return ctx.JSON(utils.ParseHttpError(err))
		}
		categoryId = ptr.Of(uint64(categoryIdInt))
	}

	dto := request.ProductExportDTO{
		StoreID:    uint64(storeId),
		UserID:     uint64(userId),
		CategoryId: categoryId,
	}

	result, err := h.usecase.Export(ctx.Request().Context(), dto)
	if err != nil {
		return ctx.JSON(utils.ParseHttpError(err))
	}

	return utils.ResponseJSON(
		ctx,
		http.StatusOK,
		"Product Export On Progress",
		result,
		nil,
	)
}
