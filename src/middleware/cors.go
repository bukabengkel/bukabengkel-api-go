package middleware

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func (m *Middleware) CORSMiddleware() echo.MiddlewareFunc {
	return middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{
			"http://localhost:8080",
			"http://localhost:8081",
			"https://admin-dev.bukabengkel.id",
			"https://admin.bukabengkel.id",
			"https://user-dev.bukabengkel.id",
			"https://user-beta.bukabengkel.id",
			"https://user.bukabengkel.id",
		},
		AllowMethods: []string{http.MethodGet},
	})
}
