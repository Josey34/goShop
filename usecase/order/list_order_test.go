package order_test

import (
	"context"
	"errors"
	"testing"

	"github.com/Josey34/goshop/domain/entity"
	"github.com/Josey34/goshop/domain/valueobject"
	ucorder "github.com/Josey34/goshop/usecase/order"
)

func TestListOrderUseCase_Execute(t *testing.T) {
	t.Run("returns customer orders", func(t *testing.T) {
		repo := &mockOrderRepo{
			listResult: []*entity.Order{
				makeOrder("o1", valueobject.OrderStatusPending),
				makeOrder("o2", valueobject.OrderStatusConfirmed),
			},
		}
		uc := ucorder.NewListProductUseCase(repo)

		orders, err := uc.Execute(context.Background(), "c1", valueobject.NewPagination(1, 10))
		if err != nil {
			t.Fatalf("unexpected err: %v", err)
		}
		if len(orders) != 2 {
			t.Errorf("got %d orders, want=2", len(orders))
		}
	})

	t.Run("repo error propagates", func(t *testing.T) {
		repo := &mockOrderRepo{listErr: errors.New("db error")}
		uc := ucorder.NewListProductUseCase(repo)

		_, err := uc.Execute(context.Background(), "c1", valueobject.NewPagination(1, 10))
		if err == nil {
			t.Error("expected error")
		}
	})
}
