package file_service

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
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

func newAWSS3Service(config *config.Config) FileServiceInterface {
	sess := session.Must(session.NewSession(&aws.Config{
		Region:      aws.String("ap-southeast-1"),
		Credentials: credentials.NewStaticCredentials(config.Storage.AccessKey, config.Storage.SecretKey, ""),
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
	if path == "" {
		return "https://ik.imagekit.io/bukabengkel/Assets/Bukabengkel%20Placeholder.png?tr=w-500,h-500"
	}

	if s.UseImageKit == "true" {
		if width != 0 && height != 0 {
			return fmt.Sprintf("%s/%s?tr=w-%d,h-%d", s.ImageKitUrl, path, width, height)
		}
		return fmt.Sprintf("%s/%s", s.ImageKitUrl, path)
	}
	return fmt.Sprintf("%s/%s", s.BaseURL, path)
}

func (s *S3Service) Upload(category string, fileUrl string) (*FileUploadResponse, error) {
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

	return &FileUploadResponse{
		Filename:  fileName,
		Size:      fileSize,
		Extension: fileMime,
		Etag:      *img.ETag,
		Bucket:    s.Bucket,
		Key:       key,
	}, nil
}

func (s *S3Service) Delete(filepath string) error {
	_, err := s.Client.DeleteObject(&s3.DeleteObjectInput{
		Bucket: &s.Bucket,
		Key:    &filepath,
	})
	if err != nil {
		return err
	}
	return nil
}
