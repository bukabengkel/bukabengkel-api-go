package main

import (
	"fmt"
	"time"

	"github.com/peang/bukabengkel-api-go/src/cmd"
	"github.com/peang/bukabengkel-api-go/src/config"
	repository "github.com/peang/bukabengkel-api-go/src/repositories"
	file_service "github.com/peang/bukabengkel-api-go/src/services/file_services"
	"github.com/peang/bukabengkel-api-go/src/utils"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "my-cli-app",
	Short: "A simple CLI application built with Go and Echo",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Hello from the root command!")
	},
}

func main() {
	// Load config
	configApp := config.LoadConfig()

	// Setup logger
	appLogger := utils.NewApiLogger(configApp)
	appLogger.InitLogger()

	// Setup Databse
	db := config.LoadDatabase(configApp)
	defer db.Close()

	// services
	fileService, err := file_service.NewFileService(configApp)
	utils.PanicIfNeeded(err)

	// Repositories
	imageRepo := repository.NewImageRepository(db, fileService)
	productDistributorRepo := repository.NewProductDistributorRepository(db, fileService)
	productCategoryDistributorRepo := repository.NewProductCategoryDistributorRepository(db)

	Register(appLogger, productDistributorRepo, productCategoryDistributorRepo, imageRepo, fileService)
}

func Register(
	logger utils.Logger,
	productDistributorRepo *repository.ProductDistributorRepository,
	productCategoryDistributorRepo *repository.ProductCategoryDistributorRepository,
	imageRepo *repository.ImageRepository,
	fileService file_service.FileServiceInterface,
) {
	rootCmd := &cobra.Command{
		Use:   "",
		Short: "Ping Command",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Pong !")
		},
	}

	asian := cmd.NewSyncAsian(logger, productDistributorRepo, productCategoryDistributorRepo, imageRepo, fileService)

	syncAsianCmd := &cobra.Command{
		Use:   "sync-asian",
		Short: "Sync Asian Products",
		Run:   asian.ExecuteViaCmd,
	}

	rootCmd.AddCommand(syncAsianCmd)

	startTime := time.Now()
	rootCmd.Execute()

	fmt.Printf("Time Executed : %s", time.Since(startTime))
}
