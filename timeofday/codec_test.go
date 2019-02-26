package timeofday

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
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
		{"zero value", Zero, []byte("00:00:00")},
		{"min value", Min, []byte("00:00:00")},
		{"max value", Max, []byte("23:59:59.999999999")},
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
		{"nil buffer", nil, Zero, ErrInvalidTextDataLen},
		{"empty buffer", []byte{}, Zero, ErrInvalidTextDataLen},
		{"short buffer", []byte{1}, Zero, ErrInvalidTextDataLen},
		{"long buffer", []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19}, Zero, ErrInvalidTextDataLen},
		// malformed text
		{"incorrect format/first separator", []byte("00_00:00"), Zero, ErrInvalidTimeFormat},
		{"incorrect format/second separator", []byte("00:00_00"), Zero, ErrInvalidTimeFormat},
		{"incorrect format/fraction separator", []byte("00:00:00_0"), Zero, ErrInvalidTimeFormat},
		{"invalid value/hours overflow", []byte("24:00:00"), Zero, ErrInvalidTimeFormat},
		{"invalid value/minutes overflow", []byte("00:60:00"), Zero, ErrInvalidTimeFormat},
		{"invalid value/seconds overflow", []byte("00:00:60"), Zero, ErrInvalidTimeFormat},
		// valid text
		{"zero value", []byte("00:00:00"), Zero, nil},
		{"min value", []byte("00:00:00"), Min, nil},
		{"max value", []byte("23:59:59.999999999"), Max, nil},
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

func TestMarshalJSON(t *testing.T) {
	cases := []struct {
		name     string
		v        TimeOfDay
		expected []byte
	}{
		{"zero value", Zero, []byte(`"00:00:00"`)},
		{"min value", Min, []byte(`"00:00:00"`)},
		{"max value", Max, []byte(`"23:59:59.999999999"`)},
		{"12:34:56.789012345", Must(FromUnits(12, 34, 56, 789012345)), []byte(`"12:34:56.789012345"`)},
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

func TestUnmarshalJSON(t *testing.T) {
	cases := []struct {
		name     string
		d        []byte
		expected TimeOfDay
		err      error
	}{
		{"00:00:00", []byte(`"00:00:00"`), Zero, nil},
		{"23:59:59.999999999", []byte(`"23:59:59.999999999"`), Max, nil},
		{"12:34:56.789012345", []byte(`"12:34:56.789012345"`), Must(FromUnits(12, 34, 56, 789012345)), nil},
		{"24:00:00", []byte(`"24:00:00"`), Zero, ErrInvalidTimeFormat},
		{"garbage input", []byte(`"nafklsd8234as"`), Zero, ErrInvalidTimeFormat},
		{"empty string", []byte(`""`), Zero, ErrInvalidTextDataLen},
		{"short input", []byte(`"12"`), Zero, ErrInvalidTextDataLen},
		{"long input", []byte(`"1234567890123456789"`), Zero, ErrInvalidTextDataLen},
		{"JSON 'null'", []byte(`null`), Zero, nil},
		{"JSON number", []byte("42"), Zero, ErrInvalidTextData},
		{"JSON array", []byte("[]"), Zero, ErrInvalidTextData},
		{"JSON object", []byte("{}"), Zero, ErrInvalidTextData},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(tt *testing.T) {
			var got TimeOfDay
			err := json.Unmarshal(tc.d, &got)
			if errors.Cause(err) != tc.err {
				tt.Errorf("Expected error %v, got %v", tc.err, err)
			}
			if got != tc.expected {
				tt.Errorf("Expected %s, got %s", tc.expected, got)
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
		{"zero value", Zero, genBinaryDataFromDuration(time.Duration(0))},
		{"min value", Min, genBinaryDataFromDuration(time.Duration(0))},
		{"max value", Max, genBinaryDataFromDuration(time.Duration(24*time.Hour - time.Nanosecond))},
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
		{"nil-buffer", nil, Zero, ErrInvalidBinaryDataLen},
		{"empty-buffer", []byte{}, Zero, ErrInvalidBinaryDataLen},
		{"short-buffer", []byte{1}, Zero, ErrInvalidBinaryDataLen},
		{"long-buffer", []byte{1, 2, 3, 4, 5, 6, 7, 8, 9}, Zero, ErrInvalidBinaryDataLen},
		{"invalid-duration-value/negative-underflow", genBinaryDataFromDuration(time.Duration(-1)), Zero, ErrInvalidDuration},
		{"invalid-duration-value/positive-overflow", genBinaryDataFromDuration(24 * time.Hour), Zero, ErrInvalidDuration},
		{"zero-value", genBinaryDataFromDuration(time.Duration(0)), Zero, nil},
		{"min-value", genBinaryDataFromDuration(time.Duration(0)), Min, nil},
		{"max-value", genBinaryDataFromDuration(time.Duration(24*time.Hour - time.Nanosecond)), Max, nil},
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
