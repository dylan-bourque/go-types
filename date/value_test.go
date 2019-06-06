// Copyright 2019 Dylan Bourque. All rights reserved.
//
// Use of this source code is governed by the MIT open source license that can be found in the LICENSE file.

package date

import (
	"fmt"
	"testing"
	"time"
)

func TestFromUnits(tt *testing.T) {
	cases := []struct {
		name  string
		year  int
		month int
		day   int
		valid bool
	}{
		{"units before minimum date", 1750, 1, 1, false},
		{"units after maximum date", 10000, 1, 1, false},
		{"minimum date", 1753, 1, 1, true},
		{"maximum date", 9999, 12, 31, true},
		{"invalid month/underflow", 2000, 0, 21, false},
		{"invalid month/overflow", 2000, 13, 21, false},
		{"invalid day/underflow", 2000, 1, 0, false},
		{"invalid day/overflow", 2000, 1, 32, false},
		{"leap year/02-29", 1996, 2, 29, true},
		{"non-leap year/02-29", 1997, 2, 29, false},
	}
	for m := 1; m < 13; m++ {
		for d := 1; d < 31; d++ {
			if m == 2 && d > 28 {
				continue
			}
			if (m == 9 || m == 4 || m == 6 || m == 11) && d > 30 {
				continue
			}
			c := struct {
				name  string
				year  int
				month int
				day   int
				valid bool
			}{
				fmt.Sprintf("valid units/%02d-%02d", m, d),
				2000,
				m,
				d,
				true,
			}
			cases = append(cases, c)
		}
	}
	for _, tc := range cases {
		tt.Run(tc.name, func(t *testing.T) {
			dt, e := FromUnits(tc.year, tc.month, tc.day)
			if tc.valid {
				if e != nil {
					t.Errorf("Unexpected failure: %v", e)
				} else {
					y, m, d := ToUnits(dt)
					if y != tc.year || m != tc.month || d != tc.day {
						t.Errorf("Expected %04d-%02d-%02d, got %04d-%02d-%02d", tc.year, tc.month, tc.day, y, m, d)
					}
				}
			} else {
				if e == nil {
					t.Errorf("Unexpected success")
				}
				y, m, d := ToUnits(dt)
				if dt != Nil {
					t.Errorf("Expected date.Nil, got %04d-%02d-%02d", y, m, d)
				}
			}
		})
	}
}

func TestFromTime(t *testing.T) {
	var (
		now     = time.Now()
		y, m, d = now.Date()
	)

	cases := []struct {
		name     string
		time     time.Time
		valid    bool
		expected Value
	}{
		{"invalid input/zero time", time.Time{}, false, Nil},
		{"invalid input/Min - 1 day", time.Date(1752, time.December, 31, 0, 0, 0, 0, time.UTC), false, Nil},
		{"valid input/Min value", time.Date(1753, time.January, 1, 0, 0, 0, 0, time.UTC), true, Min},
		{"valid input/now", now, true, Must(FromUnits(y, int(m), d))},
		{"valid input/Max value", time.Date(9999, time.December, 31, 0, 0, 0, 0, time.UTC), true, Max},
	}
	for _, c := range cases {
		t.Run(c.name, func(tt *testing.T) {
			got, err := FromTime(c.time)
			if c.valid {
				if err != nil {
					tt.Errorf("Unexpected error: %v", err)
				}
				if !got.Equal(c.expected) {
					tt.Errorf("Unexpected result: expected %v, got %v", c.expected, got)
				}
			} else {
				if err == nil {
					tt.Errorf("Expected error, got <nil>")
				}
			}
		})
	}
}
