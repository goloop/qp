# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [2.0.0]

Version 2 is a redesign. The module path is now `github.com/goloop/qp/v2` and
requires Go 1.24. The API changed; v1 code does not compile unchanged.

### Added
- `Query` holder: `qp.New(u)` decodes the query string once and serves every
  read from that snapshot. `qp.From(url.Values)` reuses an already-decoded map.
- Functional options for defaulting and validation: `Default`, `Min`, `Max`,
  `Between`, `OneOf` (scalars) and `Each`, `Sep`, `Fallback`, `KeepEmpty`,
  `ErrorOnEmpty` (slices).
- Per-element slice validation via `Each` (range/set checks on every element).
- Bare reader family (`Int`, `String`, `IntSlice`, …) returning the value
  directly, alongside the `Get`, `Pull`, and `Parse` families.
- Custom slice separators via `Sep`.
- Runnable examples and fuzz tests.

### Changed
- Validation arguments moved from positional variadics (`opt ...T`) to explicit
  options; the hidden "second argument means range" rule is gone.
- Out-of-range / disallowed values now have one clear contract: `Value` falls
  back to the default and `Error` is set; `Get` reports `ok=false`.
- Slice parsing trims each element and drops empty elements by default; the
  comma form, the repeated form, and mixtures of the two are normalised
  identically.
- `Empty` now distinguishes an absent parameter from a present-but-empty one.
- Boolean parsing reports a clear error for unrecognised input.

### Fixed
- Panic on a `nil` `*url.URL`.
- Corrupted error message for string validation (wrong format verb).
- Slice parsing failing on mixed comma/repeated forms and on surrounding
  whitespace.
- Slices could not validate their elements.

### Removed
- The external dependency: the module is now dependency-free.
- The `Value` type constraint and the `Min`/`Max`/`Others` fields on `Result`
  (constraints now live in the options, not the result).
