package file_service

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/google/uuid"
	"github.com/peang/bukabengkel-api-go/src/config"
)

type S3Service struct {
	Client      *s3.S3
	Env         string
	Bucket      string
	UseImageKit string
	ImageKitUrl string
	BaseURL     string
}

type S3UploadResponse struct {
	Filename  string
	Size      int64
	Extension string
	Etag      string
	Bucket    string
	Key       string
}

func NewAWSS3Service(config *config.Config) *S3Service {
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String("ap-southeast-1"),
	}))

	return &S3Service{
		Client:      s3.New(sess),
		Env:         config.Env,
		Bucket:      config.Storage.Bucket,
		UseImageKit: config.Storage.ImageKit,
		ImageKitUrl: config.Storage.ImageKitURL,
		BaseURL:     "https://bukabengkel.s3.ap-southeast-1.amazonaws.com",
	}
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

func (s *S3Service) Upload(category string, fileUrl string) (*S3UploadResponse, error) {
	resp, err := http.Get(fileUrl)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	tmp, err := os.CreateTemp("", "tempfile-")
	if err != nil {
		return nil, err
	}
	defer os.Remove(tmp.Name())
	// defer tmp.Close()

	_, err = io.Copy(tmp, resp.Body)
	if err != nil {
		return nil, err
	}

	_, err = tmp.Seek(0, io.SeekStart)
	if err != nil {
		return nil, err
	}

	fileSize := resp.ContentLength
	fileMime := resp.Header.Get("Content-Type")
	fileName := uuid.NewString()

	// Get File Extensions
	ext := filepath.Ext(fileUrl)
	ext = strings.TrimPrefix(ext, ".")

	key := fmt.Sprintf("%v/%v/%v.%v", s.Env, category, fileName, ext)

	img, err := s.Client.PutObject(&s3.PutObjectInput{
		Bucket: &s.Bucket,
		Key:    &key,
		Body:   tmp,
	})
	if err != nil {
		return nil, err
	}

	return &S3UploadResponse{
		Filename:  fileName,
		Size:      fileSize,
		Extension: fileMime,
		Etag:      *img.ETag,
		Bucket:    s.Bucket,
		Key:       key,
	}, nil
}
