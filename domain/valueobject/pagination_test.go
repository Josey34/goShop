package valueobject_test

import (
	"testing"

	"github.com/Josey34/goshop/domain/valueobject"
)

func TestNewPagination(t *testing.T) {
	cases := []struct {
		name      string
		page      int
		limit     int
		wantPage  int
		wantLimit int
	}{
		{"valid", 2, 20, 2, 20},
		{"page below 1 defaults to 1", 0, 10, 1, 10},
		{"negative page defaults to 1", -5, 10, 1, 10},
		{"limit below 1 defaults to 10", 1, 0, 1, 10},
		{"limit above 100 capped to 100", 1, 200, 1, 100},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			p := valueobject.NewPagination(tc.page, tc.limit)
			if p.Page != tc.wantPage {
				t.Errorf("got page=%d, want=%d", p.Page, tc.wantPage)
			}
			if p.Limit != tc.wantLimit {
				t.Errorf("got limit=%d, want=%d", p.Limit, tc.wantLimit)
			}
		})
	}
}

func TestPagination_Offset(t *testing.T) {
	cases := []struct {
		page       int
		limit      int
		wantOffset int
	}{
		{1, 10, 0},
		{2, 10, 10},
		{3, 20, 40},
	}

	for _, tc := range cases {
		p := valueobject.NewPagination(tc.page, tc.limit)
		if got := p.Offset(); got != tc.wantOffset {
			t.Errorf("page=%d limit=%d: got offset=%d, want=%d", tc.page, tc.limit, got, tc.wantOffset)
		}
	}
}
