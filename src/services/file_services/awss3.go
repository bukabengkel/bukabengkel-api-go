package file_service

import (
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/peang/bukabengkel-api-go/src/config"
)

type S3Service struct {
	Client      *s3.S3
	Session     *session.Session
	Bucket      string
	UseImageKit string
	ImageKitUrl string
	BaseURL     string
}

func NewAWSS3Service(config *config.Config) *S3Service {
	sess, err := session.NewSession(&aws.Config{
		// Region: aws.String(""),
		Credentials: credentials.NewStaticCredentials(
			config.Storage.AccessKey,
			config.Storage.SecretKey,
			"",
		),
	})

	if err != nil {
		log.Fatal("Fail Initialize S3")
	}

	return &S3Service{
		Client:      s3.New(sess),
		Session:     sess,
		Bucket:      config.Storage.Bucket,
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

func (s *S3Service) Upload(filename string) {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}

	defer file.Close()

	uploader := s3manager.NewUploader(s.Session)

	_, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(s.Bucket),
		Key:    aws.String(filename),
		Body:   file,
	})

	return err
}
