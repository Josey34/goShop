package main

import (
	"context"
	"testing"
)

func TestProcessPaymentHandler_Handle(t *testing.T) {
	t.Run("returns valid payment status", func(t *testing.T) {
		h := NewHandler()

		resp, err := h.Handle(context.Background(), ProcessPaymentEvent{OrderID: "o1", Amount: 5000})
		if err != nil {
			t.Fatalf("unexpected err: %v", err)
		}
		if resp.OrderID != "o1" {
			t.Errorf("order_id=%s, want=o1", resp.OrderID)
		}
		if resp.PaymentStatus != "success" && resp.PaymentStatus != "failed" {
			t.Errorf("unexpected payment_status=%s", resp.PaymentStatus)
		}
		if resp.Message == "" {
			t.Error("expected non-empty message")
		}
	})
}
