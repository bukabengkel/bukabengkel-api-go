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

type DistributorHandler struct {
	usecase usecase.DistributorUsecase
}

func NewDistributorHandler(
	e *echo.Echo,
	middleware *middleware.Middleware,
	usecase usecase.DistributorUsecase,
) {
	handler := &DistributorHandler{
		usecase: usecase,
	}

	apiV1 := e.Group("/v1/distributors")
	apiV1.GET("", handler.List, middleware.RBAC())
}

func (h *DistributorHandler) List(ctx echo.Context) (err error) {
	dto := request.DistributorListDTO{
		Page:    ctx.QueryParam("page"),
		PerPage: ctx.QueryParam("perPage"),
		Sort:    ctx.QueryParam("sort"),
		Name:    ctx.QueryParam("name"),
	}

	distributors, count, err := h.usecase.List(ctx.Request().Context(), dto)
	if err != nil {
		return ctx.JSON(utils.ParseHttpError(err))
	}

	return utils.ResponseJSON(
		ctx,
		http.StatusOK,
		"Distributor List",
		response.DistributorListResponse(distributors),
		utils.BuildMeta(dto.Page, dto.PerPage, count),
	)
}
