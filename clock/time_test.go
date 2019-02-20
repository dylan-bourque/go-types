package clock

import (
	"math/rand"
	"testing"
	"time"
)

func TestConstructTimeFromValidUnits(t *testing.T) {
	for h := 0; h < 24; h++ {
		for m := 0; m < 60; m++ {
			for s := 0; s < 60; s++ {
				ns := rand.Int63n(10 ^ 9)
				got, err := FromUnits(h, m, s, ns)
				if err != nil {
					t.Errorf("Unexpected error %v for valid unit values (%02d:%02d:%02d.%d)", err, h, m, s, ns)
				}
				hh, mm, ss, nn := got.ToUnits()
				if hh != h || mm != m || ss != s || nn != ns {
					t.Errorf("Expected %02d:%02d:%02d.%d, got %02d:%02d:%02d.%d", h, m, s, ns, hh, mm, ss, nn)
				}
			}
		}
	}
}

func TestConstructTimeFromInvalidUnits(t *testing.T) {
	cases := []struct {
		name    string
		h, m, s int
		ns      int64
	}{
		{"hours underflow", -1, 0, 0, 0},
		{"hours overflow", 24, 0, 0, 0},
		{"minutes underflow", 0, -1, 0, 0},
		{"minutes overflow", 0, 60, 0, 0},
		{"seconds underflow", 0, 0, -1, 0},
		{"seconds overflow", 0, 0, 60, 0},
		{"nanoseconds underflow", 0, 0, 0, -1},
		{"nanoseconds overflow", 0, 0, 0, 1000000000},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(tt *testing.T) {
			got, err := FromUnits(tc.h, tc.m, tc.s, tc.ns)
			if got != Zero || err != ErrInvalidUnit {
				t.Errorf("%02d:%02d:%02d.%d - Expected error, got (%s, <nil>)", tc.h, tc.m, tc.s, tc.ns, got.d)
			}
		})
	}
}

func TestConstructTimeFromValidDuration(t *testing.T) {
	for d := int64(0); d < 24*int64(time.Hour); d += int64(time.Second) {
		d += rand.Int63n(10 ^ 9)
		dur := time.Duration(d)
		got, err := FromDuration(dur)
		if err != nil {
			t.Errorf("Unexpected error %v for valid duration %s", err, dur)
		}
		hh, mm, ss, nn := got.ToUnits()
		dd := time.Duration(int64(hh)*int64(time.Hour) + int64(mm)*int64(time.Minute) + int64(ss)*int64(time.Second) + nn)
		if dd != dur {
			t.Errorf("Expected %s, got %s", dur, dd)
		}
	}
}

func TestConstructTimeFromInvalidDuration(t *testing.T) {
	cases := []struct {
		name string
		dur  time.Duration
	}{
		{"negative underflow", time.Duration(-1)},
		{"positive overflow", time.Duration(24*time.Hour + time.Nanosecond)},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(tt *testing.T) {
			got, err := FromDuration(tc.dur)
			if err != ErrInvalidDuration {
				t.Errorf("Expected %v, got %v", ErrInvalidDuration, err)
			}
			if got != Zero {
				t.Errorf("Expected %v, got %v", Zero, got)
			}
		})
	}
}

func TestValidateDuration(t *testing.T) {
	cases := []struct {
		name  string
		dur   time.Duration
		valid bool
	}{
		{"negative underflow", time.Duration(-1), false},
		{"minimum value", time.Duration(0), true},
		{"midrange value", 12 * time.Hour, true},
		{"maximum value", 24*time.Hour - time.Nanosecond, true},
		{"positive overflow", 24 * time.Hour, false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(tt *testing.T) {
			got := IsValidDuration(tc.dur)
			if got != tc.valid {
				tt.Errorf("Expected %t, got %t", tc.valid, got)
			}
		})
	}
}
