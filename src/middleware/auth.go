package middleware

import (
	"net/http"

	"github.com/labstack/echo/v4"
	utils "github.com/peang/bukabengkel-api-go/src/utils"
)

func (m *Middleware) JWTAuth() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {

			ctx := c.Request().Context()
			tokenStr := c.Request().Header.Get("Authorization")

			if len(tokenStr) != 1 {
				return c.JSON(http.StatusUnauthorized, utils.NewUnauthorizedError("invalid authorization token"))
			}

			tokenInfo, err := m.jwtservice.ValidateToken(ctx, tokenStr)
			// fmt.Println(tokenInfo)
			// if err != nil {
			// 	return c.JSON(
			// 		http.StatusUnauthorized,
			// 		utils.NewUnauthorizedError("invalid authorization token"),
			// 	)
			// }

			// c.Set("user_id", int(tokenInfo.UserID))
			// c.Set("role", tokenInfo.Role)
			// c.Set("scope", tokenInfo.Scope)

			return next(c)
		}
	}
}
