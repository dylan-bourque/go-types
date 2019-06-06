// Copyright 2019 Dylan Bourque. All rights reserved.
//
// Use of this source code is governed by the MIT open source license that can be found in the LICENSE file.

package date

import (
	"time"

	"github.com/pkg/errors"
)

// Value represents a calendar date, stored as an integer value containing the number
// of days since the beginning of the Julian calendar, 1/1/1753
type Value int64

var (
	// Nil represents a nil/null/undefined date
	Nil = Value(-2)
	// NilUnit represents the year, month and day unit values for date.Nil
	NilUnit = -2
	// Min represents the minimum supported date value, which is day 0 on the Julian calendar or
	// 1753-01-01 on the Gregorian calendar.
	Min = Value(2361331)
	// Max represents the maximum supported date value, which is day 3012153 on the Julian calendar or
	// 9999-12-31 on the Gregorian calendar.
	Max = Value(5373484)
)

var (
	// ErrInvalidDateUnit is returned when an out-of-range date unit value is used
	ErrInvalidDateUnit = errors.Errorf("One or more of the specified date units were invalid")
)

var (
	// the number of days in each month in non-leap years
	// . index 0 is not used so that months values can start at 1
	baseDaysInMonth = [13]int{0, 31, 28, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31}
)

func gregorianToJulian(y, m, d int) (result int64) {
	// adjust to Julian calendar
	if m > 2 {
		m -= 3
	} else {
		m += 9
		y--
	}
	// convert to Julian date
	c := uint64(y / 100)
	yr := (uint64(y) - (100 * c))
	result = int64(((146097 * c) >> 2) + ((1461 * yr) >> 2) + (153*uint64(m)+2)/5 + uint64(d) + 1721119)
	return result
}

func julianToGregorian(v int64) (y, m, d int) {
	// convert to Gregorian date
	jt := uint64(v - 1721119)
	y = int((((jt << 2) - 1) / 146097))
	jt = uint64((jt << 2) - 1 - (146097 * uint64(y)))
	ud := uint64(jt >> 2)
	jt = uint64(((ud << 2) + 3) / 1461)
	ud = uint64((ud << 2) + 3 - (1461 * jt))
	ud = (ud + 4) >> 2
	m = int(((5 * ud) - 3) / 153)
	ud = uint64((5 * ud) - 3 - (153 * uint64(m)))
	d = int(((ud + 5) / 5))
	y = int((uint64(100*y) + jt))
	// adjust to Gregorian calendar
	if m < 10 {
		m += 3
	} else {
		m -= 9
		y++
	}
	return y, m, d
}

// Must panics if the passed-in error is non-nil; otherwise, it returns the passed-in date.Value
func Must(v Value, err error) Value {
	if err != nil {
		panic(err)
	}
	return v
}

// Equal returns true if the specified date.Value values are equal (represent the same date) and false if they do not.
func Equal(v1, v2 Value) bool {
	return int64(v1) == int64(v2)
}

// Equal returns true if the specified date.Value value is equal to the receiver (represent the same date)
// and false if it is not.
func (v Value) Equal(v2 Value) bool {
	return Equal(v, v2)
}

// FromTime returns a Value value that is equivalent to the date portion of the specified time.Time value
func FromTime(t time.Time) (Value, error) {
	y, m, d := t.Date()
	return FromUnits(y, int(m), d)
}

// FromUnits returns a Value value that is equivalent to the specified date units
func FromUnits(y, m, d int) (Value, error) {
	// validate unit values
	if !isValidUnits(y, m, d) {
		return Nil, ErrInvalidDateUnit
	}

	return Value(gregorianToJulian(y, m, d)), nil
}

// ToTime returns a time.Time instance with the year, month and day fields populated from the receiver
// and the time portion set to midnight UTC
func (v Value) ToTime() time.Time {
	if v == Nil {
		return time.Time{}
	}
	y, m, d := ToUnits(v)
	return time.Date(y, time.Month(m), d, 0, 0, 0, 0, time.UTC)
}

// ToUnits returns the year, month and day components, on the Gregorian calendar,
// of the specified date
func ToUnits(d Value) (year, month, day int) {
	if d == Nil {
		return -2, -2, -2
	}
	return julianToGregorian(int64(d))
}

// Year returns the year (between 1753 and 9999) or 0 if this is a nil date
func (dt Value) Year() int {
	if dt == Nil {
		return 0
	}
	y, _, _ := ToUnits(dt)
	return y
}

// Month returns the month (between 1 and 12) or 0 if this is a nil date
func (dt Value) Month() int {
	if dt == Nil {
		return 0
	}
	_, m, _ := ToUnits(dt)
	return m
}

// Day returns the day of the month (between 1 and 31) or 0 if this is a nil date
func (dt Value) Day() int {
	if dt == Nil {
		return 0
	}
	_, _, d := ToUnits(dt)
	return d
}

// IsValid returns true if the specified date.Value is valid (between date.Min and date.Max, inclusive)
// and false if it is not.
func IsValid(d Value) bool {
	return Min <= d && d <= Max
}

func isValidUnits(y, m, d int) bool {
	return y >= 1753 && y <= 9999 && m > 0 && m < 13 && d > 0 && d <= daysInMonth(y, m)
}
func isLeapYear(y int) bool {
	return ((y%4) == 0 && (y%100) != 0) || ((y % 400) == 0)
}

func daysInMonth(y, m int) int {
	d := baseDaysInMonth[m]
	if m == 2 && isLeapYear(y) {
		d++
	}
	return d
}
