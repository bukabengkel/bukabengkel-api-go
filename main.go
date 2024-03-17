package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"

	_ "net/http/pprof"

	"github.com/labstack/echo/v4"
	"github.com/peang/bukabengkel-api-go/src/cmd"
	"github.com/peang/bukabengkel-api-go/src/config"
	"github.com/peang/bukabengkel-api-go/src/handlers"
	"github.com/peang/bukabengkel-api-go/src/middleware"
	repository "github.com/peang/bukabengkel-api-go/src/repositories"
	file_service "github.com/peang/bukabengkel-api-go/src/services/file_services"
	usecase "github.com/peang/bukabengkel-api-go/src/usecases"
	utils "github.com/peang/bukabengkel-api-go/src/utils"
	"github.com/robfig/cron/v3"
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
	s3service := file_service.NewAWSS3Service(configApp)

	middleware := middleware.NewMiddleware(enfocer, appLogger, jwtService)

	// Repositories
	imageRepo := repository.NewImageRepository(db, s3service)
	productRepo := repository.NewProductRepository(db, imageRepo)
	productDistRepo := repository.NewProductDistributorRepository(db, s3service)
	productCatDistRepo := repository.NewProductCategoryDistributorRepository(db)

	// Usecases
	productUsecase := usecase.NewProductUsecase(productRepo)
	productDistributorUsecase := usecase.NewProductDistributorUsecase(productDistRepo)

	e := echo.New()
	// e.Use(mw.CORS())
	e.Use(middleware.CORSMiddleware())

	handlers.NewProductHandler(e, middleware, productUsecase)
	handlers.NewProductDistributorHandler(e, middleware, productDistributorUsecase)

	e.Use(middleware.Logger())
	e.Use(middleware.JWTAuth())

	c := cron.New()
	_, err = c.AddFunc("0 0 * * *", func() {
		fmt.Println("Executing Sync Asian Products")
		asian := cmd.NewSyncAsian(appLogger, productDistRepo, productCatDistRepo, imageRepo, s3service)

		asian.Execute()
	})
	if err != nil {
		log.Fatal("Fail to Register Cron")
	}
	c.Start()

	defer func() {
		if r := recover(); r != nil {
			log.Println("Application failed to start:", r)
			os.Exit(1)
		}
	}()

	go func() {
		if err := e.Start(fmt.Sprintf(":%s", configApp.Port)); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
}
