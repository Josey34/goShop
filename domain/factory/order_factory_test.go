package factory_test

import (
	"testing"

	"github.com/Josey34/goshop/domain/factory"
	"github.com/Josey34/goshop/domain/valueobject"
)

func TestNewOrder(t *testing.T) {
	t.Run("valid order calculates total and attaches event", func(t *testing.T) {
		order, err := factory.NewOrder("o1", factory.CreateOrderInput{
			CustomerID: "c1",
			Items: []factory.CreateOrderItemInput{
				{ProductID: "p1", Name: "Shirt", Price: valueobject.NewMoney(1000), Quantity: 2},
				{ProductID: "p2", Name: "Pants", Price: valueobject.NewMoney(500), Quantity: 1},
			},
		})

		if err != nil {
			t.Fatalf("unexpected err: %v", err)
		}
		if order.ID != "o1" {
			t.Errorf("got id=%s, want=o1", order.ID)
		}
		if order.Total.Value() != 2500 {
			t.Errorf("got total=%d, want=2500", order.Total.Value())
		}
		if order.Status != valueobject.OrderStatusPending {
			t.Errorf("got status=%s, want=PENDING", order.Status)
		}
		events := order.PullEvents()
		if len(events) != 1 {
			t.Errorf("got %d events, want=1", len(events))
		}
	})

	t.Run("empty items returns error", func(t *testing.T) {
		_, err := factory.NewOrder("o1", factory.CreateOrderInput{
			CustomerID: "c1",
			Items:      []factory.CreateOrderItemInput{},
		})
		if err == nil {
			t.Error("expected error for empty items")
		}
	})
}
