package clock

import (
	"bytes"
	"encoding/binary"
	"time"

	"github.com/pkg/errors"
)

// MarshalBinary implements the encoding.BinaryMarshaler interface for clock.Time values.
//
// The resulting data is a 64-bit integer in big-endian byte order that contains
// the number of nanoseconds in the underlying time.Duration value.
func (t Time) MarshalBinary() ([]byte, error) {
	var buf bytes.Buffer
	if err := binary.Write(&buf, binary.BigEndian, t.d.Nanoseconds()); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// UnmarshalBinary implements the encoding.BinaryUnmarshaler interface for clock.Time values.
func (t *Time) UnmarshalBinary(data []byte) error {
	if len(data) != 8 {
		return errors.Errorf("clock.Time: value must be 8 bytes, got %d", len(data))
	}
	var d int64
	if err := binary.Read(bytes.NewReader(data), binary.BigEndian, &d); err != nil {
		return err
	}
	dur := time.Duration(d)
	if !IsValidDuration(dur) {
		return errors.Errorf("clock.Time: the provided value was outside the supported range")
	}
	t.d = dur
	return nil
}
