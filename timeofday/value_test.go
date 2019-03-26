// Copyright 2019 Dylan Bourque. All rights reserved.
//
// Use of this source code is governed by the MIT open source license that can be found in the LICENSE file.

package timeofday

import (
	"fmt"
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

func TestMustPanics(t *testing.T) {
	expected := errors.Errorf("test")
	defer func() {
		got := recover().(error)
		if got != expected {
			t.Errorf("Expected to recover() %v, got %v", expected, got)
		}
	}()
	Must(Value{}, expected)
	t.Error("Expected Must() to panic on error")
}

func TestToDateTimeLocation(t *testing.T) {
	now := time.Now()
	cases := []struct {
		name string
		t    Value
		y    int
		m    time.Month
		d    int
	}{
		{"zero value/beginning of time", Zero, 1, time.January, 1},
		{"zero value/today", Zero, now.Year(), now.Month(), now.Day()},
		{"zero value/end of time", Zero, 9999, time.December, 31},
		{"min value/beginning of time", Min, 1, time.January, 1},
		{"min value/today", Min, now.Year(), now.Month(), now.Day()},
		{"min value/end of time", Min, 9999, time.December, 31},
		{"max value/beginning of time", Max, 1, time.January, 1},
		{"max value/today", Max, now.Year(), now.Month(), now.Day()},
		{"max value/end of time", Max, 9999, time.December, 31},
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
		t        Value
		y        int
		m        time.Month
		d        int
		tz       *time.Location
		expected time.Time
	}{
		{"zero value/beginning of time", Zero, 1, time.January, 1, time.UTC, time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC)},
		{"zero value/today", Zero, now.Year(), now.Month(), now.Day(), now.Location(), time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())},
		{"zero value/end of time", Zero, 9999, time.December, 31, time.UTC, time.Date(9999, time.December, 31, 0, 0, 0, 0, time.UTC)},
		{"min value/beginning of time", Min, 1, time.January, 1, time.UTC, time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC)},
		{"min value/today", Min, now.Year(), now.Month(), now.Day(), now.Location(), time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())},
		{"min value/end of time", Min, 9999, time.December, 31, time.UTC, time.Date(9999, time.December, 31, 0, 0, 0, 0, time.UTC)},
		{"max value/beginning of time", Max, 1, time.January, 1, time.UTC, time.Date(1, time.January, 1, 23, 59, 59, 999999999, time.UTC)},
		{"max value/today", Max, now.Year(), now.Month(), now.Day(), now.Location(), time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 999999999, now.Location())},
		{"max value/end of time", Max, 9999, time.December, 31, time.UTC, time.Date(9999, time.December, 31, 23, 59, 59, 999999999, time.UTC)},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(tt *testing.T) {
			got := tc.t.ToDateTimeInLocation(tc.y, tc.m, tc.d, tc.tz)
			if !got.Equal(tc.expected) {
				tt.Errorf("Expected %s, got %s", tc.expected.Format(time.RFC3339Nano), got.Format(time.RFC3339Nano))
			}
		})
	}
}

func TestAdd(t *testing.T) {
	cases := []struct {
		name     string
		t        Value
		delta    time.Duration
		expected Value
	}{
		{"min value/zero duration", Min, time.Duration(0), Min},
		{"min value/one day duration", Min, 24 * time.Hour, Min},
		{"min value/positive duration", Min, time.Hour, Must(FromUnits(1, 0, 0, 0))},
		{"min value/positive duration/overflow", Min, 25 * time.Hour, Must(FromUnits(1, 0, 0, 0))},
		{"min value/positive duration/multiday overflow", Min, 49 * time.Hour, Must(FromUnits(1, 0, 0, 0))},
		{"min value/negative duration", Min, -1 * time.Hour, Must(FromUnits(23, 0, 0, 0))},
		{"min value/negative duration/overflow", Min, -25 * time.Hour, Must(FromUnits(23, 0, 0, 0))},
		{"min value/negative duration/multiday overflow", Min, -49 * time.Hour, Must(FromUnits(23, 0, 0, 0))},
		{"mid-range value/zero duration", Must(FromUnits(12, 0, 0, 0)), time.Duration(0), Must(FromUnits(12, 0, 0, 0))},
		{"mid-range value/one day duration", Must(FromUnits(12, 0, 0, 0)), 24 * time.Hour, Must(FromUnits(12, 0, 0, 0))},
		{"mid-range value/positive duration", Must(FromUnits(12, 0, 0, 0)), time.Hour, Must(FromUnits(13, 0, 0, 0))},
		{"mid-range value/positive duration/overflow", Must(FromUnits(12, 0, 0, 0)), 13 * time.Hour, Must(FromUnits(1, 0, 0, 0))},
		{"mid-range value/positive duration/multiday overflow", Must(FromUnits(12, 0, 0, 0)), 37 * time.Hour, Must(FromUnits(1, 0, 0, 0))},
		{"mid-range value/negative duration", Must(FromUnits(12, 0, 0, 0)), -1 * time.Hour, Must(FromUnits(11, 0, 0, 0))},
		{"mid-range value/negative duration/overflow", Must(FromUnits(12, 0, 0, 0)), -13 * time.Hour, Must(FromUnits(23, 0, 0, 0))},
		{"mid-range value/negative duration/multiday overflow", Must(FromUnits(12, 0, 0, 0)), -40 * time.Hour, Must(FromUnits(20, 0, 0, 0))},
		{"max value/zero duration", Max, time.Duration(0), Max},
		{"max value/one day duration", Max, 24 * time.Hour, Max},
		{"max value/positive duration", Max, 2 * time.Hour, Must(FromUnits(1, 59, 59, 999999999))},
		{"max value/positive duration/overflow", Max, 25 * time.Hour, Must(FromUnits(0, 59, 59, 999999999))},
		{"max value/positive duration/multiday overflow", Max, 49 * time.Hour, Must(FromUnits(0, 59, 59, 999999999))},
		{"max value/negative duration", Max, -12 * time.Hour, Must(FromUnits(11, 59, 59, 999999999))},
		{"max value/negative duration/overflow", Max, -25 * time.Hour, Must(FromUnits(22, 59, 59, 999999999))},
		{"max value/negative duration/multiday overflow", Max, -49 * time.Hour, Must(FromUnits(22, 59, 59, 999999999))},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(tt *testing.T) {
			got := tc.t.Add(tc.delta)
			if got != tc.expected {
				tt.Errorf("Expected %v, got %v", tc.expected, got)
			}
		})
	}
}

func TestSubtract(t *testing.T) {
	cases := []struct {
		name     string
		t        Value
		delta    time.Duration
		expected Value
	}{
		{"min value/zero duration", Min, time.Duration(0), Min},
		{"min value/one day duration", Min, 24 * time.Hour, Min},
		{"min value/positive duration", Min, time.Hour, Must(FromUnits(23, 0, 0, 0))},
		{"min value/positive duration/overflow", Min, 25 * time.Hour, Must(FromUnits(23, 0, 0, 0))},
		{"min value/positive duration/multiday overflow", Min, 49 * time.Hour, Must(FromUnits(23, 0, 0, 0))},
		{"min value/negative duration", Min, -1 * time.Hour, Must(FromUnits(1, 0, 0, 0))},
		{"min value/negative duration/overflow", Min, -25 * time.Hour, Must(FromUnits(1, 0, 0, 0))},
		{"min value/negative duration/multiday overflow", Min, -49 * time.Hour, Must(FromUnits(1, 0, 0, 0))},
		{"mid-range value/zero duration", Must(FromUnits(12, 0, 0, 0)), time.Duration(0), Must(FromUnits(12, 0, 0, 0))},
		{"mid-range value/one day duration", Must(FromUnits(12, 0, 0, 0)), 24 * time.Hour, Must(FromUnits(12, 0, 0, 0))},
		{"mid-range value/positive duration", Must(FromUnits(12, 0, 0, 0)), time.Hour, Must(FromUnits(11, 0, 0, 0))},
		{"mid-range value/positive duration/overflow", Must(FromUnits(12, 0, 0, 0)), 13 * time.Hour, Must(FromUnits(23, 0, 0, 0))},
		{"mid-range value/positive duration/multiday overflow", Must(FromUnits(12, 0, 0, 0)), 37 * time.Hour, Must(FromUnits(23, 0, 0, 0))},
		{"mid-range value/negative duration", Must(FromUnits(12, 0, 0, 0)), -1 * time.Hour, Must(FromUnits(13, 0, 0, 0))},
		{"mid-range value/negative duration/overflow", Must(FromUnits(12, 0, 0, 0)), -13 * time.Hour, Must(FromUnits(1, 0, 0, 0))},
		{"mid-range value/negative duration/multiday overflow", Must(FromUnits(12, 0, 0, 0)), -40 * time.Hour, Must(FromUnits(4, 0, 0, 0))},
		{"max value/zero duration", Max, time.Duration(0), Max},
		{"max value/one day duration", Max, 24 * time.Hour, Max},
		{"max value/positive duration", Max, 2 * time.Hour, Must(FromUnits(21, 59, 59, 999999999))},
		{"max value/positive duration/overflow", Max, 25 * time.Hour, Must(FromUnits(22, 59, 59, 999999999))},
		{"max value/positive duration/multiday overflow", Max, 49 * time.Hour, Must(FromUnits(22, 59, 59, 999999999))},
		{"max value/negative duration", Max, -12 * time.Hour, Must(FromUnits(11, 59, 59, 999999999))},
		{"max value/negative duration/overflow", Max, -25 * time.Hour, Must(FromUnits(0, 59, 59, 999999999))},
		{"max value/negative duration/multiday overflow", Max, -49 * time.Hour, Must(FromUnits(0, 59, 59, 999999999))},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(tt *testing.T) {
			got := tc.t.Sub(tc.delta)
			if got != tc.expected {
				tt.Errorf("Expected %v, got %v", tc.expected, got)
			}
		})
	}
}

func TestValidation(t *testing.T) {
	type testCase struct {
		name  string
		d     time.Duration
		valid bool
	}
	cases := []testCase{
		{"negative duration", time.Duration(-1), false},
		{"overflow duration", 24 * time.Hour, false},
	}
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for h := 0; h < 24; h++ {
		for m := 0; m < 60; m++ {
			s := rng.Intn(60)
			ns := rng.Int63n(1000000000)
			cases = append(cases, testCase{
				fmt.Sprintf("%02d:%02d:%02d.%09d", h, m, s, ns),
				time.Duration(int64((time.Duration(h)*time.Hour)+(time.Duration(m)*time.Minute)+(time.Duration(s)*time.Second)) + ns),
				true,
			})
		}
	}
	for _, tc := range cases {
		t.Run(tc.name, func(tt *testing.T) {
			got := Value{d: tc.d}
			if got.IsValid() != tc.valid {
				t.Errorf(
					"Expected %v to be %svalid",
					got,
					func() string {
						if tc.valid {
							return ""
						}
						return "in"
					}(),
				)
			}
		})
	}
}

func TestToDuration(t *testing.T) {
	type testCase struct {
		name     string
		t        Value
		expected time.Duration
	}
	cases := []testCase{
		{"invalid value", Value{d: -1 * time.Hour}, time.Duration(-1)},
	}
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for h := 0; h < 24; h++ {
		for m := 0; m < 60; m++ {
			s := rng.Intn(60)
			ns := rng.Int63n(1000000000)
			cases = append(cases, testCase{
				fmt.Sprintf("%02d:%02d:%02d.%09d", h, m, s, ns),
				Must(FromUnits(h, m, s, ns)),
				time.Duration(int64((time.Duration(h)*time.Hour)+(time.Duration(m)*time.Minute)+(time.Duration(s)*time.Second)) + ns),
			})
		}
	}
	for _, tc := range cases {
		t.Run(tc.name, func(tt *testing.T) {
			got := ToDuration(tc.t)
			if got != tc.expected {
				tt.Errorf("Expected %v, got %v", tc.expected, got)
			}
		})
	}
}
