package middleware

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/peang/bukabengkel-api-go/src/config"
)

func (m *Middleware) RBAC() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			payload := c.Get("payload").(config.Payload)
			scope := strconv.FormatInt(int64(c.Get("scope").(int)), 10)
			role := strconv.FormatInt(int64(payload.StoreRole), 10)

			obj := c.Path()
			act := c.Request().Method

			if ok, err := m.enforcer.Enforce(scope, role, obj, act); err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, "Failed to enforce scope policy")
			} else if !ok {
				return echo.NewHTTPError(http.StatusForbidden, "Unauthorized")
			}

			return next(c)
		}
	}

}
