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

type ProductHandlerV2 struct {
	usecase usecase.ProductUsecase
}

func NewProductHandlerV2(
	e *echo.Echo,
	middleware *middleware.Middleware,
	usecase usecase.ProductUsecase,
) {
	handler := &ProductHandlerV2{
		usecase: usecase,
	}

	apiV1 := e.Group("/v2/products")
	apiV1.GET("", handler.List, middleware.RBAC())
}

func (h *ProductHandlerV2) List(ctx echo.Context) (err error) {
	storeId, err := strconv.Atoi(ctx.Get("store_id").(string))
	if err != nil {
		return ctx.JSON(utils.ParseHttpError(err))
	}

	dto := request.ProductListDTOV2{
		StoreID:       uint(storeId),
		Limit:         ctx.QueryParam("limit"),
		Sort:          ctx.QueryParam("sort"),
		Keyword:       ctx.QueryParam("keyword"),
		CategoryId:    ctx.QueryParam("categoryId"),
		Name:          ctx.QueryParam("name"),
		StockMoreThan: ctx.QueryParam("stockMoreThan"),
		Status:        ctx.QueryParam("status"),
		Cursor:        ctx.QueryParam("cursor"),
	}

	products, nextCursor, err := h.usecase.ListV2(ctx.Request().Context(), dto)
	if err != nil {
		return ctx.JSON(utils.ParseHttpError(err))
	}

	return utils.ResponseJSONV2(
		ctx,
		http.StatusOK,
		"Product List",
		response.ProductListResponse(products),
		utils.BuildMetaV2(nextCursor),
	)
}
