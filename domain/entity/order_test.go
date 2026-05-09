package entity_test

import (
	"testing"

	"github.com/Josey34/goshop/domain/entity"
	"github.com/Josey34/goshop/domain/valueobject"
)

func TestOrder_CalculateTotal(t *testing.T) {
	o := &entity.Order{
		Items: []entity.OrderItem{
			{Price: valueobject.NewMoney(1000), Quantity: 2},
			{Price: valueobject.NewMoney(500), Quantity: 1},
		},
	}
	o.CalculateTotal()
	if o.Total.Value() != 2500 {
		t.Errorf("got total=%d, want=2500", o.Total.Value())
	}
}

func TestOrder_TransitionTo(t *testing.T) {
	cases := []struct {
		name    string
		from    valueobject.OrderStatus
		to      valueobject.OrderStatus
		wantErr bool
	}{
		{"pending to confirmed", valueobject.OrderStatusPending, valueobject.OrderStatusConfirmed, false},
		{"pending to cancelled", valueobject.OrderStatusPending, valueobject.OrderStatusCancelled, false},
		{"pending to shipped", valueobject.OrderStatusPending, valueobject.OrderStatusShipped, true},
		{"delivered to cancelled", valueobject.OrderStatusDelivered, valueobject.OrderStatusCancelled, true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			o := &entity.Order{Status: tc.from}
			err := o.TransitionTo(tc.to)
			if (err != nil) != tc.wantErr {
				t.Errorf("got err=%v, wantErr=%v", err, tc.wantErr)
			}
			if !tc.wantErr && o.Status != tc.to {
				t.Errorf("got status=%s, want=%s", o.Status, tc.to)
			}
		})
	}
}

func TestOrder_PullEvents(t *testing.T) {
	o := &entity.Order{}
	if events := o.PullEvents(); len(events) != 0 {
		t.Errorf("expected empty events, got %d", len(events))
	}
}
