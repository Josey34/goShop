package valueobject_test

import (
	"testing"

	"github.com/Josey34/goshop/domain/valueobject"
)

func TestNewImageURL(t *testing.T) {
	cases := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"valid http", "http://example.com/image.jpg", false},
		{"valid https", "https://s3.amazonaws.com/bucket/image.jpg", false},
		{"empty", "", true},
		{"no scheme", "example.com/image.jpg", true},
		{"ftp scheme", "ftp://example.com/image.jpg", true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			url, err := valueobject.NewImageURL(tc.input)
			if (err != nil) != tc.wantErr {
				t.Errorf("got err=%v, wantErr=%v", err, tc.wantErr)
			}
			if !tc.wantErr && url.Value() != tc.input {
				t.Errorf("got=%s, want=%s", url.Value(), tc.input)
			}
		})
	}
}
