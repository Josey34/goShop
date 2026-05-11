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

func seedOrder(repo *memory.OrderRepo, id string, total int64) {
	repo.Create(context.Background(), &entity.Order{
		ID:         id,
		CustomerID: "c1",
		Items: []entity.OrderItem{
			{ID: "i1", OrderID: id, ProductID: "p1", Name: "Shirt", Price: valueobject.NewMoney(total), Quantity: 1},
		},
		Total:     valueobject.NewMoney(total),
		Status:    valueobject.OrderStatusPending,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
}

func TestCalculateTotalHandler_Handle(t *testing.T) {
	t.Run("returns correct subtotal, tax, and total", func(t *testing.T) {
		repo := memory.NewOrderRepo()
		seedOrder(repo, "o1", 10000)
		h := NewHandler(ucorder.NewGetOrderUseCase(repo))

		resp, err := h.Handle(context.Background(), CalculateTotalEvent{OrderID: "o1"})
		if err != nil {
			t.Fatalf("unexpected err: %v", err)
		}
		if resp.Subtotal != 10000 {
			t.Errorf("subtotal=%d, want=10000", resp.Subtotal)
		}
		if resp.Tax != 1100 {
			t.Errorf("tax=%d, want=1100", resp.Tax)
		}
		if resp.Total != 11100 {
			t.Errorf("total=%d, want=11100", resp.Total)
		}
		if resp.OrderID != "o1" {
			t.Errorf("order_id=%s, want=o1", resp.OrderID)
		}
	})

	t.Run("order not found returns error", func(t *testing.T) {
		repo := memory.NewOrderRepo()
		h := NewHandler(ucorder.NewGetOrderUseCase(repo))

		_, err := h.Handle(context.Background(), CalculateTotalEvent{OrderID: "missing"})
		if err == nil {
			t.Error("expected error for missing order")
		}
	})
}
