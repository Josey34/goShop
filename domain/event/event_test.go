package event_test

import (
	"testing"

	"github.com/Josey34/goshop/domain/event"
)

func TestOrderCreated(t *testing.T) {
	e := event.NewOrderCreated("o1", "c1", 5000)
	if e.EventName() != "order.created" {
		t.Errorf("got=%s, want=order.created", e.EventName())
	}
	if e.OrderID != "o1" {
		t.Errorf("got orderID=%s, want=o1", e.OrderID)
	}
	if e.Total != 5000 {
		t.Errorf("got total=%d, want=5000", e.Total)
	}
	if e.OccurredAt().IsZero() {
		t.Error("expected non-zero OccurredAt")
	}
}

func TestOrderStatusChanged(t *testing.T) {
	e := event.NewOrderStatusChanged("o1", "PENDING", "CONFIRMED")
	if e.EventName() != "order.status_changed" {
		t.Errorf("got=%s, want=order.status_changed", e.EventName())
	}
	if e.From != "PENDING" || e.To != "CONFIRMED" {
		t.Errorf("got from=%s to=%s", e.From, e.To)
	}
}

func TestPaymentProcessed(t *testing.T) {
	e := event.NewPaymentProcessed("o1", "pay1", 5000, "PAID")
	if e.EventName() != "payment.processed" {
		t.Errorf("got=%s, want=payment.processed", e.EventName())
	}
	if e.Amount != 5000 {
		t.Errorf("got amount=%d, want=5000", e.Amount)
	}
	if e.Status != "PAID" {
		t.Errorf("got status=%s, want=PAID", e.Status)
	}
}

func TestStockReduced(t *testing.T) {
	e := event.NewStockReduced("p1", 3, 7)
	if e.EventName() != "stock.reduced" {
		t.Errorf("got=%s, want=stock.reduced", e.EventName())
	}
	if e.Quantity != 3 || e.Remaining != 7 {
		t.Errorf("got qty=%d remaining=%d", e.Quantity, e.Remaining)
	}
}
