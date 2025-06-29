package interfaces

import (
	"context"
	"io"
)

type FileUploader interface {
	Upload(ctx context.Context, file io.Reader, filename string, contentType string) (string, string, error) // Returns public URL, object name, error
	Delete(ctx context.Context, fileURL string) error
	GetSignedURL(ctx context.Context, filename string) (string, error)
}
