package main

import (
	"fmt"
	"log"
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
	productExportLogRepo := repository.NewProductExportLogRepository(db)
	orderRepo := repository.NewOrderRepository(db)

	// Usecases
	reportUsecase := usecase.NewReportUsecase(orderRepo)
	productUsecase := usecase.NewProductUsecase(productRepo)
	productDistributorUsecase := usecase.NewProductDistributorUsecase(productDistRepo)
	productExportLogUsecase := usecase.NewProductExportLogUsecase(productExportLogRepo)

	e := echo.New()
	e.Use(middleware.CORSMiddleware())

	handlers.NewReportHandler(e, middleware, reportUsecase)
	handlers.NewProductHandler(e, middleware, productUsecase)
	handlers.NewProductDistributorHandler(e, middleware, productDistributorUsecase)
	handlers.NewProductExportLogHandler(e, middleware, productExportLogUsecase)

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

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", configApp.Port)))

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
}
