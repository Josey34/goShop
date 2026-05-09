package product_test

import (
	"context"
	"errors"
	"testing"

	ucproduct "github.com/Josey34/goshop/usecase/product"
)

func TestGetProductUseCase_Execute(t *testing.T) {
	t.Run("found returns product", func(t *testing.T) {
		repo := &mockProductRepo{findResult: newProduct("p1", "Shirt", 1000)}
		uc := ucproduct.NewGetProductUseCase(repo)

		p, err := uc.Execute(context.Background(), "p1")
		if err != nil {
			t.Fatalf("unexpected err: %v", err)
		}
		if p.ID != "p1" {
			t.Errorf("got id=%s, want=p1", p.ID)
		}
	})

	t.Run("not found propagates error", func(t *testing.T) {
		repo := &mockProductRepo{findErr: errors.New("not found")}
		uc := ucproduct.NewGetProductUseCase(repo)

		_, err := uc.Execute(context.Background(), "missing")
		if err == nil {
			t.Error("expected error for missing product")
		}
	})
}
