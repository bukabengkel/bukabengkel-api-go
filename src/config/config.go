package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config ...
type Config struct {
	BaseURL              string
	DatabaseURL          string
	CacheURL             string
	LoggerLevel          string
	ContextTimeout       int
	JWTSecretKey         string
	CasbinModelFilePath  string
	CasbinPolicyFilePath string
}

// LoadConfig will load config from environment variable
func LoadConfig() (config *Config) {
	if err := godotenv.Load(); err != nil {
		panic(err)
	}

	baseURL := os.Getenv("BASE_URL")
	databaseURL := os.Getenv("DATABASE_URL")
	cacheURL := os.Getenv("CACHE_URL")
	loggerLevel := os.Getenv("LOGGER_LEVEL")
	contextTimeout, _ := strconv.Atoi(os.Getenv("CONTEXT_TIMEOUT"))
	jwtSecretKey := os.Getenv("JWT_SECRET_KEY")
	casbinModelFilePath := os.Getenv("CASBIN_MODEL_FILE_PATH")
	casbinPolicyFilePath := os.Getenv("CASBIN_POLICY_FILE_PATH")

	return &Config{
		BaseURL:              baseURL,
		DatabaseURL:          databaseURL,
		CacheURL:             cacheURL,
		LoggerLevel:          loggerLevel,
		ContextTimeout:       contextTimeout,
		JWTSecretKey:         jwtSecretKey,
		CasbinModelFilePath:  casbinModelFilePath,
		CasbinPolicyFilePath: casbinPolicyFilePath,
	}
}
