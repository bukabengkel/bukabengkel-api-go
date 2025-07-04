package admin_handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/peang/bukabengkel-api-go/src/middleware"
	admin_request "github.com/peang/bukabengkel-api-go/src/transport/request/admin"
	admin_usecase "github.com/peang/bukabengkel-api-go/src/usecases/admin"
	"github.com/peang/bukabengkel-api-go/src/utils"
)

type DashboardHandler struct {
	usecase admin_usecase.DashboardUsecase
}

func NewDashboardHandler(
	e *echo.Echo,
	middleware *middleware.Middleware,
	usecase admin_usecase.DashboardUsecase,
) {
	handler := &DashboardHandler{
		usecase: usecase,
	}

	apiV1 := e.Group("/v1/admin")
	apiV1.GET("/dashboard", handler.Dashboard, middleware.RBAC())
}

func (h *DashboardHandler) Dashboard(ctx echo.Context) (err error) {
	dashboard, err := h.usecase.Dashboard(ctx.Request().Context(), &admin_request.AdminDashboardDTO{})
	if err != nil {
		return utils.ResponseJSON(ctx, http.StatusInternalServerError, "Failed to get dashboard data", nil, err)
	}

	return utils.ResponseJSON(
		ctx,
		http.StatusOK,
		"Admin Dashboard Data",
		dashboard,
		nil,
	)
}