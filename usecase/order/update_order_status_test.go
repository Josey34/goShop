package order_test

import (
	"context"
	"errors"
	"testing"

	"github.com/Josey34/goshop/domain/valueobject"
	ucorder "github.com/Josey34/goshop/usecase/order"
)

func TestUpdateOrderStatusUseCase_Execute(t *testing.T) {
	t.Run("valid transition updates status", func(t *testing.T) {
		repo := &mockOrderRepo{findResult: makeOrder("o1", valueobject.OrderStatusPending)}
		uc := ucorder.NewUpdateOrderStatusUseCase(repo)

		o, err := uc.Execute(context.Background(), ucorder.UpdateOrderStatusInput{
			ID:     "o1",
			Status: valueobject.OrderStatusConfirmed,
		})
		if err != nil {
			t.Fatalf("unexpected err: %v", err)
		}
		if o.Status != valueobject.OrderStatusConfirmed {
			t.Errorf("got status=%s, want=CONFIRMED", o.Status)
		}
	})

	t.Run("invalid transition returns error", func(t *testing.T) {
		repo := &mockOrderRepo{findResult: makeOrder("o1", valueobject.OrderStatusDelivered)}
		uc := ucorder.NewUpdateOrderStatusUseCase(repo)

		_, err := uc.Execute(context.Background(), ucorder.UpdateOrderStatusInput{
			ID:     "o1",
			Status: valueobject.OrderStatusCancelled,
		})
		if err == nil {
			t.Error("expected error for invalid transition")
		}
	})

	t.Run("not found returns error", func(t *testing.T) {
		repo := &mockOrderRepo{findErr: errors.New("not found")}
		uc := ucorder.NewUpdateOrderStatusUseCase(repo)

		_, err := uc.Execute(context.Background(), ucorder.UpdateOrderStatusInput{ID: "bad"})
		if err == nil {
			t.Error("expected error")
		}
	})
}
