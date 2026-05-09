package product_test

import (
	"context"
	"errors"
	"testing"

	ucproduct "github.com/Josey34/goshop/usecase/product"
)

func TestDeleteProductUseCase_Execute(t *testing.T) {
	t.Run("found deletes successfully", func(t *testing.T) {
		repo := &mockProductRepo{findResult: newProduct("p1", "Shirt", 1000)}
		uc := ucproduct.NewDeleteProductUseCase(repo)

		err := uc.Execute(context.Background(), "p1")
		if err != nil {
			t.Fatalf("unexpected err: %v", err)
		}
	})

	t.Run("not found returns error", func(t *testing.T) {
		repo := &mockProductRepo{findErr: errors.New("not found")}
		uc := ucproduct.NewDeleteProductUseCase(repo)

		err := uc.Execute(context.Background(), "missing")
		if err == nil {
			t.Error("expected error for missing product")
		}
	})
}
