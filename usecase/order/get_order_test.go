package order_test

import (
	"context"
	"errors"
	"testing"

	"github.com/Josey34/goshop/domain/valueobject"
	ucorder "github.com/Josey34/goshop/usecase/order"
)

func TestGetOrderUseCase_Execute(t *testing.T) {
	t.Run("found returns order", func(t *testing.T) {
		repo := &mockOrderRepo{findResult: makeOrder("o1", valueobject.OrderStatusPending)}
		uc := ucorder.NewGetOrderUseCase(repo)

		o, err := uc.Execute(context.Background(), "o1")
		if err != nil {
			t.Fatalf("unexpected err: %v", err)
		}
		if o.ID != "o1" {
			t.Errorf("got id=%s, want=o1", o.ID)
		}
	})

	t.Run("not found propagates error", func(t *testing.T) {
		repo := &mockOrderRepo{findErr: errors.New("not found")}
		uc := ucorder.NewGetOrderUseCase(repo)

		_, err := uc.Execute(context.Background(), "missing")
		if err == nil {
			t.Error("expected error")
		}
	})
}
