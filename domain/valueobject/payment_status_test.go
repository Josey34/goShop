package valueobject_test

import (
	"testing"

	"github.com/Josey34/goshop/domain/valueobject"
)

func TestPaymentStatus_Values(t *testing.T) {
	cases := []struct {
		status valueobject.PaymentStatus
		want   string
	}{
		{valueobject.PaymentStatusPending, "PENDING"},
		{valueobject.PaymentStatusPaid, "PAID"},
		{valueobject.PaymentStatusFailed, "FAILED"},
		{valueobject.PaymentStatusRefunded, "REFUNDED"},
	}

	for _, tc := range cases {
		if string(tc.status) != tc.want {
			t.Errorf("got=%s, want=%s", tc.status, tc.want)
		}
	}
}
