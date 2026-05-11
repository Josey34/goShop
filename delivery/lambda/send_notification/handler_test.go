package main

import (
	"context"
	"testing"
)

func TestSendNotificationHandler_Handle(t *testing.T) {
	t.Run("always returns sent=true", func(t *testing.T) {
		h := NewHandler()

		resp, err := h.Handle(context.Background(), SendNotificationEvent{
			OrderID:    "o1",
			CustomerID: "c1",
			Message:    "your order has been placed",
		})
		if err != nil {
			t.Fatalf("unexpected err: %v", err)
		}
		if !resp.Sent {
			t.Error("expected Sent=true")
		}
		if resp.OrderID != "o1" {
			t.Errorf("order_id=%s, want=o1", resp.OrderID)
		}
	})
}
