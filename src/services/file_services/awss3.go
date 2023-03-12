package file_service

import (
	"fmt"

	"github.com/peang/bukabengkel-api-go/src/config"
)

type awss3Service struct {
	UseImageKit string
	ImageKitUrl string
	BaseURL     string
}

func NewAWSS3Service(config *config.Config) *awss3Service {
	return &awss3Service{
		UseImageKit: config.Storage.ImageKit,
		ImageKitUrl: config.Storage.ImageKitURL,
		BaseURL:     "https://bukabengkel.s3.ap-southeast-1.amazonaws.com",
	}
}

func (s *awss3Service) GetBaseURL() string {
	return s.BaseURL
}

func (s *awss3Service) BuildUrl(path string, width int, height int) string {
	if s.UseImageKit == "true" {
		if width != 0 && height != 0 {
			return fmt.Sprintf("%s/%s?tr=w-%d,h-%d", s.ImageKitUrl, path, width, height)
		}
		return fmt.Sprintf("%s/%s", s.ImageKitUrl, path)
	}
	return fmt.Sprintf("%s/%s", s.BaseURL, path)
}
