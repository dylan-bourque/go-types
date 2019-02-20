package clock

import (
	"bytes"
	"encoding/binary"
	"testing"
	"time"
)

func TestMarshalBinary(t *testing.T) {
	// genExpected constructs the expected binary encoding for a given clock.Time value
	// from the provided time.Duration
	// . the value is 8 bytes containing a 64-bit integer in big endian byte order,
	//   containing the count of nanoseconds
	genExpected := func(dur time.Duration) []byte {
		var buf bytes.Buffer
		binary.Write(&buf, binary.BigEndian, dur.Nanoseconds())
		return buf.Bytes()
	}
	cases := []struct {
		name     string
		v        Time
		expected []byte
	}{
		{"zero value", Zero, genExpected(time.Duration(0))},
		{"min value", Min, genExpected(time.Duration(0))},
		{"max value", Max, genExpected(time.Duration(24*time.Hour - time.Nanosecond))},
		{"12:34:56.789012345", Must(FromUnits(12, 34, 56, 789012345)), genExpected(time.Duration(12*time.Hour + 34*time.Minute + 56*time.Second + 789012345))},
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
