# qp — reference

The full reference for the `qp` package: the mental model, the `Query` holder,
the four method families, options, the validation contract, custom types,
slices, errors and practical recipes.

Ukrainian version: **[DOC.UK.md](DOC.UK.md)**.

## Contents

- [Mental model](#mental-model)
- [The Query holder](#the-query-holder)
- [Four method families](#four-method-families)
- [Options](#options)
- [The validation contract](#the-validation-contract)
- [Presence: Contains and Empty](#presence-contains-and-empty)
- [Custom types](#custom-types)
- [Slices](#slices)
- [Errors](#errors)
- [Top-level shortcuts](#top-level-shortcuts)
- [Recipes and tips](#recipes-and-tips)

## Mental model

`qp` reads URL query parameters into typed Go values, with optional defaulting
and validation, and without panicking on malformed input. Instead of the manual
`r.URL.Query().Get("age")` + `strconv.Atoi` + bounds-check dance, you read a
parameter into the right type in one call.

Three ideas define the package:

1. **Read once, serve many.** A `Query` holder decodes the raw query string
   exactly once and answers every read from the decoded map.
2. **Composable options.** Defaults and constraints are passed as functional
   options (`qp.Default`, `qp.Between`, `qp.OneOf`, …) rather than positional
   arguments, so a call reads like a sentence.
3. **Safe by default.** `Value` is always usable; a bad value falls back to the
   default and is *reported* via `Error`, never silently masked. A `nil` URL
   never panics. Zero dependencies.

```go
import "github.com/goloop/qp/v2"
```

## The Query holder

```go
func New(u *url.URL) *Query
func (q *Query) Values() url.Values
func (q *Query) Contains(key string) bool
func (q *Query) Empty(key string) bool
```

A query string is rarely read one parameter at a time, so the central type is a
holder that decodes the raw query once and serves every read from it:

```go
q := qp.New(r.URL) // nil-safe; decodes r.URL.Query() exactly once

page  := q.Int("page", qp.Default(1), qp.Min(1))
limit := q.Int("limit", qp.Default(20), qp.Between(1, 100))
sort  := q.String("sort", qp.Default("created_at"))
tags  := q.StringSlice("tag")
```

## Four method families

Every type (`Int`, `Float`, `String`, `Bool`, `Duration`, `Time`) and its slice
form offers four ways to read a parameter, differing only in the return shape:

| Family | Example | Returns |
|--------|---------|---------|
| Bare   | `q.Int("age", qp.Default(18))` | the value (default on absent/invalid) |
| Get    | `q.GetInt("age")`              | `(value, ok)` — `ok=false` on absent/invalid |
| Pull   | `q.PullInt("age")`             | `*value` — `nil` if absent |
| Parse  | `q.ParseInt("age")`            | full `*Result` with flags and error |

Pick by what the caller needs: the bare form for "just give me a usable value",
`Get` when you branch on validity, `Pull` for optional fields (a `nil` pointer
means "not supplied"), and `Parse` when you need the full story.

`Result[T]` carries everything the parse learned:

```go
type Result[T any] struct {
    Key      string // the parameter name
    Value    T      // the parsed value, or the default if absent/empty/invalid
    Default  T      // the default applied when the value is unusable
    Contains bool   // parameter present in the query
    Empty    bool   // absent, or present with no value
    Error    error  // non-nil if a present, non-empty value was rejected
}
```

## Options

Scalar options constrain a single value:

```go
qp.Default(v)        // value used when absent, empty, or invalid
qp.Min(v)            // open-ended lower bound
qp.Max(v)            // open-ended upper bound
qp.Between(lo, hi)   // inclusive range, bounds auto-swapped
qp.OneOf(a, b, c)    // allowed set (the way to validate strings)
```

Slice options shape a slice read:

```go
qp.Each(opts...)     // apply scalar validation to every element
qp.Sep(";")          // custom separator (default ",")
qp.Fallback(v)       // default slice when absent, empty, or invalid
qp.KeepEmpty()       // keep empty elements instead of dropping them
qp.ErrorOnEmpty()    // reject the whole slice on any empty element
```

`Min`/`Max`/`Between` require an ordered type; `OneOf` requires a comparable
type and is the idiomatic way to validate strings against an allowed set.

## The validation contract

`Value` is always safe to use. When a **present, non-empty** parameter fails to
parse or violates a constraint, `Value` falls back to the default and `Error`
explains why — the failure is reported, never silently masked:

```go
r := qp.ParseInt(u, "age", qp.Default(18), qp.Between(18, 65))
// "?age=99": r.Value == 18, r.Error != nil
// GetInt would return (18, false); PullInt a pointer to 18.
```

An **absent** or **present-but-empty** parameter is not an error — it simply
yields the default, and `Error` stays `nil`. Use the presence helpers to tell
those apart.

## Presence: Contains and Empty

```go
if q.Contains("id") {
    // "id" was supplied (possibly as "?id=")
}
if q.Empty("id") {
    // "id" is absent, or supplied with an empty value
}
```

`Contains` reports whether the key appears at all; `Empty` reports whether it is
absent or present with no value. Together they distinguish "not sent", "sent
empty" and "sent with a value".

## Custom types

For types beyond `int`/`float64`/`string`/`bool`, the `As` family takes a parse
function; `Duration` and `Time` are ready-made:

```go
import "net/netip"

addr, ok := qp.GetAs(q, "ip", netip.ParseAddr)

ttl  := q.Duration("ttl", qp.Default(30*time.Second), qp.Max(time.Hour))
from := q.Time("from", time.RFC3339)
```

`GetAs`/`PullAs` (and `GetSliceAs`/`PullSliceAs`) accept any
`func(string) (T, error)`, so any type with a parser plugs into the same
options, defaulting and validation.

## Slices

Slice parameters accept the comma form, the repeated form, and mixtures of the
two, all normalised the same way: every value is split by the separator, each
element is trimmed, and empty elements are dropped by default. A single invalid
element rejects the whole slice.

```go
u, _ := url.Parse("http://example.com?ids=1,2&ids=3")
ids := qp.IntSlice(u, "ids", qp.Each(qp.Between(1, 100))) // [1 2 3]
```

`KeepEmpty` retains empty elements; `ErrorOnEmpty` turns any empty element into
a rejection; `Sep` changes the separator; `Each` applies scalar options to every
element; `Fallback` supplies the slice used when the parameter is unusable.

## Errors

Errors wrap typed sentinels plus a `*ParseError` carrying the key and offending
value:

```go
var (
    ErrInvalid      // value could not be parsed as the target type
    ErrNotAllowed   // value violated a constraint (range / OneOf)
    ErrEmptyElement // a slice element was empty under ErrorOnEmpty
)

type ParseError struct {
    Key   string // the parameter name
    Value string // the offending raw value (empty for ErrEmptyElement)
    Err   error  // the wrapped sentinel
}
```

Match them with `errors.Is`, and pull the detail with `errors.As`:

```go
if errors.Is(r.Error, qp.ErrNotAllowed) {
    var pe *qp.ParseError
    errors.As(r.Error, &pe) // pe.Key, pe.Value
}
```

## Top-level shortcuts

For a single read there are package-level functions that build a throwaway
holder from a `*url.URL`, mirroring the holder methods:

```go
age, ok := qp.GetInt(r.URL, "age", qp.Between(18, 65))
role    := qp.String(r.URL, "role", qp.Default("guest"), qp.OneOf("admin", "user"))
ids     := qp.IntSlice(r.URL, "ids", qp.Each(qp.Between(1, 100)))
```

Prefer `qp.New(u)` when you read several parameters from the same URL, so the
query is decoded once.

## Recipes and tips

**Decode once per request.** Build one `qp.New(r.URL)` at the top of a handler
and read every parameter from it, rather than calling top-level functions
repeatedly.

**Optional SQL filters.** `Pull` methods make optional filters concise — a
`nil` pointer means "not supplied", so only present parameters enter the `WHERE`
clause:

```go
q := qp.New(r.URL)
active := q.PullBool("is_active")   // *bool or nil
if active != nil { /* add "is_active = ?" */ }
```

**Validate strings with `OneOf`.** There is no range for strings; constrain them
to an allowed set and pair it with a `Default` for the fallback.

**Report, don't hide.** When you care *why* a value was rejected, use the
`Parse` family and inspect `Result.Error`; the bare/`Get`/`Pull` forms already
apply the default safely.

**Custom separators and types.** Reach for `Sep` when clients send
`a;b;c`, and `GetAs` when the value is a `netip.Addr`, a UUID-like token or any
type with a `func(string) (T, error)` parser.
