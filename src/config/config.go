package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config ...
type Config struct {
	Port                 string
	BaseURL              string
	DatabaseURL          string
	CacheURL             string
	LoggerLevel          string
	ContextTimeout       int
	JWTSecretKey         string
	CasbinModelFilePath  string
	CasbinPolicyFilePath string
	Storage              StorageConfig
}

type StorageConfig struct {
	ImageKit    string
	ImageKitURL string
	StorageName string
	AccessKey   string
	SecretKey   string
	Bucket      string
}

// LoadConfig will load config from environment variable
func LoadConfig() (config *Config) {
	if err := godotenv.Load(".env"); err != nil {
		panic(err)
	}

	port := os.Getenv("APPLICATION_PORT")
	baseURL := os.Getenv("BASE_URL")
	databaseURL := os.Getenv("DATABASE_URL")
	cacheURL := os.Getenv("CACHE_URL")
	loggerLevel := os.Getenv("LOGGER_LEVEL")
	contextTimeout, _ := strconv.Atoi(os.Getenv("CONTEXT_TIMEOUT"))
	jwtSecretKey := os.Getenv("JWT_SECRET_KEY")

	casbinModelFilePath := os.Getenv("CASBIN_MODEL_FILE_PATH")
	casbinPolicyFilePath := os.Getenv("CASBIN_POLICY_FILE_PATH")

	storageUseImageKit := os.Getenv("IMAGEKIT")
	storageImageKitUrl := os.Getenv("IMAGEKIT_BASE_URL")
	storageServiceName := os.Getenv("STORAGE_SERVICE")
	storageAccessKey := os.Getenv("STORAGE_ACCESS_KEY")
	storageSecretKey := os.Getenv("STORAGE_SECRET_KEY")
	storageBucket := os.Getenv("STORAGE_BUCKET")

	return &Config{
		Port:                 port,
		BaseURL:              baseURL,
		DatabaseURL:          databaseURL,
		CacheURL:             cacheURL,
		LoggerLevel:          loggerLevel,
		ContextTimeout:       contextTimeout,
		JWTSecretKey:         jwtSecretKey,
		CasbinModelFilePath:  casbinModelFilePath,
		CasbinPolicyFilePath: casbinPolicyFilePath,
		Storage: StorageConfig{
			ImageKit:    storageUseImageKit,
			ImageKitURL: storageImageKitUrl,
			StorageName: storageServiceName,
			AccessKey:   storageAccessKey,
			SecretKey:   storageSecretKey,
			Bucket:      storageBucket,
		},
	}
}
