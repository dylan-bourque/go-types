package timeofday

import (
	"fmt"
	"time"

	"github.com/pkg/errors"
)

// TimeOfDay defines a clock time (hh:mm:ss.fffffffff), independent of any date, time zone, Daylight Savings
// TimeOfDay, etc. considerations.
//
// Internally, the value is stored as a time.Duration value in the range [0ns...24h). The clock time is
// derived by partitioning the total duration into hours, minutes, seconds and nanoseconds.
type TimeOfDay struct {
	d time.Duration
}

var (
	// ZeroTime defines a "zero" clock time, which is equivalent to clock.MinTime
	ZeroTime = TimeOfDay{}
	// MinTime defines the minimum supported clock time, which is midnight (00:00:00)
	MinTime = TimeOfDay{d: 0}
	// MaxTime defines the maximum supported clock time, which is 1 nanosecond before midnight (23:59:59.999999999)
	MaxTime = TimeOfDay{d: time.Duration(24*time.Hour - time.Nanosecond)}
)
var (
	// ErrInvalidUnit indicates that one or more of the specified unit values are out of the allowed range
	ErrInvalidUnit = errors.Errorf("One or more of the specified unit values are outside the valid range")
	// ErrInvalidDuration indicates that a time.Duration value cannot be converted to a TimeOfDay value
	ErrInvalidDuration = errors.Errorf("The specified duration is outside the valid range for a TimeOfDay value")
)

// Must is a helper that wraps a call to a function that returns (clock.TimeOfDay, error)
// and panics if err is non-nil.
func Must(t TimeOfDay, err error) TimeOfDay {
	if err != nil {
		panic(err)
	}
	return t
}

// IsValid returns true if t is a valid clock.TimeOfDay value in the range [00:00:00 .. 24:00:00), false otherwise
func IsValid(t TimeOfDay) bool {
	return IsValidDuration(t.d)
}

// IsValid returns true if t is a valid clock.TimeOfDay value in the range [00:00:00 .. 24:00:00), false otherwise
func (t TimeOfDay) IsValid() bool {
	return IsValid(t)
}

// IsValidUnits returns whether or not the specified unit values are valid for a TimeOfDay value
func IsValidUnits(h, m, s int, ns int64) bool {
	return (0 <= h) && (h < 24) && (0 <= m) && (m < 60) && (0 <= s) && (s < 60) && (0 <= ns) && (ns < 1000000000)
}

const (
	nsecsPerSecond int64 = 1000 * 1000 * 1000
	nsecsPerMinute int64 = 60 * nsecsPerSecond
	nsecsPerHour   int64 = 60 * nsecsPerMinute
)

// ToUnits returns the hour, minute, second and fractional components of a TimeOfDay value
func (t TimeOfDay) ToUnits() (h, m, s int, ns int64) {
	ns = t.d.Nanoseconds()

	uh := ns / nsecsPerHour
	ns -= uh * nsecsPerHour

	um := ns / nsecsPerMinute
	ns -= um * nsecsPerMinute

	us := ns / nsecsPerSecond
	ns -= us * nsecsPerSecond

	return int(uh), int(um), int(us), ns
}

// FromUnits constructs a TimeOfDay value from the provided unit values
//
// If the specified units cannot be converted to a time.Duration or is outside
// of the supported range - [00:00:00 - 24:00:00) - an error is returned
func FromUnits(h, m, s int, ns int64) (TimeOfDay, error) {
	if !IsValidUnits(h, m, s, ns) {
		return ZeroTime, ErrInvalidUnit
	}
	return TimeOfDay{
		d: time.Duration((int64(h) * nsecsPerHour) + (int64(m) * nsecsPerMinute) + (int64(s) * nsecsPerSecond) + ns),
	}, nil
}

// IsValidDuration returns whether or not the specified time.Duration value can be used as a TimeOfDay
func IsValidDuration(d time.Duration) bool {
	return d >= 0 && d < (24*time.Hour)
}

// ToDuration returns a time.Duration value that is equivalent to summing the hours, minutes, seconds,
// and nanoseconds in t, or a duration of -1 nanosecond if t is invalid.
func ToDuration(t TimeOfDay) time.Duration {
	if d := t.d; IsValidDuration(d) {
		return d
	}
	return time.Duration(-1)
}

// FromDuration constructs a TimeOfDay value from the specified duration
//
// If the provided duration is outside of the supported range - [00:00:00 - 24:00:00) - an error is returned.
func FromDuration(d time.Duration) (TimeOfDay, error) {
	if !IsValidDuration(d) {
		return ZeroTime, ErrInvalidDuration
	}
	return TimeOfDay{d: d}, nil
}

// ToDateTimeUTC composes a clock.TimeOfDay value with the specified year, month and day
// in the UTC time zone.
func (t TimeOfDay) ToDateTimeUTC(year int, month time.Month, day int) time.Time {
	return t.ToDateTimeInLocation(year, month, day, time.UTC)
}

// ToDateTimeLocal composes a clock.TimeOfDay value with the specified year, month and day
// in the current local time zone.
func (t TimeOfDay) ToDateTimeLocal(year int, month time.Month, day int) time.Time {
	return t.ToDateTimeInLocation(year, month, day, time.Local)
}

// ToDateTimeInLocation composes the current clock.TimeOfDay value with the specified year, month, day and
// location/time zone to generate a full time.Time value.
func (t TimeOfDay) ToDateTimeInLocation(year int, month time.Month, day int, loc *time.Location) time.Time {
	h, m, s, ns := t.ToUnits()
	return time.Date(year, month, day, h, m, s, int(ns), loc)
}

// String returns a string representation of the TimeOfDay value, formatted as "hh:mm:ss.fffffffff",
// with the fractional portion omitted if it is zero or trailing zeros trimmed otherwise
func (t TimeOfDay) String() string {
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

// Add adds the specified duration to t, normalizing the result to [00:00:00...24:00:00)
func (t TimeOfDay) Add(d time.Duration) TimeOfDay {
	res := time.Duration(t.d + d)
	// for a positive result, step backwards until we're within the supported range
	for ; res > 24*time.Hour; res -= 24 * time.Hour {
	}
	// for a negative result, step forward until we're within the supported range
	for ; res < 0; res += 24 * time.Hour {
	}
	return TimeOfDay{d: res}
}

// Sub adds the specified duration from t, normalizing the result to [00:00:00...24:00:00)
func (t TimeOfDay) Sub(d time.Duration) TimeOfDay {
	return t.Add(-1 * d)
}