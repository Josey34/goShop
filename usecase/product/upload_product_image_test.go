package product_test

import (
	"context"
	"errors"
	"testing"

	ucproduct "github.com/Josey34/goshop/usecase/product"
)

type mockFileStorage struct {
	uploadResult string
	uploadErr    error
}

func (m *mockFileStorage) Upload(ctx context.Context, key string, data []byte, contentType string) (string, error) {
	return m.uploadResult, m.uploadErr
}
func (m *mockFileStorage) GetPresignedURL(ctx context.Context, key string) (string, error) {
	return "", nil
}
func (m *mockFileStorage) Delete(ctx context.Context, key string) error { return nil }

func TestUploadProductImageUsecase_Execute(t *testing.T) {
	t.Run("uploads and returns url", func(t *testing.T) {
		repo := &mockProductRepo{findResult: newProduct("p1", "Shirt", 1000)}
		storage := &mockFileStorage{uploadResult: "https://s3.example.com/products/p1/image.jpg"}
		uc := ucproduct.NewUploadProductImageUsecase(repo, storage)

		url, err := uc.Execute(context.Background(), ucproduct.UploadProductImage{
			ProductID:   "p1",
			Data:        []byte("imgdata"),
			ContentType: "image/jpeg",
		})
		if err != nil {
			t.Fatalf("unexpected err: %v", err)
		}
		if url != "https://s3.example.com/products/p1/image.jpg" {
			t.Errorf("got url=%s", url)
		}
	})

	t.Run("product not found returns error", func(t *testing.T) {
		repo := &mockProductRepo{findErr: errors.New("not found")}
		uc := ucproduct.NewUploadProductImageUsecase(repo, &mockFileStorage{})

		_, err := uc.Execute(context.Background(), ucproduct.UploadProductImage{ProductID: "bad"})
		if err == nil {
			t.Error("expected error for missing product")
		}
	})

	t.Run("storage error returns error", func(t *testing.T) {
		repo := &mockProductRepo{findResult: newProduct("p1", "Shirt", 1000)}
		storage := &mockFileStorage{uploadErr: errors.New("s3 error")}
		uc := ucproduct.NewUploadProductImageUsecase(repo, storage)

		_, err := uc.Execute(context.Background(), ucproduct.UploadProductImage{
			ProductID: "p1", Data: []byte("img"), ContentType: "image/jpeg",
		})
		if err == nil {
			t.Error("expected error for storage failure")
		}
	})
}
