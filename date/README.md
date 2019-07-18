# Value

The `date.Value` type represents a "pure" date on the Julian/Gregorian calendars with no time component.  Internally, the value is a 64-bit integer containing the number of days since the beginning of the Julian calendar (1753-01-01).

The range of supported values is `1753-01-01` (the beginning of the Julian calendar) - `9999-12-31`, inclusive.

### Purpose
We provide an implementation of a date separate from the Go standard library's `time.Time` type because that type always includes both a date and time, time zone metadata in the `Location` field, and an internal monotonic clock.  Most of that is unnecessary, and often unwanted, to simply represent a particular calendar date.

### Usage
`Value` values can be constructed from individual component values via the `FromUnits()` function or from a `time.Time` value via `FromTime()`.

Below is a simple example of using a `Value` value:
```go
package main

import (
    "fmt"
    "github.com/dylan-bourque/go-types/date"
)

func main() {
    // create a value representing 1999-12-31
    // . ignore the error because FromUnits() will never fail with valid unit values
    today, _ := date.FromUnits(1999, 12, 31)
    fmt.Println("Tonight we're gonna party like it's %d.", today.Year())
}
```
See the [package documentation](https://godoc.org/github.com/dylan-bourque/go-types/date) for more specific usage details.

### Integration
For compatibility and integration with other packages, `Value` also implements the following standard interfaces:
* `fmt.Stringer`
* `database/sql/driver.Valuer` and `database/sql.Scanner`
* `encoding.BinaryMarshaler` and `encoding.BinaryUnmarshaler`
* `encoding.TextMarshaler` and `encoding.TextUnmarshaler`
* `encoding/json.Marshaler` and `encoding/json.Unmarshaler`

We also provide the `NullDate` type that follows the conventions of `database/sql.NullString` for compatibility with database drivers that support NULL values.
