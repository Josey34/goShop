package product_test

import (
	"context"
	"errors"
	"testing"

	"github.com/Josey34/goshop/domain/entity"
	"github.com/Josey34/goshop/domain/valueobject"
	ucproduct "github.com/Josey34/goshop/usecase/product"
)

func TestListProductsUseCase_Execute(t *testing.T) {
	t.Run("returns products", func(t *testing.T) {
		repo := &mockProductRepo{
			listResult: []*entity.Product{
				newProduct("p1", "Shirt", 1000),
				newProduct("p2", "Pants", 2000),
			},
		}
		uc := ucproduct.NewListProductUseCase(repo)

		products, err := uc.Execute(context.Background(), valueobject.NewPagination(1, 10))
		if err != nil {
			t.Fatalf("unexpected err: %v", err)
		}
		if len(products) != 2 {
			t.Errorf("got %d products, want=2", len(products))
		}
	})

	t.Run("repo error propagates", func(t *testing.T) {
		repo := &mockProductRepo{listErr: errors.New("db error")}
		uc := ucproduct.NewListProductUseCase(repo)

		_, err := uc.Execute(context.Background(), valueobject.NewPagination(1, 10))
		if err == nil {
			t.Error("expected error")
		}
	})
}
