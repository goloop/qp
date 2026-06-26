// Package qp parses query parameters from URLs into typed Go values.
//
// Instead of the manual r.URL.Query().Get("age") + strconv.Atoi + bounds-check
// dance, qp reads a parameter into the right type in one call, with optional
// defaulting and validation, and without panicking on malformed input.
//
// # The Query holder
//
// A query string is rarely read one parameter at a time, so the central type is
// a holder that decodes the raw query once and serves every read from it:
//
//	q := qp.New(r.URL) // nil-safe; decodes r.URL.Query() exactly once
//
//	page  := q.Int("page", qp.Default(1), qp.Min(1))
//	limit := q.Int("limit", qp.Default(20), qp.Between(1, 100))
//	sort  := q.String("sort", qp.Default("created_at"))
//	tags  := q.StringSlice("tag")
//
// For a single read there are top-level shortcuts that build a throwaway holder:
//
//	age, ok := qp.GetInt(r.URL, "age", qp.Between(18, 65))
//
// # Four method families
//
// Every type (int, float64, string, bool) and its slice form offers four ways
// to read a parameter, differing only in the return shape:
//
//   - Bare (Int, String, IntSlice, …): returns just the value, falling back to
//     the default when the parameter is absent, empty, or invalid. The common
//     case.
//   - Get (GetInt, GetString, …): returns (value, ok); ok is true only when the
//     parameter was present, non-empty, and valid.
//   - Pull (PullInt, PullString, …): returns a pointer that is nil when the
//     parameter is absent, otherwise points to the value. Ideal for optional
//     filters where "unset" must be distinguished from a real value.
//   - Parse (ParseInt, ParseString, …): returns the full Result with the value,
//     presence flags, the applied default, and any validation error.
//
// # Options
//
// Reads are tuned with functional options. Scalar options:
//
//	qp.Default(v)        // value used when absent, empty, or invalid
//	qp.Min(v) qp.Max(v)  // open-ended bounds (ordered types)
//	qp.Between(lo, hi)   // inclusive range, bounds auto-swapped
//	qp.OneOf(a, b, c)    // allowed set (the way to validate strings)
//
// Slice options:
//
//	qp.Each(opts...)     // apply scalar validation to every element
//	qp.Sep(";")          // custom separator (default ",")
//	qp.Fallback(v)       // default slice when absent, empty, or invalid
//	qp.KeepEmpty()       // keep empty elements instead of dropping them
//	qp.ErrorOnEmpty()    // reject the slice on any empty element
//
// # Validation contract
//
// Value is always safe to use. When a present, non-empty parameter fails to
// parse or violates a constraint, Value falls back to the default and Error
// explains why — the failure is reported, never silently masked:
//
//	r := qp.ParseInt(u, "age", qp.Default(18), qp.Between(18, 65))
//	// "?age=99": r.Value == 18, r.Error != nil
//	// GetInt would return (18, false); PullInt a pointer to 18.
//
// # Slices
//
// Slice parameters accept both the comma form ("?ids=1,2,3") and the repeated
// form ("?ids=1&ids=2&ids=3"), and mixtures of the two, all normalised the same
// way: every value is split by the separator, each element is trimmed, and
// empty elements are dropped by default. A single unparseable or disallowed
// element rejects the whole slice.
//
//	ids := qp.IntSlice(u, "ids", qp.Each(qp.Between(1, 100)))
//
// # Boolean forms
//
// Booleans accept true/false, yes/no, on/off and 1/0, case-insensitively.
//
// # Example: optional SQL filters
//
// Pull methods make optional filters concise — a nil pointer means "not
// supplied", so only the parameters actually present enter the WHERE clause:
//
//	q := qp.New(r.URL)
//	isActive := q.PullBool("is_active")
//	isStaff := q.PullBool("is_staff")
//	isSuperuser := q.PullBool("is_superuser") // nil when absent
package qp
