package s3

import (
	"context"
	"io"
)

// Storage defines the S3 operations used by the application.
type Storage interface {
	Upload(ctx context.Context, key string, r io.Reader, contentType string) (string, error)
	Delete(ctx context.Context, key string) error
	GetURL(key string) string
}
