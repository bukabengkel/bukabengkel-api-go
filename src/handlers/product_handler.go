package handlers

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/peang/bukabengkel-api-go/src/domain/entity"
	"github.com/peang/bukabengkel-api-go/src/middleware"
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

	return
}

func (h *ProductHandler) List(ctx echo.Context) (err error) {
	// ctx := ctx.Request().Context()

	filter := entity.ProductEntityRepositoryFilter{}
	sort := "id"

	page, err := strconv.Atoi(ctx.QueryParam("page"))
	if err != nil || page < 1 {
		page = 1
	}

	perPage, err := strconv.Atoi(ctx.QueryParam("perPage"))
	if err != nil || perPage < 1 || perPage > 100 {
		perPage = 10
	}

	if ctx.QueryParam("name") != "" && len(ctx.QueryParam("name")) >= 3 {
		filter.Name = ctx.QueryParam("name")
	}

	// if ctx.QueryParam("status") != "" {
	// 	filter.Status = ctx.QueryParam("status")
	// }

	if ctx.QueryParam("sort") != "" {
		sort = ctx.QueryParam("orders")
	}

	products, count, err := h.usecase.List(ctx.Request().Context(), page, perPage, sort, filter)
	if err != nil {
		return ctx.JSON(utils.ParseHttpError(err))
	}

	meta := utils.BuildMeta(page, perPage, count)
	return utils.ResponseJSON(ctx, http.StatusOK, "Product List", response.ProductListResponse(products), meta)
}
