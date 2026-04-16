package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	App AppConfig
	DB  DBConfig
	JWT JWTConfig
	S3  S3Config
	SQS SQSConfig
	AWS AWSConfig
}

type AppConfig struct {
	Env  string
	Port string
}

type DBConfig struct {
	Path string
}

type JWTConfig struct {
	Secret string
	Expiry int
}

func Load() (*Config, error) {
	_ = godotenv.Load(".env")

	return &Config{
		App: AppConfig{
			Env:  getEnv("APP_ENV", "development"),
			Port: getEnv("APP_PORT", "8080"),
		},
		DB: DBConfig{
			Path: getEnv("DB_PATH", "./db/main.db"),
		},
		JWT: JWTConfig{
			Secret: getEnv("JWT_SECRET", "secret"),
			Expiry: getEnvInt("JWT_EXPIRY", 24),
		},
		S3: S3Config{
			Endpoint:   getEnv("S3_ENDPOINT", ""),
			BucketName: getEnv("S3_BUCKET_NAME", ""),
			Region:     getEnv("S3_REGION", "us-east-1"),
		},
		SQS: SQSConfig{
			Endpoint: getEnv("SQS_ENDPOINT", ""),
			QueueURL: getEnv("SQS_QUEUE_URL", ""),
		},
		AWS: AWSConfig{
			Region: getEnv("AWS_REGION", "us-east-1"),
		},
	}, nil
}

func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}

	return defaultVal
}

func getEnvInt(key string, defaultVal int) int {
	if val := os.Getenv(key); val != "" {
		if i, err := strconv.Atoi(val); err == nil {
			return i
		}
	}
	return defaultVal
}

type S3Config struct {
	Endpoint   string
	BucketName string
	Region     string
}

type SQSConfig struct {
	Endpoint string
	QueueURL string
}

type AWSConfig struct {
	Region string
}
