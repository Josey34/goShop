package memory

import (
	"context"
	"fmt"
	"sync"
)

type FileStorage struct {
	mu    sync.RWMutex
	files map[string][]byte
}

func NewFileStorage() *FileStorage {
	return &FileStorage{
		files: make(map[string][]byte),
	}
}

func (s *FileStorage) Upload(ctx context.Context, key string, data []byte, contentType string) (string, error) {
	s.files[key] = data

	return fmt.Sprintf("memory://%s", key), nil
}

func (s *FileStorage) GetPresignedURL(ctx context.Context, key string) (string, error) {
	return fmt.Sprintf("memory://%s", key), nil
}

func (s *FileStorage) Delete(ctx context.Context, key string) error {
	delete(s.files, key)

	return nil
}
