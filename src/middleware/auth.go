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
				storeId, storeKey := strings.Split(tokenInfo.Payload.StoreID, "-")[0], strings.Join(strings.Split(tokenInfo.Payload.StoreID, "-")[1:], "-")
				c.Set("store_id", storeId)
				c.Set("store_key", storeKey)
			}

			if tokenInfo.Payload.ID != "" {
				userId, userKey := strings.Split(tokenInfo.Payload.ID, "-")[0], strings.Join(strings.Split(tokenInfo.Payload.ID, "-")[1:], "-")
				c.Set("user_id", userId)
				c.Set("user_key", userKey)
			}

			return next(c)
		}
	}
}
