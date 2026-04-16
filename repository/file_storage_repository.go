package repository

import "context"

type FileStorage interface {
	Upload(ctx context.Context, key string, data []byte, contentType string) (string, error)
	GetPresignedURL(ctx context.Context, key string) (string, error)
	Delete(ctx context.Context, key string) error
}
