package validator_test

import (
	"testing"

	"github.com/Josey34/goshop/pkg/validator"
)

func TestRequiredString(t *testing.T) {
	cases := []struct {
		name    string
		value   string
		field   string
		wantErr bool
	}{
		{"non-empty value", "hello", "name", false},
		{"empty value", "", "name", true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := validator.RequiredString(tc.value, tc.field)
			if (err != nil) != tc.wantErr {
				t.Errorf("got err=%v, wantErr=%v", err, tc.wantErr)
			}
		})
	}
}
