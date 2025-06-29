# tradia-file-uploader-service

Tradia File Uploader Service
A Go library for uploading, deleting, and generating signed URLs for files on Google Cloud Storage (GCS), AWS S3, or DigitalOcean Spaces. The library abstracts provider-specific logic, allowing seamless switching between providers via configuration.
Installation
Add the library to your Go project:
go get github.com/Tradia-HQ/tradia-file-uploader-service@v0.1.0

Ensure you have Go 1.22 or later.
Usage
The library provides a FileUploader interface for uploading files, generating signed URLs, and deleting files. Use the services.NewUploader function to create an uploader instance based on configuration.
Example
package main

import (
"context"
"github.com/Tradia-HQ/tradia-file-uploader-service/internal/config"
"github.com/Tradia-HQ/tradia-file-uploader-service/internal/services"
"log"
"os"
)

func main() {
// Load configuration from environment variables or .env file
cfg, err := config.LoadConfig()
if err != nil {
log.Fatalf("Failed to load config: %v", err)
}

    // Initialize uploader (GCS, S3, or DigitalOcean Spaces)
    uploader, err := services.NewUploader(cfg)
    if err != nil {
        log.Fatalf("Failed to initialize uploader: %v", err)
    }

    // Open a file for upload
    file, err := os.Open("example.txt")
    if err != nil {
        log.Fatalf("Failed to open file: %v", err)
    }
    defer file.Close()

    // Upload file
    publicURL, objectName, err := uploader.Upload(context.Background(), file, "example.txt", "text/plain")
    if err != nil {
        log.Fatalf("Failed to upload: %v", err)
    }
    log.Printf("Uploaded to %s: %s (Object: %s)", cfg.Provider, publicURL, objectName)

    // Get signed URL
    signedURL, err := uploader.GetSignedURL(context.Background(), objectName)
    if err != nil {
        log.Fatalf("Failed to get signed URL: %v", err)
    }
    log.Printf("Signed URL: %s", signedURL)

    // Delete file
    if err := uploader.Delete(context.Background(), publicURL); err != nil {
        log.Fatalf("Failed to delete file: %v", err)
    }
    log.Println("File deleted successfully")
}

Configuration
Configure the library using environment variables or a .env file. Set STORAGE_PROVIDER to gcs, s3, or do to select the provider.
Example .env
STORAGE_PROVIDER=gcs
GCS_PROJECT_ID=your-project-id
GCS_BUCKET_NAME=your-bucket-name
GCS_CREDENTIALS_PATH=your-secret-path
AWS_REGION=us-east-1
AWS_BUCKET_NAME=your-bucket
AWS_ACCESS_KEY_ID=your-access-key
AWS_SECRET_ACCESS_KEY=your-secret-key
DO_REGION=nyc3
DO_BUCKET_NAME=your-do-bucket
DO_ACCESS_KEY_ID=your-do-access-key
DO_SECRET_ACCESS_KEY=your-do-secret-key
DO_ENDPOINT=nyc3.digitaloceanspaces.com


GCS: Requires a service account JSON file at GCS_CREDENTIALS_PATH.
S3/DigitalOcean: Requires access and secret keys.
Secrets: Store credentials securely (e.g., secrets manager) and exclude .env/secrets/ from version control.

Testing
Test the library locally by creating a test file and running the example above. For unit tests, mock the FileUploader interface:
package main

import (
"context"
"github.com/Tradia-HQ/tradia-file-uploader-service/internal/services/interfaces"
"io"
"testing"
)

type mockUploader struct{}

func (m *mockUploader) Upload(ctx context.Context, file io.Reader, filename, contentType string) (string, string, error) {
return "mock-url", "mock-object", nil
}
func (m *mockUploader) Delete(ctx context.Context, fileURL string) error {
return nil
}
func (m *mockUploader) GetSignedURL(ctx context.Context, filename string) (string, error) {
return "mock-signed-url", nil
}

func TestMockUploader(t *testing.T) {
uploader := &mockUploader{}
publicURL, objectName, err := uploader.Upload(context.Background(), nil, "test.txt", "text/plain")
if err != nil || publicURL == "" || objectName == "" {
t.Error("Expected successful upload")
}
}

Run tests:
go test ./...

Notes

File Types: Supports PNG, JPEG, TXT, PDF, and more (content type set by client).
File Size: Enforce limits in your application (e.g., 10MB) before calling Upload.
Versioning: Use tagged releases (e.g., v0.1.0) for stability.
Security: Add authentication in consuming services for production use.

Contributing
Submit pull requests to the staging branch. Ensure tests pass and include documentation updates.