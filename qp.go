package qp

import (
	"net/url"
	"slices"
)

// Result holds the outcome of parsing a single query parameter.
//
// Value is always safe to use: on a missing, empty, or invalid parameter it
// equals the configured default (or the zero value of T when no default is
// set). Error is non-nil only when the parameter was present and non-empty but
// could not be parsed or failed validation — in that case Value still holds the
// default, while Error explains why the input was rejected.
//
// The combination of flags distinguishes every state:
//
//	Contains  Empty  Error    meaning
//	false     true   nil      parameter absent
//	true      true   nil      parameter present but empty ("?k=")
//	true      false  nil      parameter present and valid
//	true      false  non-nil  parameter present but invalid (out of range,
//	                          not in set, unparseable)
type Result[T any] struct {
	Key   string // the query parameter name
	Value T      // the parsed value, or the default if absent/empty/invalid

	Default T // the default value applied when the value is unusable

	Contains bool  // true if the parameter is present in the query
	Empty    bool  // true if the parameter is absent or present with no value
	Error    error // non-nil if a present, non-empty value was rejected
}

// Query is an immutable view over the parameters of a URL query string. It
// decodes the raw query exactly once at construction, so reading many
// parameters from the same request costs a single parse instead of one parse
// per call.
//
// A Query is safe for concurrent reads. The zero value is not usable; create
// one with New or From.
type Query struct {
	values url.Values
}

// New builds a Query from the query string of u. It is nil-safe: New(nil)
// returns a Query with no parameters rather than panicking, so handlers never
// crash on a malformed request.
//
// Example usage:
//
//	q := qp.New(r.URL)
//	page := q.Int("page", qp.Default(1), qp.Min(1))
//	limit := q.Int("limit", qp.Default(20), qp.Between(1, 100))
func New(u *url.URL) *Query {
	if u == nil {
		return &Query{values: url.Values{}}
	}

	return &Query{values: u.Query()}
}

// From builds a Query from an already-decoded url.Values. Use it when the
// values are obtained elsewhere (for example r.Form) to avoid re-parsing. A nil
// map yields a Query with no parameters.
func From(v url.Values) *Query {
	if v == nil {
		return &Query{values: url.Values{}}
	}

	return &Query{values: v}
}

// Values returns an independent copy of the decoded query parameters. Mutating
// the returned map never affects the Query, so it is safe to hand out. Reading
// many keys is cheaper through the typed methods, which share the Query's
// snapshot without copying.
func (q *Query) Values() url.Values {
	out := make(url.Values, len(q.values))
	for k, v := range q.values {
		out[k] = slices.Clone(v)
	}

	return out
}

// Contains reports whether the parameter is present in the query, regardless of
// whether it carries a value.
//
// Example usage:
//
//	if q.Contains("id") {
//	    // the "id" parameter was supplied (possibly as "?id=")
//	}
func (q *Query) Contains(key string) bool {
	_, ok := q.values[key]
	return ok
}

// Empty reports whether the parameter is present but carries no value (as in
// "?id="). It returns false for an absent parameter; use Contains to tell the
// two apart.
//
// Example usage:
//
//	if q.Empty("id") {
//	    // "id" was supplied with an empty value
//	}
func (q *Query) Empty(key string) bool {
	v, ok := q.values[key]
	return ok && (len(v) == 0 || v[0] == "")
}

// Contains reports whether the parameter is present in u's query. It is a
// one-shot convenience for New(u).Contains(key) and is nil-safe.
func Contains(u *url.URL, key string) bool {
	return New(u).Contains(key)
}

// Empty reports whether the parameter is present in u's query but carries no
// value. It is a one-shot convenience for New(u).Empty(key) and is nil-safe.
func Empty(u *url.URL, key string) bool {
	return New(u).Empty(key)
}
