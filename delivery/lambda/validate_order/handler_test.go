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

func seedOrder(repo *memory.OrderRepo, id string, items []entity.OrderItem) {
	repo.Create(context.Background(), &entity.Order{
		ID:         id,
		CustomerID: "c1",
		Items:      items,
		Status:     valueobject.OrderStatusPending,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	})
}

func TestValidateOrderHandler_Handle(t *testing.T) {
	t.Run("valid order with items", func(t *testing.T) {
		repo := memory.NewOrderRepo()
		seedOrder(repo, "o1", []entity.OrderItem{
			{ID: "i1", OrderID: "o1", ProductID: "p1", Name: "Shirt", Price: valueobject.NewMoney(1000), Quantity: 1},
		})
		h := NewHandler(ucorder.NewGetOrderUseCase(repo))

		resp, err := h.Handle(context.Background(), ValidateOrderEvent{OrderID: "o1"})
		if err != nil {
			t.Fatalf("unexpected err: %v", err)
		}
		if !resp.IsValid {
			t.Errorf("expected IsValid=true, got errors: %v", resp.Errors)
		}
	})

	t.Run("order not found returns error", func(t *testing.T) {
		repo := memory.NewOrderRepo()
		h := NewHandler(ucorder.NewGetOrderUseCase(repo))

		resp, err := h.Handle(context.Background(), ValidateOrderEvent{OrderID: "missing"})
		if err == nil {
			t.Error("expected error for missing order")
		}
		if resp.IsValid {
			t.Error("expected IsValid=false")
		}
	})

	t.Run("order with no items is invalid", func(t *testing.T) {
		repo := memory.NewOrderRepo()
		seedOrder(repo, "o1", []entity.OrderItem{})
		h := NewHandler(ucorder.NewGetOrderUseCase(repo))

		resp, err := h.Handle(context.Background(), ValidateOrderEvent{OrderID: "o1"})
		if err != nil {
			t.Fatalf("unexpected err: %v", err)
		}
		if resp.IsValid {
			t.Error("expected IsValid=false for empty items")
		}
	})
}
