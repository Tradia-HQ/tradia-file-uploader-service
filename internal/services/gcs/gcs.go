package gcs

import (
	"cloud.google.com/go/storage"
	"context"
	"fmt"
	"github.com/google/uuid"
	"google.golang.org/api/option"
	"io"
	"strings"
	"time"
)

type GCSService struct {
	client     *storage.Client
	bucketName string
}

func NewGCSService(projectID, bucketName, credentialsPath string) (*GCSService, error) {
	ctx := context.Background()
	var client *storage.Client
	var err error

	if credentialsPath != "" {
		client, err = storage.NewClient(ctx, option.WithCredentialsFile(credentialsPath))
	} else {
		client, err = storage.NewClient(ctx)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to create GCS client: %w", err)
	}
	return &GCSService{
		client:     client,
		bucketName: bucketName,
	}, nil
}

func (g *GCSService) Upload(ctx context.Context, file io.Reader, filename string, contentType string) (string, string, error) {
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

	bucket := g.client.Bucket(g.bucketName)
	fName := fmt.Sprintf("%s-%s", uuid.New().String(), filename)
	object := bucket.Object(fName)

	writer := object.NewWriter(ctx)
	writer.ContentType = contentType

	if _, err := io.Copy(writer, file); err != nil {
		return "", "", fmt.Errorf("failed to upload file: %w", err)
	}
	if err := writer.Close(); err != nil {
		return "", "", fmt.Errorf("failed to close writer: %w", err)
	}

	attrs, err := object.Attrs(ctx)
	if err != nil {
		return "", "", fmt.Errorf("failed to get object attributes: %w", err)
	}
	if attrs.Size == 0 {
		return "", "", fmt.Errorf("uploaded file is empty")
	}

	publicURL := fmt.Sprintf("https://storage.googleapis.com/%s/%s", g.bucketName, fName)
	return publicURL, fName, nil
}

func (g *GCSService) Delete(ctx context.Context, fileURL string) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	filename := extractFilename(fileURL, g.bucketName)
	if filename == "" {
		return fmt.Errorf("invalid file URL: %s", fileURL)
	}

	bucket := g.client.Bucket(g.bucketName)
	object := bucket.Object(filename)
	if err := object.Delete(ctx); err != nil {
		return fmt.Errorf("failed to delete file: %w", err)
	}
	return nil
}

func (g *GCSService) GetSignedURL(ctx context.Context, filename string) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	url, err := g.client.Bucket(g.bucketName).SignedURL(filename, &storage.SignedURLOptions{
		Method:  "GET",
		Expires: time.Now().Add(15 * time.Minute),
	})
	if err != nil {
		return "", fmt.Errorf("failed to generate signed URL: %w", err)
	}
	return url, nil
}

// extractFilename extracts the object name from the public URL
func extractFilename(fileURL, bucketName string) string {
	prefix := fmt.Sprintf("https://storage.googleapis.com/%s/", bucketName)
	if strings.HasPrefix(fileURL, prefix) {
		return strings.TrimPrefix(fileURL, prefix)
	}
	return ""
}
