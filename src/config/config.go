package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config ...
type Config struct {
	Env                  string
	Port                 string
	BaseURL              string
	DatabaseURL          string
	LoggerLevel          string
	ContextTimeout       int
	JWTSecretKey         string
	CasbinModelFilePath  string
	CasbinPolicyFilePath string
	Storage              StorageConfig
	Cache                CacheConfig
}

type StorageConfig struct {
	ImageKit    string
	ImageKitURL string
	StorageName string
	AccessKey   string
	SecretKey   string
	Bucket      string
}

type CacheConfig struct {
	CacheServiceName string
	CacheHost        string
	CacheUsername    string
	CachePassword    string
	CachePort        int
}

// LoadConfig will load config from environment variable
func LoadConfig() (config *Config) {
	if err := godotenv.Load(".env"); err != nil {
		panic(err)
	}
	var contextTimeout, cachePort int
	var err error

	env := os.Getenv("ENV")
	port := os.Getenv("APPLICATION_PORT")
	baseURL := os.Getenv("BASE_URL")
	databaseURL := os.Getenv("DATABASE_URL")
	loggerLevel := os.Getenv("LOGGER_LEVEL")
	jwtSecretKey := os.Getenv("JWT_SECRET_KEY")

	contextTimeoutString := os.Getenv("CONTEXT_TIMEOUT")
	if contextTimeoutString != "" {
		contextTimeout, err = strconv.Atoi(os.Getenv("CONTEXT_TIMEOUT"))
		if err != nil {
			panic(err)
		}
	}

	casbinModelFilePath := os.Getenv("CASBIN_MODEL_FILE_PATH")
	casbinPolicyFilePath := os.Getenv("CASBIN_POLICY_FILE_PATH")

	storageUseImageKit := os.Getenv("IMAGEKIT")
	storageImageKitUrl := os.Getenv("IMAGEKIT_BASE_URL")
	storageServiceName := os.Getenv("STORAGE_SERVICE")
	storageAccessKey := os.Getenv("STORAGE_ACCESS_KEY")
	storageSecretKey := os.Getenv("STORAGE_SECRET_KEY")
	storageBucket := os.Getenv("STORAGE_BUCKET")

	cacheServiceName := os.Getenv("CACHE_SERVICE")
	cacheHost := os.Getenv("CACHE_HOST")
	cacheUsername := os.Getenv("CACHE_USERNAME")
	cachePassword := os.Getenv("CACHE_PASSWORD")
	cachePortString := os.Getenv("CACHE_PORT")
	if cachePortString != "" {
		cachePort, err = strconv.Atoi(cachePortString)
		if err != nil {
			panic(err)
		}
	}

	return &Config{
		Env:                  env,
		Port:                 port,
		BaseURL:              baseURL,
		DatabaseURL:          databaseURL,
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
		Cache: CacheConfig{
			CacheServiceName: cacheServiceName,
			CacheHost:        cacheHost,
			CacheUsername:    cacheUsername,
			CachePassword:    cachePassword,
			CachePort:        cachePort,
		},
	}
}
