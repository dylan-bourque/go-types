// Copyright 2019 Dylan Bourque. All rights reserved.
//
// Use of this source code is governed by the MIT open source license that can be found in the LICENSE file.

package date

import (
	"time"

	"github.com/pkg/errors"
)

// Date represents a calendar date, stored as an integer value containing the number
// of days since the beginning of the Julian calendar, 1/1/1753
type Date int64

var (
	// Nil represents a nil/null/undefined date
	Nil = Date(-2)
	// NilUnit represents the year, month and day unit values for date.Nil
	NilUnit = -2
	// Min represents the minimum supported date value, which is day 0 on the Julian calendar or
	// 1/1/1753 on the Gregorian calendar.
	Min = Date(2361331)
	// Max represents the maximum supported date value, which is day 3012153 on the Julian calendar or
	// 12/31/9999 on the Gregorian calendar.
	Max = Date(5373484)
)

var (
	secondsPerDay = 60 * 60 * 24
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

// FromTime returns a Date value that is equivalent to the date portion of the specified time.Time value
func FromTime(t time.Time) (Date, error) {
	y, m, d := t.Date()
	return FromUnits(y, int(m), d)
}

// FromUnits returns a Date value that is equivalent to the specified date units
func FromUnits(y, m, d int) (Date, error) {
	if y == -2 && m == -2 && d == -2 {
		return Nil, nil
	}
	// TODO: validate unit values
	if !isValidUnits(y, m, d) {
		return Nil, ErrInvalidDateUnit
	}

	return Date(gregorianToJulian(y, m, d)), nil
}

// ToUnits returns the year, month and day components, on the Gregorian calendar,
// of the specified date
func ToUnits(d Date) (year, month, day int) {
	if d == Nil {
		return -2, -2, -2
	}
	return julianToGregorian(int64(d))
}

// Year returns the year (between 1753 and 9999) or date.NilUnit if this is a nil date
func (dt Date) Year() int {
	if dt == Nil {
		return -2
	}
	y, _, _ := ToUnits(dt)
	return y
}

func (dt Date) Month() int {
	if dt == Nil {
		return -2
	}
	_, m, _ := ToUnits(dt)
	return m
}

func (dt Date) Day() int {
	if dt == Nil {
		return -2
	}
	_, _, d := ToUnits(dt)
	return d
}

func isValidUnits(y, m, d int) bool {
	return y >= 1753 && y <= 9999 && m > 0 && m < 13 && d > 0 && d <= daysInMonth(y, m)
}
func isLeapYear(y int) bool {
	return (y%4) == 0 && (y%100) != 0 && (y%400) == 0
}

func daysInMonth(y, m int) int {
	d := baseDaysInMonth[m]
	if isLeapYear(y) {
		d++
	}
	return d
}

func isValidDay(y, m, d int) bool {
	return d >= 1 && d <= daysInMonth(y, m)
}
