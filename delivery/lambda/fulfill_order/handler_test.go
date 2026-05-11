package main

import (
	"context"
	"testing"
	"time"

	"github.com/Josey34/goshop/domain/entity"
	"github.com/Josey34/goshop/domain/valueobject"
	"github.com/Josey34/goshop/repository/memory"
	ucorder "github.com/Josey34/goshop/usecase/order"
)

func seedOrder(repo *memory.OrderRepo, id string, status valueobject.OrderStatus) {
	repo.Create(context.Background(), &entity.Order{
		ID:         id,
		CustomerID: "c1",
		Items: []entity.OrderItem{
			{ID: "i1", OrderID: id, ProductID: "p1", Name: "Shirt", Price: valueobject.NewMoney(1000), Quantity: 1},
		},
		Total:     valueobject.NewMoney(1000),
		Status:    status,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
}

func TestFulfillOrderHandler_Handle(t *testing.T) {
	t.Run("shipped order transitions to delivered", func(t *testing.T) {
		repo := memory.NewOrderRepo()
		seedOrder(repo, "o1", valueobject.OrderStatusShipped)
		h := NewHandler(ucorder.NewGetOrderUseCase(repo), ucorder.NewUpdateOrderStatusUseCase(repo))

		resp, err := h.Handle(context.Background(), FulfillOrderEvent{OrderID: "o1"})
		if err != nil {
			t.Fatalf("unexpected err: %v", err)
		}
		if resp.Status != string(valueobject.OrderStatusDelivered) {
			t.Errorf("status=%s, want=DELIVERED", resp.Status)
		}
		if resp.OrderID != "o1" {
			t.Errorf("order_id=%s, want=o1", resp.OrderID)
		}
	})

	t.Run("order not found returns error", func(t *testing.T) {
		repo := memory.NewOrderRepo()
		h := NewHandler(ucorder.NewGetOrderUseCase(repo), ucorder.NewUpdateOrderStatusUseCase(repo))

		_, err := h.Handle(context.Background(), FulfillOrderEvent{OrderID: "missing"})
		if err == nil {
			t.Error("expected error for missing order")
		}
	})

	t.Run("invalid status transition returns error", func(t *testing.T) {
		repo := memory.NewOrderRepo()
		seedOrder(repo, "o1", valueobject.OrderStatusPending)
		h := NewHandler(ucorder.NewGetOrderUseCase(repo), ucorder.NewUpdateOrderStatusUseCase(repo))

		_, err := h.Handle(context.Background(), FulfillOrderEvent{OrderID: "o1"})
		if err == nil {
			t.Error("expected error for invalid status transition")
		}
	})
}
