package date_test

import (
	"testing"
	"time"

	"github.com/dylan-bourque/types/date"
)

func TestXxx(t *testing.T) {
	now := time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC)
	srcYear, srcMonth, srcDay := now.Date()
	d, e := date.FromUnits(srcYear, int(srcMonth), srcDay)
	if e != nil {
		t.Errorf("Expected no error, got %v", e)
	} else {

		year, month, day := date.ToUnits(d)
		if year != 2000 || month != 1 || day != 1 {
			t.Errorf("Expected 2000-01-01, got %04v-%02v-%02v", year, month, day)
		}
	}
}

func TestFromUnits(tt *testing.T) {
	cases := []struct {
		name  string
		year  int
		month int
		day   int
		valid bool
	}{
		{
			"units before minimum date",
			1750,
			1,
			1,
			false,
		},
		{
			"units after maximum date",
			10000,
			1,
			1,
			false,
		},
		{
			"minimum date",
			1753,
			1,
			1,
			true,
		},
		{
			"maximum date",
			9999,
			12,
			31,
			true,
		},
	}
	for _, tc := range cases {
		tt.Run(tc.name, func(t *testing.T) {
			dt, e := date.FromUnits(tc.year, tc.month, tc.day)
			if tc.valid {
				if e != nil {
					t.Errorf("Unexpected failure: %v", e)
				} else {
					y, m, d := date.ToUnits(dt)
					if y != tc.year || m != tc.month || d != tc.day {
						t.Errorf("Expected %04d-%02d-%02d, got %04d-%02d-%02d", tc.year, tc.month, tc.day, y, m, d)
					}
				}
			} else {
				if e == nil {
					t.Errorf("Unexpected success")
				}
				y, m, d := date.ToUnits(dt)
				if dt != date.Nil {
					t.Errorf("Expected date.Nil, got %04d-%02d-%02d", y, m, d)
				}
			}
		})
	}
}

func TestGenerationFromCurrentDate(t *testing.T) {

}
