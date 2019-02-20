package timeofday

import (
	"bytes"
	"encoding/binary"
	"testing"
	"time"

	"github.com/pkg/errors"
)

func TestMarshalText(t *testing.T) {
	cases := []struct {
		name     string
		v        TimeOfDay
		expected []byte
	}{
		{"zero value", ZeroTime, []byte("00:00:00")},
		{"min value", MinTime, []byte("00:00:00")},
		{"max value", MaxTime, []byte("23:59:59.999999999")},
		{"12:34:56.789012345", Must(FromUnits(12, 34, 56, 789012345)), []byte("12:34:56.789012345")},
		{"12:34:56", Must(FromUnits(12, 34, 56, 0)), []byte("12:34:56")},
		{"12:34:56.1", Must(FromUnits(12, 34, 56, 100000000)), []byte("12:34:56.1")},
		{"12:34:56.11", Must(FromUnits(12, 34, 56, 110000000)), []byte("12:34:56.11")},
		{"12:34:56.111", Must(FromUnits(12, 34, 56, 111000000)), []byte("12:34:56.111")},
		{"12:34:56.1111", Must(FromUnits(12, 34, 56, 111100000)), []byte("12:34:56.1111")},
		{"12:34:56.11111", Must(FromUnits(12, 34, 56, 111110000)), []byte("12:34:56.11111")},
		{"12:34:56.111111", Must(FromUnits(12, 34, 56, 111111000)), []byte("12:34:56.111111")},
		{"12:34:56.1111111", Must(FromUnits(12, 34, 56, 111111100)), []byte("12:34:56.1111111")},
		{"12:34:56.11111111", Must(FromUnits(12, 34, 56, 111111110)), []byte("12:34:56.11111111")},
		{"12:34:56.111111111", Must(FromUnits(12, 34, 56, 111111111)), []byte("12:34:56.111111111")},
		{"12:34:56.011111111", Must(FromUnits(12, 34, 56, 11111111)), []byte("12:34:56.011111111")},
		{"12:34:56.001111111", Must(FromUnits(12, 34, 56, 1111111)), []byte("12:34:56.001111111")},
		{"12:34:56.000111111", Must(FromUnits(12, 34, 56, 111111)), []byte("12:34:56.000111111")},
		{"12:34:56.000011111", Must(FromUnits(12, 34, 56, 11111)), []byte("12:34:56.000011111")},
		{"12:34:56.000001111", Must(FromUnits(12, 34, 56, 1111)), []byte("12:34:56.000001111")},
		{"12:34:56.000000111", Must(FromUnits(12, 34, 56, 111)), []byte("12:34:56.000000111")},
		{"12:34:56.000000011", Must(FromUnits(12, 34, 56, 11)), []byte("12:34:56.000000011")},
		{"12:34:56.000000001", Must(FromUnits(12, 34, 56, 1)), []byte("12:34:56.000000001")},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(tt *testing.T) {
			got, _ := tc.v.MarshalText()
			if !bytes.Equal(got, tc.expected) {
				tt.Errorf("Expected %v, got %v", tc.expected, got)
			}
		})
	}
}

func TestUnmarshalText(t *testing.T) {
	cases := []struct {
		name     string
		d        []byte
		expected TimeOfDay
		err      error
	}{
		// invalid buffer
		{"nil buffer", nil, ZeroTime, ErrInvalidTextDataLen},
		{"empty buffer", []byte{}, ZeroTime, ErrInvalidTextDataLen},
		{"short buffer", []byte{1}, ZeroTime, ErrInvalidTextDataLen},
		{"long buffer", []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19}, ZeroTime, ErrInvalidTextDataLen},
		// malformed text
		{"incorrect format/first separator", []byte("00_00:00"), ZeroTime, ErrInvalidTimeFormat},
		{"incorrect format/second separator", []byte("00:00_00"), ZeroTime, ErrInvalidTimeFormat},
		{"incorrect format/fraction separator", []byte("00:00:00_0"), ZeroTime, ErrInvalidTimeFormat},
		{"invalid value/hours overflow", []byte("24:00:00"), ZeroTime, ErrInvalidTimeFormat},
		{"invalid value/minutes overflow", []byte("00:60:00"), ZeroTime, ErrInvalidTimeFormat},
		{"invalid value/seconds overflow", []byte("00:00:60"), ZeroTime, ErrInvalidTimeFormat},
		// valid text
		{"zero value", []byte("00:00:00"), ZeroTime, nil},
		{"min value", []byte("00:00:00"), MinTime, nil},
		{"max value", []byte("23:59:59.999999999"), MaxTime, nil},
		{"12:34:56.789012345", []byte("12:34:56.789012345"), Must(FromUnits(12, 34, 56, 789012345)), nil},
		{"12:34:56", []byte("12:34:56"), Must(FromUnits(12, 34, 56, 0)), nil},
		{"12:34:56.1", []byte("12:34:56.1"), Must(FromUnits(12, 34, 56, 100000000)), nil},
		{"12:34:56.11", []byte("12:34:56.11"), Must(FromUnits(12, 34, 56, 110000000)), nil},
		{"12:34:56.111", []byte("12:34:56.111"), Must(FromUnits(12, 34, 56, 111000000)), nil},
		{"12:34:56.1111", []byte("12:34:56.1111"), Must(FromUnits(12, 34, 56, 111100000)), nil},
		{"12:34:56.11111", []byte("12:34:56.11111"), Must(FromUnits(12, 34, 56, 111110000)), nil},
		{"12:34:56.111111", []byte("12:34:56.111111"), Must(FromUnits(12, 34, 56, 111111000)), nil},
		{"12:34:56.1111111", []byte("12:34:56.1111111"), Must(FromUnits(12, 34, 56, 111111100)), nil},
		{"12:34:56.11111111", []byte("12:34:56.11111111"), Must(FromUnits(12, 34, 56, 111111110)), nil},
		{"12:34:56.111111111", []byte("12:34:56.111111111"), Must(FromUnits(12, 34, 56, 111111111)), nil},
		{"12:34:56.011111111", []byte("12:34:56.011111111"), Must(FromUnits(12, 34, 56, 11111111)), nil},
		{"12:34:56.001111111", []byte("12:34:56.001111111"), Must(FromUnits(12, 34, 56, 1111111)), nil},
		{"12:34:56.000111111", []byte("12:34:56.000111111"), Must(FromUnits(12, 34, 56, 111111)), nil},
		{"12:34:56.000011111", []byte("12:34:56.000011111"), Must(FromUnits(12, 34, 56, 11111)), nil},
		{"12:34:56.000001111", []byte("12:34:56.000001111"), Must(FromUnits(12, 34, 56, 1111)), nil},
		{"12:34:56.000000111", []byte("12:34:56.000000111"), Must(FromUnits(12, 34, 56, 111)), nil},
		{"12:34:56.000000011", []byte("12:34:56.000000011"), Must(FromUnits(12, 34, 56, 11)), nil},
		{"12:34:56.000000001", []byte("12:34:56.000000001"), Must(FromUnits(12, 34, 56, 1)), nil},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(tt *testing.T) {
			var got TimeOfDay

			err := got.UnmarshalText(tc.d)
			if tc.err != errors.Cause(err) {
				tt.Errorf("Expected error %v, got %v", tc.err, err)
			}
			if got != tc.expected {
				tt.Errorf("Expected %v, got %v", tc.expected, got)
			}
		})
	}
}

func TestMarshalBinary(t *testing.T) {
	cases := []struct {
		name     string
		v        TimeOfDay
		expected []byte
	}{
		{"zero value", ZeroTime, genBinaryDataFromDuration(time.Duration(0))},
		{"min value", MinTime, genBinaryDataFromDuration(time.Duration(0))},
		{"max value", MaxTime, genBinaryDataFromDuration(time.Duration(24*time.Hour - time.Nanosecond))},
		{"12:34:56.789012345", Must(FromUnits(12, 34, 56, 789012345)), genBinaryDataFromDuration(time.Duration(12*time.Hour + 34*time.Minute + 56*time.Second + 789012345))},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(tt *testing.T) {
			d, err := tc.v.MarshalBinary()
			if err != nil {
				tt.Errorf("Unexpected error %v", err)
			}
			if !bytes.Equal(d, tc.expected) {
				tt.Errorf("Expected %v, got %v", tc.expected, d)
			}
		})
	}
}

func TestUnmarshalBinary(t *testing.T) {
	cases := []struct {
		name     string
		d        []byte
		expected TimeOfDay
		err      error
	}{
		{"nil-buffer", nil, ZeroTime, ErrInvalidBinaryDataLen},
		{"empty-buffer", []byte{}, ZeroTime, ErrInvalidBinaryDataLen},
		{"short-buffer", []byte{1}, ZeroTime, ErrInvalidBinaryDataLen},
		{"long-buffer", []byte{1, 2, 3, 4, 5, 6, 7, 8, 9}, ZeroTime, ErrInvalidBinaryDataLen},
		{"invalid-duration-value/negative-underflow", genBinaryDataFromDuration(time.Duration(-1)), ZeroTime, ErrInvalidDuration},
		{"invalid-duration-value/positive-overflow", genBinaryDataFromDuration(24 * time.Hour), ZeroTime, ErrInvalidDuration},
		{"zero-value", genBinaryDataFromDuration(time.Duration(0)), ZeroTime, nil},
		{"min-value", genBinaryDataFromDuration(time.Duration(0)), MinTime, nil},
		{"max-value", genBinaryDataFromDuration(time.Duration(24*time.Hour - time.Nanosecond)), MaxTime, nil},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(tt *testing.T) {
			var got TimeOfDay

			err := got.UnmarshalBinary(tc.d)
			if tc.err != errors.Cause(err) {
				tt.Errorf("Expected error %v, got %v", tc.err, err)
			}
			if got != tc.expected {
				tt.Errorf("Expected %v, got %v", tc.expected, got)
			}
		})
	}
}

// genBinaryDataFromDuration constructs the expected binary encoding for a given clock.TimeOfDay value
// from the provided time.Duration
// . the value is 8 bytes containing a 64-bit integer in big endian byte order, containing the count
//   of nanoseconds
func genBinaryDataFromDuration(dur time.Duration) []byte {
	var buf bytes.Buffer
	_ = binary.Write(&buf, binary.BigEndian, dur.Nanoseconds())
	return buf.Bytes()
}
