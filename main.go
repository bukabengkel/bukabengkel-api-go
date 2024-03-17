package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
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

	// Usecases
	productUsecase := usecase.NewProductUsecase(productRepo)

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.JWTAuth())

	handlers.NewProductHandler(e, middleware, productUsecase)

	c := registerCron()
	c.Start()
	// defer c.Stop()

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

func registerCron() *cron.Cron {
	c := cron.New()

	_, err := c.AddFunc("0 0 * * *", func() {
		dir, _ := os.Getwd()

		cmd := exec.Command("go", "run", "cli/main.go", "sync-asian")
		cmd.Dir = dir

		out, err := cmd.CombinedOutput()
		if err != nil {
			log.Printf("Error running command: %v\n", err)
			return
		}
		fmt.Printf("Command output: %s\n", out)

	})

	if err != nil {
		log.Fatal("Error adding cron job:", err)
		return nil
	}

	return c
}
