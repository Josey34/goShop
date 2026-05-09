package entity_test

import (
	"testing"

	"github.com/Josey34/goshop/domain/entity"
	"github.com/Josey34/goshop/domain/valueobject"
)

func TestOrderItem_LineTotal(t *testing.T) {
	cases := []struct {
		name     string
		price    int64
		qty      int
		wantTotal int64
	}{
		{"single item", 1000, 1, 1000},
		{"multiple qty", 500, 3, 1500},
		{"zero qty", 1000, 0, 0},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			item := entity.OrderItem{
				Price:    valueobject.NewMoney(tc.price),
				Quantity: tc.qty,
			}
			if got := item.LineTotal().Value(); got != tc.wantTotal {
				t.Errorf("got=%d, want=%d", got, tc.wantTotal)
			}
		})
	}
}
