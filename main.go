package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"

	_ "net/http/pprof"

	"github.com/labstack/echo/v4"
	"github.com/peang/bukabengkel-api-go/src/config"
	"github.com/peang/bukabengkel-api-go/src/handlers"
	"github.com/peang/bukabengkel-api-go/src/middleware"
	repository "github.com/peang/bukabengkel-api-go/src/repositories"
	file_service "github.com/peang/bukabengkel-api-go/src/services/file_services"
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
	db := config.LoadDatabase(configApp)
	defer db.Close()

	// Setup Casbin Enfocer
	enfocer, err := config.NewCasbinEnfocer(configApp)
	utils.PanicIfNeeded(err)

	// services
	jwtService := config.NewJWTService(configApp.JWTSecretKey, configApp.BaseURL)
	fileService := file_service.NewAWSS3Service(configApp)

	middleware := middleware.NewMiddleware(enfocer, appLogger, jwtService)

	// Repositories
	imageRepo := repository.NewImageRepository(db)
	productRepo := repository.NewProductRepository(db, fileService, imageRepo)

	// Usecases
	productUsecase := usecase.NewProductUsecase(productRepo)

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.JWTAuth())

	// handlers.NewPporfHandler(e, middleware)
	// handlers.NewHealthHandler(e, middleware)
	handlers.NewProductHandler(e, middleware, productUsecase)

	defer func() {
		if r := recover(); r != nil {
			log.Println("Application failed to start:", r)
			os.Exit(1)
		}
	}()

	// Start server
	go func() {
		if err := e.Start(fmt.Sprintf(":%s", configApp.Port)); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal(err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 10 seconds.
	// Use a buffered channel to avoid missing signals as recommended for signal.Notify
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
}
