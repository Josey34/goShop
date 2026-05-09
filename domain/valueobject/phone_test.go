package valueobject_test

import (
	"testing"

	"github.com/Josey34/goshop/domain/valueobject"
)

func TestNewPhone(t *testing.T) {
	cases := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"valid 08xx", "08123456789", false},
		{"valid +62xx", "+628123456789", false},
		{"empty", "", true},
		{"invalid prefix", "628123456789", true},
		{"random string", "hello", true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			phone, err := valueobject.NewPhone(tc.input)
			if (err != nil) != tc.wantErr {
				t.Errorf("got err=%v, wantErr=%v", err, tc.wantErr)
			}
			if !tc.wantErr && phone.Value() != tc.input {
				t.Errorf("got=%s, want=%s", phone.Value(), tc.input)
			}
		})
	}
}
