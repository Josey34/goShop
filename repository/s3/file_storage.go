package s3

import (
	"bytes"
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type FileStorage struct {
	client     *s3.Client
	bucketName string
	region     string
}

func NewFileStorage(client *s3.Client, bucketName, region string) *FileStorage {
	return &FileStorage{
		client:     client,
		bucketName: bucketName,
		region:     region,
	}
}

func (fs *FileStorage) Upload(ctx context.Context, key string, body []byte, contentType string) (string, error) {
	_, err := fs.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(fs.bucketName),
		Key:         aws.String(key),
		Body:        bytes.NewReader(body),
		ContentType: aws.String(contentType),
	})

	if err != nil {
		return "", fmt.Errorf("failed to upload file: %w", err)
	}

	presignedURL, err := fs.GetPresignedURL(ctx, key)
	if err != nil {
		return "", fmt.Errorf("failed to generate presigned URL: %w", err)
	}

	return presignedURL, nil
}

func (fs *FileStorage) GetPresignedURL(ctx context.Context, key string) (string, error) {
	presigner := s3.NewPresignClient(fs.client)

	presignedURL, err := presigner.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(fs.bucketName),
		Key:    aws.String(key),
	}, s3.WithPresignExpires(15*time.Minute))

	if err != nil {
		return "", fmt.Errorf("failed to generate presigned URL: %w", err)
	}

	return presignedURL.URL, nil
}

func (fs *FileStorage) Delete(ctx context.Context, key string) error {
	_, err := fs.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(fs.bucketName),
		Key:    aws.String(key),
	})

	if err != nil {
		return fmt.Errorf("failed to delete file: %w", err)
	}

	return nil
}
