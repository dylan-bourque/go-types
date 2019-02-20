package timeofday

import (
	"math/rand"
	"testing"
	"time"

	"github.com/pkg/errors"
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
			if got != ZeroTime || err != ErrInvalidUnit {
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
			if got != ZeroTime {
				t.Errorf("Expected %v, got %v", ZeroTime, got)
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

func TestMustPanics(t *testing.T) {
	expected := errors.Errorf("test")
	defer func() {
		got := recover().(error)
		if got != expected {
			t.Errorf("Expected to recover() %v, got %v", expected, got)
		}
	}()
	Must(TimeOfDay{}, expected)
	t.Error("Expected Must() to panic on error")
}

func TestToDateTimeLocation(t *testing.T) {
	now := time.Now()
	cases := []struct {
		name string
		t    TimeOfDay
		y    int
		m    time.Month
		d    int
	}{
		{"zero value/beginning of time", ZeroTime, 1, time.January, 1},
		{"zero value/today", ZeroTime, now.Year(), now.Month(), now.Day()},
		{"zero value/end of time", ZeroTime, 9999, time.December, 31},
		{"min value/beginning of time", MinTime, 1, time.January, 1},
		{"min value/today", MinTime, now.Year(), now.Month(), now.Day()},
		{"min value/end of time", MinTime, 9999, time.December, 31},
		{"max value/beginning of time", MaxTime, 1, time.January, 1},
		{"max value/today", MaxTime, now.Year(), now.Month(), now.Day()},
		{"max value/end of time", MaxTime, 9999, time.December, 31},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(tt *testing.T) {
			tt.Run("local", func(ttt *testing.T) {
				h, m, s, ns := tc.t.ToUnits()
				expected := time.Date(tc.y, tc.m, tc.d, h, m, s, int(ns), time.Local)
				got := tc.t.ToDateTimeLocal(tc.y, tc.m, tc.d)
				if !got.Equal(expected) {
					t.Errorf("Expected %s, got %s", expected.Format(time.RFC3339Nano), got.Format(time.RFC3339Nano))
				}
			})
			tt.Run("utc", func(ttt *testing.T) {
				h, m, s, ns := tc.t.ToUnits()
				expected := time.Date(tc.y, tc.m, tc.d, h, m, s, int(ns), time.UTC)
				got := tc.t.ToDateTimeUTC(tc.y, tc.m, tc.d)
				if !got.Equal(expected) {
					t.Errorf("Expected %s, got %s", expected.Format(time.RFC3339Nano), got.Format(time.RFC3339Nano))
				}
			})
		})
	}
}

func TestToDateTime(t *testing.T) {
	now := time.Now()
	cases := []struct {
		name     string
		t        TimeOfDay
		y        int
		m        time.Month
		d        int
		tz       *time.Location
		expected time.Time
	}{
		{"zero value/beginning of time", ZeroTime, 1, time.January, 1, time.UTC, time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC)},
		{"zero value/today", ZeroTime, now.Year(), now.Month(), now.Day(), now.Location(), time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())},
		{"zero value/end of time", ZeroTime, 9999, time.December, 31, time.UTC, time.Date(9999, time.December, 31, 0, 0, 0, 0, time.UTC)},
		{"min value/beginning of time", MinTime, 1, time.January, 1, time.UTC, time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC)},
		{"min value/today", MinTime, now.Year(), now.Month(), now.Day(), now.Location(), time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())},
		{"min value/end of time", MinTime, 9999, time.December, 31, time.UTC, time.Date(9999, time.December, 31, 0, 0, 0, 0, time.UTC)},
		{"max value/beginning of time", MaxTime, 1, time.January, 1, time.UTC, time.Date(1, time.January, 1, 23, 59, 59, 999999999, time.UTC)},
		{"max value/today", MaxTime, now.Year(), now.Month(), now.Day(), now.Location(), time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 999999999, now.Location())},
		{"max value/end of time", MaxTime, 9999, time.December, 31, time.UTC, time.Date(9999, time.December, 31, 23, 59, 59, 999999999, time.UTC)},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(tt *testing.T) {
			got := tc.t.ToDateTimeInLocation(tc.y, tc.m, tc.d, tc.tz)
			if !got.Equal(tc.expected) {
				t.Errorf("Expected %s, got %s", tc.expected.Format(time.RFC3339Nano), got.Format(time.RFC3339Nano))
			}
		})
	}
}
