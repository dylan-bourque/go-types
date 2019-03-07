// Copyright 2019 Dylan Bourque. All rights reserved.
//
// Use of this source code is governed by the MIT open source license that can be found in the LICENSE file.

package timeofday

import (
	"bytes"
	"encoding"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/pkg/errors"
)

var (
	// ErrInvalidBinaryDataLen is returned from timeofday.Value.UnmarshalBinary() then the passed-in byte slice
	// is not exactly 8 bytes long
	ErrInvalidBinaryDataLen = errors.Errorf("timeofday.Value: binary data must be 8 bytes")
	// ErrInvalidTextDataLen is returned from timeofday.Value.UnmarshalText() when the passed-in byte slice
	// is not between 8 and 18 bytes long
	ErrInvalidTextDataLen = errors.Errorf("timeofday.Value: text data must be bewteen 8 and 18 bytes")
	// ErrInvalidTextData is returned from timeofday.Value.UnmarshalJSON() when the passed-in byte slice
	// does not contain a string
	ErrInvalidTextData = errors.Errorf("timeofday.Value: can only decode JSON strings")
	// ErrInvalidTimeFormat is returned from timeofday.Value.UnmarshalText() when the passed-in byte slice
	// is not formatted correctly
	ErrInvalidTimeFormat = errors.Errorf("timeofday.Value: text data was not in the correct format")
)

// interface validations
var _ encoding.TextMarshaler = (*Value)(nil)
var _ encoding.TextUnmarshaler = (*Value)(nil)
var _ encoding.BinaryMarshaler = (*Value)(nil)
var _ encoding.BinaryUnmarshaler = (*Value)(nil)
var _ json.Marshaler = (*Value)(nil)
var _ json.Unmarshaler = (*Value)(nil)

// MarshalText implements the encoding.TextMarshaler interface for timeofday.Value values.
//
// The encoded value is the same as is returned by the String() method
func (t Value) MarshalText() ([]byte, error) {
	return []byte(t.String()), nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface for timeofday.Value values.
//
// The supported format is "hh:mm:ss.ffffffff" with the following constraints:
// . "hh" must be 2 decimal digits between 00 and 23, representing the hour of the day
// . "mm" must be 2 decimal digits between 00 and 59, representing the minute of the hour
// . "ss" must be 2 decimal digits between 00 and 59, representing the second of the minute
// . ".fffffffff" is optional, if specified it must be between 1 and 9 decimal digits, respresenting
//   the fractional seconds up to nanosecond-level resolution
func (t *Value) UnmarshalText(text []byte) error {
	if l := len(text); l < 8 || l > 18 {
		return ErrInvalidTextDataLen
	}
	// defer to stdlib to parse the time string in UTC (so no DST)
	tv, err := time.ParseInLocation(`15:04:05.999999999`, string(text), time.UTC)
	if err != nil {
		return ErrInvalidTimeFormat
	}
	// extract the time unit values, construct a timeofday.Value from them and return
	// . no error checking needed in the call to FromUnits() below b/c time.ParseInLocation() would
	//   have failed above if there were invalid unit values
	hh, mm, ss := tv.Clock()
	ns := int64(tv.Nanosecond())
	v, _ := FromUnits(hh, mm, ss, ns)
	t.d = v.d
	return nil
}

// MarshalBinary implements the encoding.BinaryMarshaler interface for timeofday.Value values.
//
// The resulting data is a 64-bit integer in big-endian byte order that contains
// the number of nanoseconds in the underlying time.Duration value.
func (t Value) MarshalBinary() ([]byte, error) {
	var buf bytes.Buffer
	// this can't fail b/c we can always write a 64-bit into into 8 bytes
	_ = binary.Write(&buf, binary.BigEndian, t.d.Nanoseconds())
	return buf.Bytes(), nil
}

// UnmarshalBinary implements the encoding.BinaryUnmarshaler interface for timeofday.Value values.
//
// The provided value must be 8 bytes and contain a 64-bit integer value in big-endian byte order between
// 0 (00:00:00) and 86,399,999,000,000 (23:59:59.999999999).
//
// If data is not 8 bytes, ErrInvalidBinaryDataLen is returned.  If the unmarshalled integer value is
// out of range, ErrInvalidDuration is returned.
func (t *Value) UnmarshalBinary(data []byte) error {
	if len(data) != 8 {
		return ErrInvalidBinaryDataLen
	}
	// this can't fail b/c any 8 bytes can be read into an int64 value
	var d int64
	_ = binary.Read(bytes.NewReader(data), binary.BigEndian, &d)
	// convert to time.Duration and validate range
	dur := time.Duration(d)
	if !IsValidDuration(dur) {
		return ErrInvalidDuration
	}
	// all is well
	t.d = dur
	return nil
}

// MarshalJSON implements the json.Marshaler interface for timeofday.Value values.  The JSON
// encoding is the same as MarshalText().
func (t Value) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("%q", t)), nil
}

// UnmarshalJSON implements the json.Unmarshaler interface for timeofday.Value values.
//
// If the value is the special JSON null token, t is set to timeofday.Zero.  All other values are
// delegated to UnmarshalText().
func (t *Value) UnmarshalJSON(p []byte) error {
	if bytes.Equal(p, []byte("null")) {
		t.d = time.Duration(0)
		return nil
	}
	var s string
	if err := json.NewDecoder(bytes.NewReader(p)).Decode(&s); err != nil {
		return errors.Wrapf(ErrInvalidTextData, "%v", err)
	}
	return t.UnmarshalText([]byte(strings.Trim(s, `"`)))
}
