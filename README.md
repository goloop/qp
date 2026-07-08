![Go Report Card](https://img.shields.io/badge/go%20report-retired-lightgrey) [![License](https://img.shields.io/badge/license-MIT-brightgreen)](https://github.com/goloop/qp/blob/master/LICENSE) [![Go Reference](https://pkg.go.dev/badge/github.com/goloop/qp/v2.svg)](https://pkg.go.dev/github.com/goloop/qp/v2) [![Stay with Ukraine](https://img.shields.io/static/v1?label=Stay%20with&message=Ukraine%20♥&color=ffD700&labelColor=0057B8&style=flat)](https://u24.gov.ua/)

# qp

`qp` reads URL query parameters into typed Go values — with optional defaulting
and validation, and without panicking on malformed input.

Instead of the manual `r.URL.Query().Get("age")` + `strconv.Atoi` +
bounds-check dance, you read a parameter into the right type in one call. A
`Query` holder decodes the raw query once and serves every read; a bad value
falls back to the default and is *reported* via an error, never silently masked.

## Features

- Read into `int`, `float64`, `string`, `bool`, `time.Duration`, `time.Time`,
  and their slices.
- Defaulting and validation via composable options — `Default`, `Min`, `Max`,
  `Between`, `OneOf`, plus per-element `Each` for slices.
- Both `?ids=1,2,3` and `?ids=1&ids=2&ids=3` (and mixtures) handled the same way.
- Boolean forms: `true/false`, `yes/no`, `on/off`, `1/0` (case-insensitive).
- Custom types via `GetAs`/`PullAs` and a parse function.
- nil-safe: a `nil` URL never panics. Zero dependencies.

## Installation

```sh
go get github.com/goloop/qp/v2
```

```go
import "github.com/goloop/qp/v2"
```

Requires Go 1.24 or newer. The package has no third-party dependencies.

## Quick start

```go
q := qp.New(r.URL) // nil-safe; decodes r.URL.Query() exactly once

page  := q.Int("page", qp.Default(1), qp.Min(1))
limit := q.Int("limit", qp.Default(20), qp.Between(1, 100))
sort  := q.String("sort", qp.Default("created_at"), qp.OneOf("created_at", "name"))
tags  := q.StringSlice("tag")
```

For a single read, top-level shortcuts build a throwaway holder:

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

`Value` is always safe to use. When a present, non-empty value fails to parse or
violates a constraint, it falls back to the default and the error explains why:

```go
r := qp.ParseInt(u, "age", qp.Default(18), qp.Between(18, 65))
// "?age=99": r.Value == 18, r.Error != nil (wraps qp.ErrNotAllowed)
```

## Documentation

- Full reference and recipes: [DOC.md](DOC.md) · [DOC.UK.md](DOC.UK.md)
- Package API: [pkg.go.dev/github.com/goloop/qp/v2](https://pkg.go.dev/github.com/goloop/qp/v2)
- Changes between versions: [CHANGELOG.md](CHANGELOG.md)

## Contributing

Contributions are welcome. Please run `go test ./...`, `go vet ./...` and
`gofmt -l .` before submitting a pull request.

## License

`qp` is released under the MIT License. See [LICENSE](LICENSE).
