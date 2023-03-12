package services

type FileServiceInterface interface {
	GetBaseURL() string
	BuildUrl(path string, width int, height int) string
}
