# Value

The `timeofday.Value` type represents a "clock time" identified by hour, minute, second and fraction components.  Internally, the value is a `time.Duration` made up of the sum of the four component values with the following formula:

    value = (hh * 3.6e12) + (mm * 6e10) + (ss * 1e9) + fff)

The supported range of values is `00:00:00` (midnight) - `23:59:59.999999999` (one nanosecond before midnight).  For simplicity, `Value` only supports 24-hour clocks.

### Purpose
We provide an implementation of a "clock time" separate from the Go standard library's `time.Time` type because that type always includes both a date and time, time zone metadata in the `Location` field, and an internal monotonic clock.  Most of that is unnecessary, and often unwanted, to simply represent a time of "08:30 am".

Using `Value` also avoids any complications for times that occur twice, or not at all, on Daylight Savings Time boundaries.  A `Value` value of `01:45:00` always exists at exactly one point in the day.  It does not occur twice when transitioning from STD to DST and it is not skipped when moving from DST to STD since there is no concept of Standard Time versus Daylight Savings Time.

### Usage
`Value` values can be constructed from individual component values via the `FromUnits()` function or from a `time.Duration` value via `FromDuration()`.

Below is a simple example of using a `Value` value:
```go
package main

import (
    "fmt"
    "github.com/dylan-bourque/go-types/timeofday"
)

func main() {
    // create a value representing 1:30 am
    // . ignore the error because FromUnits() will never fail with valid unit values
    tod, _ := timeofday.FromUnits(1, 30, 0, 0)
    fmt.Println("The clock shows:", tod)
}
```
See the [package documentation](https://godoc.org/github.com/dylan-bourque/go-types/timeofday) for more specific usage details.

### Integration
For compatibility and integration with other packages, `Value` also implements the following standard interfaces:
* `fmt.Stringer`
* `database/sql/driver.Valuer` and `database/sql.Scanner`
* `encoding.BinaryMarshaler` and `encoding.BinaryUnmarshaler`
* `encoding.TextMarshaler` and `encoding.TextUnmarshaler`
* `encoding/json.Marshaler` and `encoding/json.Unmarshaler`

We also provide the `NullTimeOfDay` type that follows the conventions of `database/sql.NullString` for compatibility with database drivers that support NULL values.
