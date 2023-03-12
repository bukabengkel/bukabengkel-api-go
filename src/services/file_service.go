package services

import (
	"github.com/peang/bukabengkel-api-go/src/config"
	svc "github.com/peang/bukabengkel-api-go/src/domain/services"
	file_service "github.com/peang/bukabengkel-api-go/src/services/file_services"
	"github.com/peang/bukabengkel-api-go/src/utils"
)

type ImageTransform struct {
	imageWidth  int
	imageHeight int
}

func NewFileService(config *config.Config) (service svc.FileServiceInterface, err error) {
	if config.Storage.StorageName == "awss3" {
		service = file_service.NewAWSS3Service(config)
	} else {
		err = utils.NewInternalServerError("Invalid Storage Name")
	}

	return
}
