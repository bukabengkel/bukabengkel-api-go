package file_service

import (
	"fmt"

	"github.com/peang/bukabengkel-api-go/src/config"
)

type S3Service struct {
	UseImageKit string
	ImageKitUrl string
	BaseURL     string
}

func NewAWSS3Service(config *config.Config) *S3Service {
	return &S3Service{
		UseImageKit: config.Storage.ImageKit,
		ImageKitUrl: config.Storage.ImageKitURL,
		BaseURL:     "https://bukabengkel.s3.ap-southeast-1.amazonaws.com",
	}
}

func (s *S3Service) GetBaseURL() string {
	return s.BaseURL
}

func (s *S3Service) BuildUrl(path string, width int, height int) string {
	if s.UseImageKit == "true" {
		if width != 0 && height != 0 {
			return fmt.Sprintf("%s/%s?tr=w-%d,h-%d", s.ImageKitUrl, path, width, height)
		}
		return fmt.Sprintf("%s/%s", s.ImageKitUrl, path)
	}
	return fmt.Sprintf("%s/%s", s.BaseURL, path)
}
