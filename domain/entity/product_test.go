package entity_test

import (
	"testing"

	"github.com/Josey34/goshop/domain/entity"
	"github.com/Josey34/goshop/domain/valueobject"
)

func TestProduct_ReduceStock(t *testing.T) {
	cases := []struct {
		name      string
		stock     int
		reduce    int
		wantErr   bool
		wantStock int
	}{
		{"sufficient stock", 10, 3, false, 7},
		{"exact stock", 5, 5, false, 0},
		{"insufficient stock", 2, 5, true, 2},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			p := &entity.Product{ID: "p1", Stock: tc.stock}
			err := p.ReduceStock(tc.reduce)
			if (err != nil) != tc.wantErr {
				t.Errorf("got err=%v, wantErr=%v", err, tc.wantErr)
			}
			if p.Stock != tc.wantStock {
				t.Errorf("got stock=%d, want=%d", p.Stock, tc.wantStock)
			}
		})
	}
}

func TestProduct_UpdatePrice(t *testing.T) {
	cases := []struct {
		name    string
		price   int64
		wantErr bool
	}{
		{"valid price", 1000, false},
		{"zero price", 0, true},
		{"negative price", -500, true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			p := &entity.Product{ID: "p1", Price: valueobject.NewMoney(500)}
			err := p.UpdatePrice(valueobject.NewMoney(tc.price))
			if (err != nil) != tc.wantErr {
				t.Errorf("got err=%v, wantErr=%v", err, tc.wantErr)
			}
			if !tc.wantErr && p.Price.Value() != tc.price {
				t.Errorf("got price=%d, want=%d", p.Price.Value(), tc.price)
			}
		})
	}
}
