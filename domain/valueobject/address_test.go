package valueobject_test

import (
	"testing"

	"github.com/Josey34/goshop/domain/valueobject"
)

func TestNewAddress(t *testing.T) {
	cases := []struct {
		name       string
		street     string
		city       string
		province   string
		postalCode string
		wantErr    bool
	}{
		{"valid address", "Jl. Sudirman", "Jakarta", "DKI Jakarta", "10220", false},
		{"all empty returns empty address", "", "", "", "", false},
		{"missing street", "", "Jakarta", "DKI Jakarta", "10220", true},
		{"missing city", "Jl. A", "", "DKI Jakarta", "10220", true},
		{"invalid postal code length", "Jl. A", "Jakarta", "DKI", "123", true},
		{"postal code too long", "Jl. A", "Jakarta", "DKI", "123456", true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			addr, err := valueobject.NewAddress(tc.street, tc.city, tc.province, tc.postalCode)
			if (err != nil) != tc.wantErr {
				t.Errorf("got err=%v, wantErr=%v", err, tc.wantErr)
			}
			if !tc.wantErr && tc.street != "" && addr.Street != tc.street {
				t.Errorf("got street=%s, want=%s", addr.Street, tc.street)
			}
		})
	}
}
