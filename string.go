package qp

import "net/url"

// ParseString reads a string parameter and returns the full Result. Value holds
// the raw string, or the default when the parameter is absent, empty, or
// rejected by validation; Error is set only when a present, non-empty value is
// not in the allowed set given by OneOf.
//
// Example usage:
//
//	r := q.ParseString("role", qp.Default("guest"), qp.OneOf("admin", "user"))
//	// "?role=hacker": r.Value == "guest", r.Error != nil
func (q *Query) ParseString(key string, opts ...Option[string]) *Result[string] {
	return getScalar(q.values, key, parseString, opts)
}

// String reads a string parameter and returns its value, falling back to the
// default (empty string unless Default is given) when the parameter is absent,
// empty, or invalid.
//
// Example usage:
//
//	sort := q.String("sort", qp.Default("created_at"))
func (q *Query) String(key string, opts ...Option[string]) string {
	return q.ParseString(key, opts...).Value
}

// GetString reads a string parameter and reports whether it was present,
// non-empty, and (if OneOf is given) allowed. On failure the returned value is
// the default.
//
// Example usage:
//
//	role, ok := q.GetString("role", qp.OneOf("admin", "user"))
func (q *Query) GetString(key string, opts ...Option[string]) (string, bool) {
	r := q.ParseString(key, opts...)
	return r.Value, ok(r)
}

// PullString reads a string parameter and returns a pointer to its value, or
// nil when the parameter is absent. Use ParseString when the validation error
// matters.
func (q *Query) PullString(key string, opts ...Option[string]) *string {
	r := q.ParseString(key, opts...)
	if !r.Contains {
		return nil
	}

	return &r.Value
}

// ParseStringSlice reads a string slice parameter and returns the full Result.
// Both "?names=a,b,c" and "?names=a&names=b&names=c" are accepted; with Each a
// single disallowed element rejects the whole slice.
func (q *Query) ParseStringSlice(
	key string,
	opts ...SliceOption[string],
) *Result[[]string] {
	return getSlice(q.values, key, parseString, opts)
}

// StringSlice reads a string slice parameter and returns its values, falling
// back to the fallback (an empty slice unless Fallback is given) when the
// parameter is absent, empty, or invalid.
//
// Example usage:
//
//	tags := q.StringSlice("tag")
func (q *Query) StringSlice(key string, opts ...SliceOption[string]) []string {
	return q.ParseStringSlice(key, opts...).Value
}

// GetStringSlice reads a string slice parameter and reports whether it was
// present, non-empty, and valid.
func (q *Query) GetStringSlice(
	key string,
	opts ...SliceOption[string],
) ([]string, bool) {
	r := q.ParseStringSlice(key, opts...)
	return r.Value, ok(r)
}

// PullStringSlice reads a string slice parameter and returns its values, or nil
// when the parameter is absent (even if a fallback is configured).
func (q *Query) PullStringSlice(
	key string,
	opts ...SliceOption[string],
) []string {
	r := q.ParseStringSlice(key, opts...)
	if !r.Contains {
		return nil
	}

	return r.Value
}

// ParseString is a one-shot convenience for New(u).ParseString(key, opts...).
func ParseString(u *url.URL, key string, opts ...Option[string]) *Result[string] {
	return New(u).ParseString(key, opts...)
}

// String is a one-shot convenience for New(u).String(key, opts...).
func String(u *url.URL, key string, opts ...Option[string]) string {
	return New(u).String(key, opts...)
}

// GetString is a one-shot convenience for New(u).GetString(key, opts...).
func GetString(u *url.URL, key string, opts ...Option[string]) (string, bool) {
	return New(u).GetString(key, opts...)
}

// PullString is a one-shot convenience for New(u).PullString(key, opts...).
func PullString(u *url.URL, key string, opts ...Option[string]) *string {
	return New(u).PullString(key, opts...)
}

// ParseStringSlice is a one-shot convenience for New(u).ParseStringSlice(...).
func ParseStringSlice(
	u *url.URL,
	key string,
	opts ...SliceOption[string],
) *Result[[]string] {
	return New(u).ParseStringSlice(key, opts...)
}

// StringSlice is a one-shot convenience for New(u).StringSlice(key, opts...).
func StringSlice(u *url.URL, key string, opts ...SliceOption[string]) []string {
	return New(u).StringSlice(key, opts...)
}

// GetStringSlice is a one-shot convenience for New(u).GetStringSlice(...).
func GetStringSlice(
	u *url.URL,
	key string,
	opts ...SliceOption[string],
) ([]string, bool) {
	return New(u).GetStringSlice(key, opts...)
}

// PullStringSlice is a one-shot convenience for New(u).PullStringSlice(...).
func PullStringSlice(
	u *url.URL,
	key string,
	opts ...SliceOption[string],
) []string {
	return New(u).PullStringSlice(key, opts...)
}
