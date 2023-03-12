package middleware

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (m *Middleware) RBAC() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			scope := c.Get("scope").(string)

			obj := c.Path()
			act := c.Request().Method

			if ok, err := m.enforcer.Enforce(scope, obj, act); err != nil {
				fmt.Println(err)
				return echo.NewHTTPError(http.StatusInternalServerError, "Failed to enforce scope policy")
			} else if !ok {
				return echo.NewHTTPError(http.StatusForbidden, "Unauthorized")
			}

			return next(c)
		}
	}

}
