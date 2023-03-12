package main

import (
	"net/http"
	"os"
	"os/signal"

	"github.com/labstack/echo/v4"
	"github.com/peang/bukabengkel-api-go/src/config"
	"github.com/peang/bukabengkel-api-go/src/handlers"
	"github.com/peang/bukabengkel-api-go/src/middleware"
	repository "github.com/peang/bukabengkel-api-go/src/repositories"
	usecase "github.com/peang/bukabengkel-api-go/src/usecases"
	utils "github.com/peang/bukabengkel-api-go/src/utils"
)

func main() {
	// Load config
	configApp := config.LoadConfig()

	// Setup logger
	appLogger := utils.NewApiLogger(configApp)
	appLogger.InitLogger()

	// Setup Databse
	db, err := config.NewPgxDatabase(configApp.DatabaseURL)
	utils.PanicIfNeeded(err)

	// Setup Casbin Enfocer
	enfocer, err := config.NewCasbinEnfocer(configApp)
	utils.PanicIfNeeded(err)

	// services
	jwtService := config.NewJWTService(configApp.JWTSecretKey, configApp.BaseURL)

	middleware := middleware.NewMiddleware(enfocer, appLogger, jwtService)

	// Repositories
	productRepo := repository.NewProductRepository(db)

	// Usecases
	productUsecase := usecase.NewProductUsecase(productRepo)

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.JWTAuth())

	handlers.NewHealthHandler(e, middleware)
	handlers.NewProductHandler(e, middleware, productUsecase)

	// Start server
	go func() {
		if err := e.Start(":3000"); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal(err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 10 seconds.
	// Use a buffered channel to avoid missing signals as recommended for signal.Notify
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
}
