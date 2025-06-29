package config

import (
	"github.com/joho/godotenv"
	"os"
	"strconv"
)

type Config struct {
	Port     int
	Env      string
	Provider string // "gcs", "s3", or "do" (DigitalOcean)
	GCS      struct {
		ProjectID       string
		BucketName      string
		CredentialsPath string
	}
	S3 struct {
		Region     string
		BucketName string
		AccessKey  string
		SecretKey  string
	}
	DSpace struct { // DigitalOcean Spaces
		Region     string
		BucketName string
		AccessKey  string
		SecretKey  string
		Endpoint   string
	}
}

func LoadConfig() (*Config, error) {
	_ = godotenv.Load() // Load .env file

	cfg := &Config{
		Provider: os.Getenv("STORAGE_PROVIDER"),
	}
	cfg.GCS.ProjectID = os.Getenv("GCS_PROJECT_ID")
	cfg.GCS.BucketName = os.Getenv("GCS_BUCKET_NAME")
	cfg.GCS.CredentialsPath = os.Getenv("GCS_CREDENTIALS_PATH")
	cfg.S3.Region = os.Getenv("AWS_REGION")
	cfg.S3.BucketName = os.Getenv("AWS_BUCKET_NAME")
	cfg.S3.AccessKey = os.Getenv("AWS_ACCESS_KEY_ID")
	cfg.S3.SecretKey = os.Getenv("AWS_SECRET_ACCESS_KEY")
	cfg.DSpace.Region = os.Getenv("DO_REGION")
	cfg.DSpace.BucketName = os.Getenv("DO_BUCKET_NAME")
	cfg.DSpace.AccessKey = os.Getenv("DO_ACCESS_KEY_ID")
	cfg.DSpace.SecretKey = os.Getenv("DO_SECRET_ACCESS_KEY")
	cfg.DSpace.Endpoint = os.Getenv("DO_ENDPOINT")

	port, _ := strconv.Atoi(os.Getenv("PORT"))
	cfg.Port = port
	cfg.Env = os.Getenv("ENV")

	return cfg, nil
}
