# go-types

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://github.com/dylan-bourque/go-types/blob/master/LICENSE)
[![GoDoc](https://godoc.org/github.com/dylan-bourque/go-types?status.svg)](https://godoc.org/github.com/dylan-bourque/go-types)

Package `go-types` provides idomatic Go implementations of useful aggregate types that are not present in the Go standard library.  Each type is provided as a subpackage and, as much as possible, adheres to recommended Go idioms.

### Provided Types
This table contains a summary of the types provided by this package.  See the sub-package level `README.md` files for more details of the individual types.

| Type | Description |
|------|-------------|
| [`timeofday.Value`](timeofday/README.md) | A type that represents a specific time of day - in the range [00:00:00 .. 23:59:59.999999999) - independent of date, time zone and Daylight Savings Time concerns. |
| [`date.Value`](date/README.md) | A type that represents a calendar date with no time component, which is useful for avoiding the vaguaries of what _one day_ means after considering time zones and Daylight Savings Time.|

### Installation

Once you have [installed Go][golang-install], run this command
to install the `go-types` package:

    go get github.com/dylan-bourque/go-types


## License

This source code of this package is released under the MIT License. Please see
the [LICENSE](https://github.com/dylan-bourque/go-types/blob/master/LICENSE) for the full
content of the license.

[golang-install]: http://golang.org/doc/install.html
[sv]: http://semver.org/
