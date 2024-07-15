package file_services

import (
	"fmt"

	"github.com/peang/bukabengkel-api-go/src/config"
)

type FileService interface {
	BuildUrl(path string, width int, height int) string
	Upload(category string, fileUrl string) (*FileUploadResponse, error)
	Delete(filepath string) error
}

type FileUploadResponse struct {
	Filename  string
	Size      int64
	Extension string
	Etag      string
	Bucket    string
	Key       string
}

const (
	FILE_SERVICE_AWSS3 = "awss3"
)

func NewFileService(config *config.Config) (FileService, error) {
	switch config.Storage.StorageName {
	case FILE_SERVICE_AWSS3:
		return newAWSS3Service(config), nil
	default:
		return nil, fmt.Errorf("storage %s not supported", config.Storage.StorageName)
	}
}
