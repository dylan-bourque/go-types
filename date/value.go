// Copyright 2019 Dylan Bourque. All rights reserved.
//
// Use of this source code is governed by the MIT open source license that can be found in the LICENSE file.

package date

import (
	"fmt"
	"time"

	"github.com/pkg/errors"
)

// Value represents a calendar date, stored as an integer value representing the number
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

// Must panics if the passed-in error is non-nil; otherwise, it returns the passed-in date.Value
func Must(v Value, err error) Value {
	if err != nil {
		panic(err)
	}
	return v
}

// FromTime returns a Value value that is equivalent to the date portion of the specified time.Time value
func FromTime(t time.Time) (Value, error) {
	y, m, d := t.Date()
	return FromUnits(y, int(m), d)
}

// FromUnits returns a Value value that is equivalent to the specified date units
func FromUnits(y, m, d int) (Value, error) {
	// validate unit values
	if !IsValidUnits(y, m, d) {
		return Nil, ErrInvalidDateUnit
	}

	return Value(gregorianToJulian(y, m, d)), nil
}

// ToTime returns a time.Time instance with the year, month and day fields populated from the receiver
// and the time portion set to midnight UTC
func (v Value) ToTime() time.Time {
	if !v.IsValid() {
		return time.Time{}
	}
	y, m, d := ToUnits(v)
	return time.Date(y, time.Month(m), d, 0, 0, 0, 0, time.UTC)
}

// ToUnits returns the year, month and day components, on the Gregorian calendar,
// of the specified date
func ToUnits(d Value) (year, month, day int) {
	if d == Nil {
		return NilUnit, NilUnit, NilUnit
	}
	if !d.IsValid() {
		return -1, -1, -1
	}
	return julianToGregorian(int64(d))
}

// Year returns the year (between 1753 and 9999) or 0 if this is a nil date
func (dt Value) Year() int {
	if dt == Nil {
		return NilUnit
	}
	y, _, _ := ToUnits(dt)
	return y
}

// Month returns the month (between 1 and 12) or 0 if this is a nil date
func (dt Value) Month() int {
	if dt == Nil {
		return NilUnit
	}
	_, m, _ := ToUnits(dt)
	return m
}

// Day returns the day of the month (between 1 and 31) or 0 if this is a nil date
func (dt Value) Day() int {
	if dt == Nil {
		return NilUnit
	}
	_, _, d := ToUnits(dt)
	return d
}

// IsValid returns true if the date.Value is valid (between date.Min and date.Max, inclusive)
// and false if it is not.
func (d Value) IsValid() bool {
	return Min <= d && d <= Max
}

// Equal returns true if the specified date.Value values are equal (represent the same date) and false if they do not.
//
// *NOTE*
// The Nil value is treated specially and is not less than, equal to, or greater than any value, so this
// function returns false if either value is Nil.
func Equal(v1, v2 Value) bool {
	if v1 == Nil || v2 == Nil {
		return false
	}
	return int64(v1) == int64(v2)
}

// Equal returns true if the specified date.Value value is equal to the receiver (represent the same date)
// and false if it is not.
//
// *NOTE*
// The Nil value is treated specially and is not less than, equal to, or greater than any value, so this
// method returns false if the reciever or the specified value are Nil.
func (v Value) Equals(v2 Value) bool {
	return Equal(v, v2)
}

// Before returns true if the specified date.Value is after the receiver and false if it does not.
//
// *NOTE*
// The Nil value is treated specially and is not less than, equal to, or greater than any value, so this
// method returns false if the reciever or the specified value are Nil.
func (v Value) Before(v2 Value) bool {
	if v == Nil || v2 == Nil {
		return false
	}
	return int64(v) < int64(v2)
}

// After returns true if the specified date.Value is before the receiver and false if it does not.
//
// *NOTE*
// The Nil value is treated specially and is not less than, equal to, or greater than any value, so this
// method returns false if the reciever or the specified value are Nil.
func (v Value) After(v2 Value) bool {
	if v == Nil || v2 == Nil {
		return false
	}
	return int64(v) > int64(v2)
}

// String implements fmt.Stringer for date.Value instances.
//
// The returns string is formatted as "YYYY-MM-DD".
func (v Value) String() string {
	y, m, d := ToUnits(v)
	return fmt.Sprintf("%04d-%02d-%02d", y, m, d)
}

// Format returns a textual representation of the date value according to the same rules as
// time.Time.Format(), with the restriction that the time portion of the result will always be
// midnight UTC.
func (v Value) Format(layout string) string {
	return v.ToTime().Format(layout)
}

// Parse parses a formatted string and returns the date value that it represents according to the
// same rules as time.Parse().
func Parse(layout, value string) (Value, error) {
	t, err := time.Parse(layout, value)
	if err != nil {
		return Nil, err
	}
	return FromTime(t)
}

// DaysInMonth returns the number of days in the specified month, accounting for leap years.
//
// If the specified year or month or outside the range of valid values, NilUnit is returned.
func DaysInMonth(y, m int) int {
	if !IsValidYear(y) || !IsValidMonth(m) {
		return NilUnit
	}
	d := baseDaysInMonth[m]
	if m == 2 && IsLeapYear(y) {
		d++
	}
	return d
}

// DaysInYear returns the number of days in the specified year, accounting for leap years.
//
// If the specified year is outside the range of valid values, NilUnit is returned.
func DaysInYear(y int) int {
	if !IsValidYear(y) {
		return NilUnit
	}
	n := 365
	if IsLeapYear(y) {
		n += 1
	}
	return n
}

// StartOfYear returns a new date.Value that represents the first day of the
// year for the date represented by d.
//
// If the receiver is date.Nil, this method return date.Nil
func (d Value) StartOfYear() Value {
	if !d.IsValid() {
		return Nil
	}
	y, _, _ := ToUnits(d)
	v, _ := FromUnits(y, 1, 1)
	return v
}

// MiddleOfYear returns a new date.Value that represents the middle day of the
// year - defined as June 30 - for the date represented by d.
//
// If the receiver is date.Nil, this method return date.Nil
func (d Value) MiddleOfYear() Value {
	if !d.IsValid() {
		return Nil
	}
	y, _, _ := ToUnits(d)
	v, _ := FromUnits(y, 6, 30)
	return v
}

// EndOfYear returns a new date.Value that represents the last day of the
// year for the date represented by d.
//
// If the receiver is date.Nil, this method return date.Nil
func (d Value) EndOfYear() Value {
	if !d.IsValid() {
		return Nil
	}
	y, _, _ := ToUnits(d)
	v, _ := FromUnits(y, 12, 31)
	return v
}

// StartOfMonth returns a new date.Value that represents the first day of the
// month for the date represented by d.
//
// If the receiver is date.Nil, this method return date.Nil
func (d Value) StartOfMonth() Value {
	if !d.IsValid() {
		return Nil
	}
	y, m, _ := ToUnits(d)
	v, _ := FromUnits(y, m, 1)
	return v
}

// MiddleOfMonth returns a new date.Value that represents the "middle" day -
// defined as <days in month>/2 - of the month for the date represented by d.
//
// If the receiver is date.Nil, this method return date.Nil
func (d Value) MiddleOfMonth() Value {
	if !d.IsValid() {
		return Nil
	}
	y, m, _ := ToUnits(d)
	v, _ := FromUnits(y, m, DaysInMonth(y, m)/2)
	return v
}

// EndOfMonth returns a new date.Value that represents the last day of the
// month for the date represented by d.
//
// If the receiver is date.Nil, this method return date.Nil
func (d Value) EndOfMonth() Value {
	if !d.IsValid() {
		return Nil
	}
	y, m, _ := ToUnits(d)
	v, _ := FromUnits(y, m, DaysInMonth(y, m))
	return v
}

// NextMonth returns a new date.Value that represents the same day on a subsequent
// month.
//
// If the receiver is date.Nil, this method returns date.Nil
func (d Value) NextMonth(m int) (Value, error) {
	if !d.IsValid() {
		return Nil, nil
	}
	yr, mon, day := ToUnits(d)
	if m <= mon {
		yr++
	}
	return FromUnits(yr, m, day)
}

// NextWeekday returns a new date.Value that represents a subsequent week day relative
// to the current date.
//
// If the receiver is date.Nil, this method returns date.Nil
func (d Value) NextWeekday(wd time.Weekday) (Value, error) {
	if !d.IsValid() {
		return Nil, nil
	}
	delta := int(wd - d.Weekday())
	if delta <= 0 {
		delta += 7
	}
	return d.AddDays(delta)
}

// NextYear returns a new date.Value that represents the same month and day on a subsequent
// year relative to the current date.
//
// If the receiver is date.Nil, this method returns date.Nil
func (d Value) NextYear(yy int) (Value, error) {
	if !d.IsValid() {
		return Nil, nil
	}
	if !IsValidYear(yy) {
		return Nil, errors.Errorf("invalid year unit value: %d", yy)
	}
	yr, mon, day := ToUnits(d)
	if yr > yy {
		return Nil, errors.Errorf("the specified year, %d, is before the current year", yy)
	}
	return FromUnits(yy, mon, day)
}

// PreviousWeekday returns a new date.Value that represents a prior weekday relative to the current
// date.
//
// If the receiver is date.Nil, this method returns date.Nil
func (d Value) PreviousWeekday(w time.Weekday) (Value, error) {
	if !d.IsValid() {
		return Nil, nil
	}
	delta := int(w - d.ToTime().Weekday())
	if delta >= 0 {
		delta -= 7
	}
	return d.AddDays(delta)
}

// Weekday returns the date of the week represented by the current date.
//
// If the receiver is date.Nil, this method returns -1
func (d Value) Weekday() time.Weekday {
	if !d.IsValid() {
		return -1
	}
	return d.ToTime().Weekday()
}

// AddDays adds the specified number of days to the current date.
//
// If the receiver is date.Nil, this method returns date.Nil and no error
func (d Value) AddDays(n int) (Value, error) {
	if !d.IsValid() {
		return Nil, nil
	}
	v := int64(d) + int64(n)
	if v < int64(Min) || v > int64(Max) {
		return Nil, errors.Errorf("adding %d days would generate in an out-of-range result", n)
	}
	return Value(v), nil
}

// Add adds the specified duration to the current date.
//
// Because date.Value has no concept of time, any "partial" day information will be
// discarded in the result.
func (d Value) Add(dur time.Duration) Value {
	v, err := FromTime(d.ToTime().Add(dur))
	if err != nil {
		return Nil
	}
	return v
}

// IsValidUnits returns a value indicating whether or not the specified combination of
// date unit values represent a valid date.
func IsValidUnits(y, m, d int) bool {
	return IsValidYear(y) && IsValidMonth(m) && d > 0 && d <= DaysInMonth(y, m)
}

// IsValidYear returns a value indicating whether or not the specified year falls within
// the range of supported values: 1753 to 9999, inclusive.
func IsValidYear(y int) bool {
	return y >= 1753 && y <= 9999
}

// IsValidMonth returns a value indicating whether or not the specified month falls within
// the range of supported values: 1 to 12, inclusive.
func IsValidMonth(m int) bool {
	return m > 0 && m < 13
}

// IsLeapYear returns a value indicating whether or not the specified year is a leap year.
//
// A leap year is defined as being divisible by 4, but not divisible by 100 unless it is
// also divisible by 400.
//
// This function always returns false for an invalid year value.
func IsLeapYear(y int) bool {
	return IsValidYear(y) && (((y%4) == 0 && (y%100) != 0) || ((y % 400) == 0))
}

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
