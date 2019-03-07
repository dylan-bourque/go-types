// Copyright 2019 Dylan Bourque. All rights reserved.
//
// Use of this source code is governed by the MIT open source license that can be found in the LICENSE file.

package timeofday

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"
	"math/rand"
	"testing"
	"time"

	"github.com/pkg/errors"
)

func TestValuer(t *testing.T) {
	type testCase struct {
		name     string
		t        Value
		err      error
		expected string
	}
	var cases []testCase
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for h := 0; h < 24; h++ {
		for m := 0; m < 60; m++ {
			s := rng.Intn(60)
			ns := rng.Int63n(1000000000)

			v := Must(FromUnits(h, m, s, ns))
			cases = append(cases, testCase{
				name:     v.String(),
				t:        v,
				err:      nil,
				expected: v.String(),
			})
		}
	}

	for _, tc := range cases {
		t.Run(tc.name, func(tt *testing.T) {
			got, err := tc.t.Value()
			if err != tc.err {
				tt.Errorf("Expected error %v, got %v", tc.err, err)
			}
			if got != tc.expected {
				tt.Errorf("Expected %s, got %s", tc.expected, got)
			}
		})
	}
}

func TestScanner(t *testing.T) {
	type testCase struct {
		name     string
		d        interface{}
		expected Value
		err      error
	}
	cases := []testCase{
		{"nil input", nil, Zero, ErrUnsupportedSourceType},
		{"invalid input type", 42, Zero, ErrUnsupportedSourceType},
		{"invalid byte slice", []byte{42, 43}, Zero, ErrInvalidBinaryDataLen},
		{"valid byte slice", genBinaryDataFromDuration(8 * time.Hour), Must(FromUnits(8, 0, 0, 0)), nil},
		{"short text input", "blah", Zero, ErrInvalidTextDataLen},
		{"invalid text input", "24:00:00", Zero, ErrInvalidTimeFormat},
		{"valid text input", "12:34:56.789012345", Must(FromUnits(12, 34, 56, 789012345)), nil},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(tt *testing.T) {
			var got Value
			err := got.Scan(tc.d)
			if errors.Cause(err) != tc.err {
				tt.Errorf("Expected error %v, got %v", tc.err, err)
			}
			if got != tc.expected {
				tt.Errorf("Expected %s, got %s", tc.expected, got)
			}
		})
	}
}

func TestNullTimeOfDayValuer(t *testing.T) {
	type testCase struct {
		name     string
		v        NullTimeOfDay
		expected driver.Value
		err      error
	}
	cases := []testCase{
		{"null value", NullTimeOfDay{}, nil, nil},
	}
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	h, m, s := rng.Intn(24), rng.Intn(60), rng.Intn(60)
	ns := rng.Int63n(1000000000)

	v := Must(FromUnits(h, m, s, ns))
	cases = append(cases, testCase{
		name:     v.String(),
		v:        NullTimeOfDay{TimeOfDay: v, Valid: true},
		expected: v.String(),
		err:      nil,
	})

	for _, tc := range cases {
		t.Run(tc.name, func(tt *testing.T) {
			got, err := tc.v.Value()
			if err != tc.err {
				tt.Errorf("Expected error %v, got %v", tc.err, err)
			}
			if got != tc.expected {
				tt.Errorf("Expected %s, got %s", tc.expected, got)
			}
		})
	}
}

func TestNullTimeOfDayScanner(t *testing.T) {
	type testCase struct {
		name     string
		d        interface{}
		expected NullTimeOfDay
		err      error
	}
	cases := []testCase{
		{"nil input", nil, NullTimeOfDay{}, nil},
		{"invalid input type", 42, NullTimeOfDay{TimeOfDay: Zero}, ErrUnsupportedSourceType},
		{"invalid byte slice", []byte{42, 43}, NullTimeOfDay{TimeOfDay: Zero}, ErrInvalidBinaryDataLen},
		{"valid byte slice", genBinaryDataFromDuration(8 * time.Hour), NullTimeOfDay{TimeOfDay: Must(FromUnits(8, 0, 0, 0)), Valid: true}, nil},
		{"short text input", "blah", NullTimeOfDay{TimeOfDay: Zero}, ErrInvalidTextDataLen},
		{"invalid text input", "24:00:00", NullTimeOfDay{TimeOfDay: Zero}, ErrInvalidTimeFormat},
		{"valid text input", "12:34:56.789012345", NullTimeOfDay{TimeOfDay: Must(FromUnits(12, 34, 56, 789012345)), Valid: true}, nil},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(tt *testing.T) {
			var got NullTimeOfDay
			err := got.Scan(tc.d)
			if errors.Cause(err) != tc.err {
				tt.Errorf("Expected error %v, got %v", tc.err, err)
			}
			if got != tc.expected {
				tt.Errorf("Expected %v, got %v", tc.expected, got)
			}
		})
	}
}

func TestNullTimeOfDayMarshalJSON(t *testing.T) {
	cases := []struct {
		name     string
		v        NullTimeOfDay
		expected []byte
	}{
		{"zero value", NullTimeOfDay{}, []byte("null")},
		{"min value", NullTimeOfDay{TimeOfDay: Min, Valid: true}, []byte(`"00:00:00"`)},
		{"max value", NullTimeOfDay{TimeOfDay: Max, Valid: true}, []byte(`"23:59:59.999999999"`)},
		{"12:34:56.789012345", NullTimeOfDay{TimeOfDay: Must(FromUnits(12, 34, 56, 789012345)), Valid: true}, []byte(`"12:34:56.789012345"`)},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(tt *testing.T) {
			got, err := json.Marshal(tc.v)
			if err != nil {
				tt.Errorf("Unexpected error %v", err)
			}
			if !bytes.Equal(got, tc.expected) {
				tt.Errorf("Expected %s, got %s", string(tc.expected), string(got))
			}
		})
	}
}

func TestNullTimeOfDayUnmarshalJSON(t *testing.T) {
	cases := []struct {
		name     string
		d        []byte
		expected NullTimeOfDay
		err      error
	}{
		{"00:00:00", []byte(`"00:00:00"`), NullTimeOfDay{TimeOfDay: Zero, Valid: true}, nil},
		{"23:59:59.999999999", []byte(`"23:59:59.999999999"`), NullTimeOfDay{TimeOfDay: Max, Valid: true}, nil},
		{"12:34:56.789012345", []byte(`"12:34:56.789012345"`), NullTimeOfDay{TimeOfDay: Must(FromUnits(12, 34, 56, 789012345)), Valid: true}, nil},
		{"24:00:00", []byte(`"24:00:00"`), NullTimeOfDay{}, ErrInvalidTimeFormat},
		{"garbage input", []byte(`"nafklsd8234as"`), NullTimeOfDay{}, ErrInvalidTimeFormat},
		{"empty string", []byte(`""`), NullTimeOfDay{}, ErrInvalidTextDataLen},
		{"short input", []byte(`"12"`), NullTimeOfDay{}, ErrInvalidTextDataLen},
		{"long input", []byte(`"1234567890123456789"`), NullTimeOfDay{}, ErrInvalidTextDataLen},
		{"JSON 'null'", []byte(`null`), NullTimeOfDay{}, nil},
		{"JSON number", []byte("42"), NullTimeOfDay{}, ErrInvalidTextData},
		{"JSON array", []byte("[]"), NullTimeOfDay{}, ErrInvalidTextData},
		{"JSON object", []byte("{}"), NullTimeOfDay{}, ErrInvalidTextData},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(tt *testing.T) {
			var got NullTimeOfDay
			err := json.Unmarshal(tc.d, &got)
			if errors.Cause(err) != tc.err {
				tt.Errorf("Expected error %v, got %v", tc.err, err)
			}
			if got != tc.expected {
				tt.Errorf("Expected %v, got %v", tc.expected, got)
			}
		})
	}
}
