[![Go Report Card](https://goreportcard.com/badge/github.com/goloop/qp)](https://goreportcard.com/report/github.com/goloop/qp) [![License](https://img.shields.io/badge/license-MIT-brightgreen)](https://github.com/goloop/qp/blob/master/LICENSE) [![License](https://img.shields.io/badge/godoc-YES-green)](https://godoc.org/github.com/goloop/qp) [![Stay with Ukraine](https://img.shields.io/static/v1?label=Stay%20with&message=Ukraine%20♥&color=ffD700&labelColor=0057B8&style=flat)](https://u24.gov.ua/)


# qp

The **qp** package reads URL query parameters into typed Go values — with
optional defaulting and validation, and without panicking on malformed input.

Instead of the manual `r.URL.Query().Get("age")` + `strconv.Atoi` +
bounds-check dance, you read a parameter into the right type in one call.


## Installation

```sh
go get github.com/goloop/qp/v2
```

```go
import "github.com/goloop/qp/v2"
```

## Features

- Read into `int`, `float64`, `string`, `bool`, and their slices
- Defaulting and validation via composable options (range, allowed set)
- Per-element validation for slices
- Both `?ids=1,2,3` and `?ids=1&ids=2&ids=3` (and mixtures) handled the same way
- Boolean forms: `true/false`, `yes/no`, `on/off`, `1/0` (case-insensitive)
- Decodes the query string once per request via a holder
- nil-safe: a `nil` URL never panics
- Zero dependencies

## The Query holder

A query string is rarely read one parameter at a time, so the central type is a
holder that decodes the raw query once and serves every read from it:

```go
q := qp.New(r.URL) // nil-safe; decodes r.URL.Query() exactly once

page  := q.Int("page", qp.Default(1), qp.Min(1))
limit := q.Int("limit", qp.Default(20), qp.Between(1, 100))
sort  := q.String("sort", qp.Default("created_at"))
tags  := q.StringSlice("tag")
```

For a single read there are top-level shortcuts that build a throwaway holder:

```go
age, ok := qp.GetInt(r.URL, "age", qp.Between(18, 65))
```

## Four method families

Every type (and its slice form) offers four ways to read a parameter, differing
only in the return shape:

| Family | Example | Returns |
|--------|---------|---------|
| Bare   | `q.Int("age", qp.Default(18))` | the value (default on absent/invalid) |
| Get    | `q.GetInt("age")`              | `(value, ok)` |
| Pull   | `q.PullInt("age")`             | `*value` (nil if absent) |
| Parse  | `q.ParseInt("age")`            | full `*Result` with flags and error |

## Options

Scalar options:

```go
qp.Default(v)        // value used when absent, empty, or invalid
qp.Min(v)            // open-ended lower bound
qp.Max(v)            // open-ended upper bound
qp.Between(lo, hi)   // inclusive range, bounds auto-swapped
qp.OneOf(a, b, c)    // allowed set (the way to validate strings)
```

Slice options:

```go
qp.Each(opts...)     // apply scalar validation to every element
qp.Sep(";")          // custom separator (default ",")
qp.Fallback(v)       // default slice when absent, empty, or invalid
qp.KeepEmpty()       // keep empty elements instead of dropping them
qp.ErrorOnEmpty()    // reject the slice on any empty element
```

## Validation contract

`Value` is always safe to use. When a present, non-empty parameter fails to
parse or violates a constraint, `Value` falls back to the default and `Error`
explains why — the failure is reported, never silently masked:

```go
r := qp.ParseInt(u, "age", qp.Default(18), qp.Between(18, 65))
// "?age=99": r.Value == 18, r.Error != nil
// GetInt would return (18, false); PullInt a pointer to 18.
```

Errors wrap sentinels (`ErrInvalid`, `ErrNotAllowed`, `ErrEmptyElement`) and a
`*ParseError` with the key and offending value:

```go
if errors.Is(r.Error, qp.ErrNotAllowed) {
    var pe *qp.ParseError
    errors.As(r.Error, &pe) // pe.Key, pe.Value
}
```

## Custom types

For types beyond `int`/`float64`/`string`/`bool`, the `As` family takes a parse
function; `Duration` and `Time` are ready-made:

```go
import "net/netip"

addr, ok := qp.GetAs(q, "ip", netip.ParseAddr)

ttl  := q.Duration("ttl", qp.Default(30*time.Second), qp.Max(time.Hour))
from := q.Time("from", time.RFC3339)
```

## Slices

Slice parameters accept the comma form, the repeated form, and mixtures of the
two, all normalised the same way: every value is split by the separator, each
element is trimmed, and empty elements are dropped by default. A single invalid
element rejects the whole slice.

```go
u, _ := url.Parse("http://example.com?ids=1,2&ids=3")
ids := qp.IntSlice(u, "ids", qp.Each(qp.Between(1, 100))) // [1 2 3]
```

## Examples

### Reading typed values

```go
u, _ := url.Parse("http://example.com?age=25&temp=36.6&name=alice&active=true")
q := qp.New(u)

age    := q.Int("age", qp.Between(0, 120))
temp   := q.Float("temp")
name   := q.String("name", qp.Default("guest"))
active := q.Bool("active")
```

### Validating against an allowed set

```go
role, ok := q.GetString("role", qp.Default("guest"), qp.OneOf("admin", "user"))
// "?role=hacker" -> role == "guest", ok == false
```

### Optional SQL filters

`Pull` methods make optional filters concise — a nil pointer means "not
supplied", so only the parameters actually present enter the `WHERE` clause:

```go
func where(isActive, isStaff, isSuperuser *bool) string {
    check := [...]struct {
        name  string
        value *bool
    }{
        {name: "is_active", value: isActive},
        {name: "is_staff", value: isStaff},
        {name: "is_superuser", value: isSuperuser},
    }

    and := make([]string, 0, len(check))
    for _, m := range check {
        if m.value != nil {
            and = append(and, fmt.Sprintf("%s=%t", m.name, *m.value))
        }
    }

    if len(and) != 0 {
        return "WHERE " + strings.Join(and, ", ")
    }
    return ""
}

u, _ := url.Parse("http://example.com?is_active=true&is_staff=false")
q := qp.New(u)
query := where(
    q.PullBool("is_active"),
    q.PullBool("is_staff"),
    q.PullBool("is_superuser"), // nil
)
// Result: WHERE is_active=true, is_staff=false
```

### Presence checks

```go
if q.Contains("id") {
    // "id" was supplied (possibly as "?id=")
}
if q.Empty("id") {
    // "id" was supplied with an empty value
}
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
