package valueobject_test

import (
	"testing"

	"github.com/Josey34/goshop/domain/valueobject"
)

func TestOrderStatus_CanTransitionTo(t *testing.T) {
	cases := []struct {
		name string
		from valueobject.OrderStatus
		to   valueobject.OrderStatus
		want bool
	}{
		{"pending → confirmed", valueobject.OrderStatusPending, valueobject.OrderStatusConfirmed, true},
		{"pending → cancelled", valueobject.OrderStatusPending, valueobject.OrderStatusCancelled, true},
		{"pending → shipped (invalid)", valueobject.OrderStatusPending, valueobject.OrderStatusShipped, false},
		{"confirmed → processing", valueobject.OrderStatusConfirmed, valueobject.OrderStatusProcessing, true},
		{"confirmed → cancelled", valueobject.OrderStatusConfirmed, valueobject.OrderStatusCancelled, true},
		{"processing → shipped", valueobject.OrderStatusProcessing, valueobject.OrderStatusShipped, true},
		{"processing → confirmed (invalid)", valueobject.OrderStatusProcessing, valueobject.OrderStatusConfirmed, false},
		{"shipped → delivered", valueobject.OrderStatusShipped, valueobject.OrderStatusDelivered, true},
		{"delivered → cancelled (terminal)", valueobject.OrderStatusDelivered, valueobject.OrderStatusCancelled, false},
		{"cancelled → pending (terminal)", valueobject.OrderStatusCancelled, valueobject.OrderStatusPending, false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.from.CanTransitionTo(tc.to)
			if got != tc.want {
				t.Errorf("got=%v, want=%v", got, tc.want)
			}
		})
	}
}
