package s3

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
	"io"
	"strings"
	"time"
)

type S3Service struct {
	client     *s3.Client
	bucketName string
}

func NewS3Service(region, bucketName, accessKey, secretKey string) (*S3Service, error) {
	cfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithRegion(region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKey, secretKey, "")),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	client := s3.NewFromConfig(cfg)
	return &S3Service{
		client:     client,
		bucketName: bucketName,
	}, nil
}

func (s *S3Service) Upload(ctx context.Context, file io.Reader, filename string, contentType string) (string, string, error) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	if seeker, ok := file.(io.Seeker); ok {
		size, err := seeker.Seek(0, io.SeekEnd)
		if err != nil {
			return "", "", fmt.Errorf("failed to check file size: %w", err)
		}
		if size == 0 {
			return "", "", fmt.Errorf("input file is empty")
		}
		_, err = seeker.Seek(0, io.SeekStart)
		if err != nil {
			return "", "", fmt.Errorf("failed to reset file stream: %w", err)
		}
	}

	fName := fmt.Sprintf("%s-%s", uuid.New().String(), filename)
	_, err := s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(s.bucketName),
		Key:         aws.String(fName),
		Body:        file,
		ContentType: aws.String(contentType),
	})
	if err != nil {
		return "", "", fmt.Errorf("failed to upload file: %w", err)
	}

	publicURL := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", s.bucketName, s.client.Options().Region, fName)
	return publicURL, fName, nil
}

func (s *S3Service) Delete(ctx context.Context, fileURL string) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	filename := extractFilename(fileURL, s.bucketName, s.client.Options().Region)
	if filename == "" {
		return fmt.Errorf("invalid file URL: %s", fileURL)
	}

	_, err := s.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(filename),
	})
	if err != nil {
		return fmt.Errorf("failed to delete file: %w", err)
	}
	return nil
}

func (s *S3Service) GetSignedURL(ctx context.Context, filename string) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	presignClient := s3.NewPresignClient(s.client)
	resp, err := presignClient.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(filename),
	}, s3.WithPresignExpires(15*time.Minute))
	if err != nil {
		return "", fmt.Errorf("failed to generate signed URL: %w", err)
	}
	return resp.URL, nil
}

func extractFilename(fileURL, bucketName, region string) string {
	prefix := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/", bucketName, region)
	if strings.HasPrefix(fileURL, prefix) {
		return strings.TrimPrefix(fileURL, prefix)
	}
	return ""
}
