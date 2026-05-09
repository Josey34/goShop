package valueobject_test

import (
	"testing"

	"github.com/Josey34/goshop/domain/valueobject"
)

func TestNewEmail(t *testing.T) {
	cases := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"valid email", "user@example.com", false},
		{"empty email", "", true},
		{"missing at sign", "userexample.com", true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			email, err := valueobject.NewEmail(tc.input)
			if (err != nil) != tc.wantErr {
				t.Errorf("got err=%v, wantErr=%v", err, tc.wantErr)
			}
			if !tc.wantErr && email.Value() != tc.input {
				t.Errorf("got=%s, want=%s", email.Value(), tc.input)
			}
		})
	}
}
