package valueobject_test

import (
	"testing"

	"github.com/Josey34/goshop/domain/valueobject"
)

func TestMoney_Add(t *testing.T) {
	a := valueobject.NewMoney(300)
	b := valueobject.NewMoney(200)
	if got := a.Add(b).Value(); got != 500 {
		t.Errorf("got=%d, want=500", got)
	}
}

func TestMoney_Subtract(t *testing.T) {
	a := valueobject.NewMoney(500)
	b := valueobject.NewMoney(200)
	if got := a.Subtract(b).Value(); got != 300 {
		t.Errorf("got=%d, want=300", got)
	}
}

func TestMoney_Multiply(t *testing.T) {
	m := valueobject.NewMoney(100)
	if got := m.Multiply(5).Value(); got != 500 {
		t.Errorf("got=%d, want=500", got)
	}
}

func TestMoney_IsNegative(t *testing.T) {
	if !valueobject.NewMoney(-1).IsNegative() {
		t.Error("expected negative")
	}
	if valueobject.NewMoney(1).IsNegative() {
		t.Error("expected not negative")
	}
}

func TestMoney_IsZero(t *testing.T) {
	if !valueobject.NewMoney(0).IsZero() {
		t.Error("expected zero")
	}
	if valueobject.NewMoney(1).IsZero() {
		t.Error("expected not zero")
	}
}
