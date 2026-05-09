package entity_test

import (
	"testing"
	"time"

	"github.com/Josey34/goshop/domain/entity"
	"github.com/Josey34/goshop/domain/valueobject"
)

func TestPayment_Fields(t *testing.T) {
	p := &entity.Payment{
		ID:        "pay1",
		OrderID:   "o1",
		Amount:    valueobject.NewMoney(5000),
		Status:    valueobject.PaymentStatusPaid,
		Method:    "credit_card",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if p.ID != "pay1" {
		t.Errorf("got id=%s, want=pay1", p.ID)
	}
	if p.Amount.Value() != 5000 {
		t.Errorf("got amount=%d, want=5000", p.Amount.Value())
	}
	if p.Status != valueobject.PaymentStatusPaid {
		t.Errorf("got status=%s, want=PAID", p.Status)
	}
}
