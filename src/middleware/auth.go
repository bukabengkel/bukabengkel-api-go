package middleware

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	utils "github.com/peang/bukabengkel-api-go/src/utils"
)

func (m *Middleware) JWTAuth() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ctx := c.Request().Context()
			tokenStr := c.Request().Header.Get("Authorization")

			// if c.Request().URL.Path == "/debug/pprof" || strings.HasPrefix(c.Request().URL.Path, "/debug/pprof/") {
			// 	return next(c)
			// }

			if tokenStr == "" {
				return c.JSON(http.StatusUnauthorized, utils.NewUnauthorizedError("Invalid Auth Token"))
			}

			tokenInfo, err := m.jwtservice.GetTokenInfo(ctx, tokenStr)

			if err != nil {
				return c.JSON(
					http.StatusUnauthorized,
					utils.NewUnauthorizedError(err.Error()),
				)
			}

			c.Set("payload", tokenInfo.Payload)
			c.Set("scope", tokenInfo.Scope)

			if tokenInfo.Payload.StoreID != "" { // This Condition for Bukabengkel Admin
				user := strings.Split(tokenInfo.Payload.StoreID, "-")
				userId := string(user[0])
				c.Set("user_id", userId)

				store := strings.Split(tokenInfo.Payload.StoreID, "-")
				storeId := string(store[0])
				c.Set("store_id", storeId)
			}

			return next(c)
		}
	}
}
