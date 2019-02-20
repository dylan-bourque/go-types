package clock

import (
	"fmt"
	"time"

	"github.com/pkg/errors"
)

// Time defines a clock time, independent of any date, time zone, Daylight Savings Time, etc. considerations
//
// Internally, the clock time is stored as a time.Duration value in the range [0ns .. 23h59m59s999999999ns).
// The clock time is derived by partitioning the total duration into hours, minutes, seconds and nanoseconds.
type Time struct {
	d time.Duration
}

var (
	// Zero defines a "zero" clock time, which represents midnight on the clock
	Zero = Time{d: 0}
	// Min defines the minimum supported clock time, which is midnight (00:00:00)
	Min = Time{d: 0}
	// Max defines the maximum supported clock time, which is 1 nanosecond before midnight (23:59:59.999999999)
	Max = Time{d: time.Duration(24*time.Hour - time.Nanosecond)}
)
var (
	// ErrInvalidUnit indicates that one or more of the specified unit values are out of the allowed range
	ErrInvalidUnit = errors.Errorf("One or more of the specified unit values are outside the valid range")
	// ErrInvalidDuration indicates that a time.Duration value cannot be converted to a Time value
	ErrInvalidDuration = errors.Errorf("The specified duration is outside the valid range for a Time value")
)

// Must is a helper that wraps a call to a function that returns (clock.Time, error)
// and panics if err is non-nil.
func Must(t Time, err error) Time {
	if err != nil {
		panic(err)
	}
	return t
}

// FromUnits constructs a Time value from the provided unit values
//
// If the specified units cannot be converted to a time.Duration or is outside
// of the supported range - [00:00:00 - 24:00:00) - an error is returned
func FromUnits(h, m, s int, ns int64) (Time, error) {
	if !IsValidUnits(h, m, s, ns) {
		return Zero, ErrInvalidUnit
	}
	return Time{
		d: time.Duration((int64(h) * nsecsPerHour) + (int64(m) * nsecsPerMinute) + (int64(s) * nsecsPerSecond) + ns),
	}, nil
}

// IsValidUnits returns whether or not the specified unit values are valid for a Time value
func IsValidUnits(h, m, s int, ns int64) bool {
	return (0 <= h) && (h < 24) && (0 <= m) && (m < 60) && (0 <= s) && (s < 60) && (0 <= ns) && (ns < 1000000000)
}

// FromDuration constructs a Time value from the specified duration
//
// If the provided duration is outside of the supported range - [00:00:00 - 24:00:00) - an error is returned.
func FromDuration(d time.Duration) (Time, error) {
	if !IsValidDuration(d) {
		return Zero, ErrInvalidDuration
	}
	return Time{d: d}, nil
}

// IsValidDuration returns whether or not the specified time.Duration value can be used as a Time
func IsValidDuration(d time.Duration) bool {
	return d >= 0 && d < (24*time.Hour)
}

// String returns a string representation of the Time value, formatted as "hh:mm:ss.fffffffff",
// with the fractional portion omitted if it is zero or trailing zeros trimmed otherwise
func (t Time) String() string {
	h, m, s, ns := t.ToUnits()
	result := fmt.Sprintf("%02d:%02d:%02d", h, m, s)
	if ns > 0 {
		result += fmtFrac(uint64(ns))
	}
	return result
}

// fmtFrac formats the fraction of v/10**9 (e.g., ".12345") into a string, omitting trailing zeros.
// It omits the decimal point too if the fraction is 0.
//
// NOTE: shamelessly "borrowed" from the Go source code for formatting the fractional portion of
// time.Duration values
func fmtFrac(v uint64) string {
	// v is always in the range [0..10^9], so we need a max. of 10 characters
	buf := make([]byte, 10)
	w, print := len(buf), false
	for i := 0; i < 9; i++ {
		digit := v % 10
		print = print || digit != 0
		if print {
			w--
			buf[w] = byte(digit) + '0'
		}
		v /= 10
	}
	if print {
		w--
		buf[w] = '.'
	}
	return string(buf[w:])
}

// ToUnits returns the hour, minute, second and fractional components of a Time value
func (t Time) ToUnits() (h, m, s int, ns int64) {
	return toUnits(t)
}

// ToStandardTimeUTC composes a clock.Time value with the specified year, month and day
// in the UTC time zone.
func (t Time) ToStandardTimeUTC(year int, month time.Month, day int) time.Time {
	return t.ToStandardTimeInLocation(year, month, day, time.UTC)
}

// ToStandardTimeLocal composes a clock.Time value with the specified year, month and day
// in the current local time zone.
func (t Time) ToStandardTimeLocal(year int, month time.Month, day int) time.Time {
	return t.ToStandardTimeInLocation(year, month, day, time.Local)
}

// ToStandardTimeInLocation composes the current clock.Time value with the specified year, month, day and
// location/time zone to generate a full time.Time value.
func (t Time) ToStandardTimeInLocation(year int, month time.Month, day int, loc *time.Location) time.Time {
	h, m, s, ns := toUnits(t)
	return time.Date(year, month, day, h, m, s, int(ns), loc)
}

var (
	nsecsPerSecond = time.Second.Nanoseconds()
	nsecsPerMinute = time.Minute.Nanoseconds()
	nsecsPerHour   = time.Hour.Nanoseconds()
)

func toUnits(v Time) (h, m, s int, ns int64) {
	ns = v.d.Nanoseconds()

	uh := ns / nsecsPerHour
	ns -= uh * nsecsPerHour

	um := ns / nsecsPerMinute
	ns -= um * nsecsPerMinute

	us := ns / nsecsPerSecond
	ns -= us * nsecsPerSecond

	return int(uh), int(um), int(us), ns
}
