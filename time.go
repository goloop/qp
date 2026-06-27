package qp

import (
	"net/url"
	"time"
)

// timeParser builds a parse function for the given reference layout.
func timeParser(layout string) func(string) (time.Time, error) {
	return func(s string) (time.Time, error) {
		return time.Parse(layout, s)
	}
}

// ParseDuration reads a time.Duration parameter (as accepted by
// time.ParseDuration, e.g. "300ms", "1.5h") and returns the full Result. Since
// Duration is an ordered type, Min/Max/Between apply as usual.
//
// Example usage:
//
//	r := q.ParseDuration("ttl", qp.Default(time.Minute), qp.Max(time.Hour))
func (q *Query) ParseDuration(
	key string,
	opts ...Option[time.Duration],
) *Result[time.Duration] {
	return getScalar(q.values, key, time.ParseDuration, opts)
}

// Duration reads a time.Duration parameter and returns its value, falling back
// to the default when the parameter is absent, empty, or invalid.
func (q *Query) Duration(
	key string,
	opts ...Option[time.Duration],
) time.Duration {
	return q.ParseDuration(key, opts...).Value
}

// GetDuration reads a time.Duration parameter and reports whether it was
// present, non-empty, and valid.
func (q *Query) GetDuration(
	key string,
	opts ...Option[time.Duration],
) (time.Duration, bool) {
	r := q.ParseDuration(key, opts...)
	return r.Value, ok(r)
}

// PullDuration reads a time.Duration parameter and returns a pointer to its
// value, or nil when the parameter is absent.
func (q *Query) PullDuration(
	key string,
	opts ...Option[time.Duration],
) *time.Duration {
	r := q.ParseDuration(key, opts...)
	if !r.Contains {
		return nil
	}

	return &r.Value
}

// ParseTime reads a time.Time parameter using the given reference layout (as
// accepted by time.Parse, e.g. time.RFC3339) and returns the full Result.
// Range options do not apply to time.Time; Default and OneOf do.
//
// Example usage:
//
//	r := q.ParseTime("from", time.RFC3339)
func (q *Query) ParseTime(
	key, layout string,
	opts ...Option[time.Time],
) *Result[time.Time] {
	return getScalar(q.values, key, timeParser(layout), opts)
}

// Time reads a time.Time parameter and returns its value, falling back to the
// default (the zero time unless Default is given) when the parameter is absent,
// empty, or invalid.
func (q *Query) Time(
	key, layout string,
	opts ...Option[time.Time],
) time.Time {
	return q.ParseTime(key, layout, opts...).Value
}

// GetTime reads a time.Time parameter and reports whether it was present,
// non-empty, and valid.
func (q *Query) GetTime(
	key, layout string,
	opts ...Option[time.Time],
) (time.Time, bool) {
	r := q.ParseTime(key, layout, opts...)
	return r.Value, ok(r)
}

// PullTime reads a time.Time parameter and returns a pointer to its value, or
// nil when the parameter is absent.
func (q *Query) PullTime(
	key, layout string,
	opts ...Option[time.Time],
) *time.Time {
	r := q.ParseTime(key, layout, opts...)
	if !r.Contains {
		return nil
	}

	return &r.Value
}

// ParseDuration is a one-shot convenience for New(u).ParseDuration(...).
func ParseDuration(
	u *url.URL,
	key string,
	opts ...Option[time.Duration],
) *Result[time.Duration] {
	return New(u).ParseDuration(key, opts...)
}

// Duration is a one-shot convenience for New(u).Duration(key, opts...).
func Duration(u *url.URL, key string, opts ...Option[time.Duration]) time.Duration {
	return New(u).Duration(key, opts...)
}

// GetDuration is a one-shot convenience for New(u).GetDuration(...).
func GetDuration(
	u *url.URL,
	key string,
	opts ...Option[time.Duration],
) (time.Duration, bool) {
	return New(u).GetDuration(key, opts...)
}

// PullDuration is a one-shot convenience for New(u).PullDuration(...).
func PullDuration(
	u *url.URL,
	key string,
	opts ...Option[time.Duration],
) *time.Duration {
	return New(u).PullDuration(key, opts...)
}

// ParseTime is a one-shot convenience for New(u).ParseTime(key, layout, opts...).
func ParseTime(
	u *url.URL,
	key, layout string,
	opts ...Option[time.Time],
) *Result[time.Time] {
	return New(u).ParseTime(key, layout, opts...)
}

// Time is a one-shot convenience for New(u).Time(key, layout, opts...).
func Time(u *url.URL, key, layout string, opts ...Option[time.Time]) time.Time {
	return New(u).Time(key, layout, opts...)
}

// GetTime is a one-shot convenience for New(u).GetTime(key, layout, opts...).
func GetTime(
	u *url.URL,
	key, layout string,
	opts ...Option[time.Time],
) (time.Time, bool) {
	return New(u).GetTime(key, layout, opts...)
}

// PullTime is a one-shot convenience for New(u).PullTime(key, layout, opts...).
func PullTime(
	u *url.URL,
	key, layout string,
	opts ...Option[time.Time],
) *time.Time {
	return New(u).PullTime(key, layout, opts...)
}
