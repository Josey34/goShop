package product_test

import (
	"context"
	"errors"
	"testing"

	ucproduct "github.com/Josey34/goshop/usecase/product"
)

func TestUpdateProductUseCase_Execute(t *testing.T) {
	t.Run("updates and returns product", func(t *testing.T) {
		repo := &mockProductRepo{findResult: newProduct("p1", "Old", 500)}
		uc := ucproduct.NewUpdateProductUseCase(repo)

		p, err := uc.Execute(context.Background(), ucproduct.UpdateProductInput{
			ID: "p1", Name: "New", Price: 2000, Stock: 5,
		})
		if err != nil {
			t.Fatalf("unexpected err: %v", err)
		}
		if p.Name != "New" {
			t.Errorf("got name=%s, want=New", p.Name)
		}
		if p.Price.Value() != 2000 {
			t.Errorf("got price=%d, want=2000", p.Price.Value())
		}
	})

	t.Run("not found returns error", func(t *testing.T) {
		repo := &mockProductRepo{findErr: errors.New("not found")}
		uc := ucproduct.NewUpdateProductUseCase(repo)

		_, err := uc.Execute(context.Background(), ucproduct.UpdateProductInput{ID: "bad"})
		if err == nil {
			t.Error("expected error for missing product")
		}
	})
}
