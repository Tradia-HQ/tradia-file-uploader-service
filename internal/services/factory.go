package services

import (
	"fmt"
	"github.com/Tradia-HQ/tradia-file-uploader-service/internal/config"
	"github.com/Tradia-HQ/tradia-file-uploader-service/internal/services/digitalOcean"
	"github.com/Tradia-HQ/tradia-file-uploader-service/internal/services/gcs"
	"github.com/Tradia-HQ/tradia-file-uploader-service/internal/services/interfaces"
	"github.com/Tradia-HQ/tradia-file-uploader-service/internal/services/s3"
)

func NewUploader(cfg *config.Config) (interfaces.FileUploader, error) {
	switch cfg.Provider {
	case "gcs":
		return gcs.NewGCSService(cfg.GCS.ProjectID, cfg.GCS.BucketName, cfg.GCS.CredentialsPath)
	case "s3":
		return s3.NewS3Service(cfg.S3.Region, cfg.S3.BucketName, cfg.S3.AccessKey, cfg.S3.SecretKey)
	case "do":
		return digitalOcean.NewSpacesService(cfg.DSpace.Region, cfg.DSpace.BucketName, cfg.DSpace.AccessKey, cfg.DSpace.SecretKey, cfg.DSpace.Endpoint)
	default:
		return nil, fmt.Errorf("unsupported storage provider: %s", cfg.Provider)
	}
}
