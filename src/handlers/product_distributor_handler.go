package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/peang/bukabengkel-api-go/src/middleware"
	"github.com/peang/bukabengkel-api-go/src/transport/request"
	"github.com/peang/bukabengkel-api-go/src/transport/response"
	usecase "github.com/peang/bukabengkel-api-go/src/usecases"
	"github.com/peang/bukabengkel-api-go/src/utils"
)

type ProductDistributorHandler struct {
	usecase usecase.ProductDistributorUsecase
}

func NewProductDistributorHandler(
	e *echo.Echo,
	middleware *middleware.Middleware,
	usecase usecase.ProductDistributorUsecase,
) {
	handler := &ProductDistributorHandler{
		usecase: usecase,
	}

	apiV1 := e.Group("/v1/product-distributors")
	apiV1.GET("", handler.List, middleware.RBAC())
	apiV1.GET("/:key", handler.Detail, middleware.RBAC())
}

func (h *ProductDistributorHandler) List(ctx echo.Context) (err error) {
	dto := request.ProductDistributorListDTO{
		Page:          ctx.QueryParam("page"),
		PerPage:       ctx.QueryParam("perPage"),
		Sort:          ctx.QueryParam("sort"),
		Keyword:       ctx.QueryParam("keyword"),
		DistributorID: ctx.QueryParam("distributorId"),
		Name:          ctx.QueryParam("name"),
		Code:          ctx.QueryParam("code"),
		RemoteUpdate:  ctx.QueryParam("remoteUpdate"),
	}

	products, count, err := h.usecase.List(ctx.Request().Context(), dto)
	if err != nil {
		return ctx.JSON(utils.ParseHttpError(err))
	}

	return utils.ResponseJSON(
		ctx,
		http.StatusOK,
		"Product List",
		response.ProductDistributorListResponse(products), utils.BuildMeta(dto.Page, dto.PerPage, count),
	)
}

func (h *ProductDistributorHandler) Detail(ctx echo.Context) (err error) {
	dto := request.ProductDistributorDetailDTO{
		ID: ctx.Param("key"),
	}

	product, err := h.usecase.Detail(ctx.Request().Context(), dto)
	if err != nil {
		return ctx.JSON(utils.ParseHttpError(err))
	}

	if product == nil {
		err = utils.NewHttpError(404, "PRODUCT NOT FOUND", nil)
		return ctx.JSON(utils.ParseHttpError(err))
	}

	return utils.ResponseJSON(
		ctx,
		http.StatusOK,
		"Product Detail",
		response.ProductDistributorDetailResponse(product),
		nil,
	)
}
