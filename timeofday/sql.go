package timeofday

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"

	"github.com/pkg/errors"
)

var (
	// ErrUnsupportedSourceType is returned by .Scan() when the provided value cannot be converted to
	// a timeofday.TimeOfDay value
	ErrUnsupportedSourceType = errors.Errorf("Cannot convert the source data to a timeofday.TimeOfDay value")
)

// Value implements the driver.Valuer interface for TimeOfDay values.  The returned value is the
// default string encoding, hh:mm:ss.fffffffff.
func (t TimeOfDay) Value() (driver.Value, error) {
	return t.String(), nil
}

// Scan implements the sql.Scanner interface for TimeOfDay values.
//
// An 8-byte slice is handled by UnmarshalBinary() and a string is handled by UnmarshalText().  All other
// values will return an error
func (t *TimeOfDay) Scan(src interface{}) error {
	switch tv := src.(type) {
	case []byte:
		return t.UnmarshalBinary(tv)
	case string:
		return t.UnmarshalText([]byte(tv))
	default:
		return errors.Wrapf(ErrUnsupportedSourceType, "Unsupported type: %T", src)
	}
}

// NullTimeOfDay can be used with the standard sql package to represent a TimeOfDay value that can
// be NULL in the database.
type NullTimeOfDay struct {
	TimeOfDay TimeOfDay
	Valid     bool
}

// Value implements the driver.Valuer interface for NullTimeOfDay values
func (t NullTimeOfDay) Value() (driver.Value, error) {
	if !t.Valid {
		return nil, nil
	}
	return t.TimeOfDay.Value()
}

// Scan implements the sql.Scanner interface for NullTimeOfDay values
func (t *NullTimeOfDay) Scan(src interface{}) error {
	if src == nil {
		t.TimeOfDay, t.Valid = Zero, false
		return nil
	}
	if err := t.TimeOfDay.Scan(src); err != nil {
		return err
	}
	t.Valid = true
	return nil
}

// MarshalJSON implements the json.Marshaler interface for NullTimeOfDay values
func (t NullTimeOfDay) MarshalJSON() ([]byte, error) {
	if !t.Valid {
		return json.Marshal(nil)
	}
	return json.Marshal(t.TimeOfDay)
}

// UnmarshalJSON implements the json.Unmarshaler interface for NullTimeOfDay values
func (t *NullTimeOfDay) UnmarshalJSON(d []byte) error {
	if bytes.Equal(d, []byte("null")) {
		t.TimeOfDay, t.Valid = Zero, false
		return nil
	}

	if err := json.Unmarshal(d, &t.TimeOfDay); err != nil {
		return err
	}

	t.Valid = true
	return nil
}
