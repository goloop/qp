package qp

import "net/url"

// ParseBool reads a boolean parameter and returns the full Result. The accepted
// forms are true/false, yes/no, on/off and 1/0, case-insensitively. Value holds
// the parsed boolean, or the default when the parameter is absent, empty, or
// unparseable; Error is set only when a present, non-empty value is not a
// recognised boolean.
//
// Range options (Min/Max/Between) do not apply to booleans; only Default and
// OneOf are meaningful.
//
// Example usage:
//
//	r := q.ParseBool("active", qp.Default(true))
func (q *Query) ParseBool(key string, opts ...Option[bool]) *Result[bool] {
	return getScalar(q.values, key, parseBool, opts)
}

// Bool reads a boolean parameter and returns its value, falling back to the
// default (false unless Default is given) when the parameter is absent, empty,
// or unparseable.
//
// Example usage:
//
//	verbose := q.Bool("verbose")
func (q *Query) Bool(key string, opts ...Option[bool]) bool {
	return q.ParseBool(key, opts...).Value
}

// GetBool reads a boolean parameter and reports whether it was present,
// non-empty, and a recognised boolean. On failure the returned value is the
// default.
func (q *Query) GetBool(key string, opts ...Option[bool]) (bool, bool) {
	r := q.ParseBool(key, opts...)
	return r.Value, ok(r)
}

// PullBool reads a boolean parameter and returns a pointer to its value, or nil
// when the parameter is absent. This is the natural fit for tri-state filters
// where "unset" must be told apart from "false".
//
// Example usage:
//
//	if active := q.PullBool("is_active"); active != nil {
//	    // filter on is_active == *active
//	}
func (q *Query) PullBool(key string, opts ...Option[bool]) *bool {
	r := q.ParseBool(key, opts...)
	if !r.Contains {
		return nil
	}

	return &r.Value
}

// ParseBoolSlice reads a boolean slice parameter and returns the full Result.
// Both "?flags=true,false" and "?flags=true&flags=false" are accepted; a single
// unparseable element rejects the whole slice.
func (q *Query) ParseBoolSlice(
	key string,
	opts ...SliceOption[bool],
) *Result[[]bool] {
	return getSlice(q.values, key, parseBool, opts)
}

// BoolSlice reads a boolean slice parameter and returns its values, falling
// back to the fallback (an empty slice unless Fallback is given) when the
// parameter is absent, empty, or invalid.
func (q *Query) BoolSlice(key string, opts ...SliceOption[bool]) []bool {
	return q.ParseBoolSlice(key, opts...).Value
}

// GetBoolSlice reads a boolean slice parameter and reports whether it was
// present, non-empty, and valid.
func (q *Query) GetBoolSlice(
	key string,
	opts ...SliceOption[bool],
) ([]bool, bool) {
	r := q.ParseBoolSlice(key, opts...)
	return r.Value, ok(r)
}

// PullBoolSlice reads a boolean slice parameter and returns its values, or nil
// when the parameter is absent (even if a fallback is configured).
func (q *Query) PullBoolSlice(key string, opts ...SliceOption[bool]) []bool {
	r := q.ParseBoolSlice(key, opts...)
	if !r.Contains {
		return nil
	}

	return r.Value
}

// ParseBool is a one-shot convenience for New(u).ParseBool(key, opts...).
func ParseBool(u *url.URL, key string, opts ...Option[bool]) *Result[bool] {
	return New(u).ParseBool(key, opts...)
}

// Bool is a one-shot convenience for New(u).Bool(key, opts...).
func Bool(u *url.URL, key string, opts ...Option[bool]) bool {
	return New(u).Bool(key, opts...)
}

// GetBool is a one-shot convenience for New(u).GetBool(key, opts...).
func GetBool(u *url.URL, key string, opts ...Option[bool]) (bool, bool) {
	return New(u).GetBool(key, opts...)
}

// PullBool is a one-shot convenience for New(u).PullBool(key, opts...).
func PullBool(u *url.URL, key string, opts ...Option[bool]) *bool {
	return New(u).PullBool(key, opts...)
}

// ParseBoolSlice is a one-shot convenience for New(u).ParseBoolSlice(...).
func ParseBoolSlice(
	u *url.URL,
	key string,
	opts ...SliceOption[bool],
) *Result[[]bool] {
	return New(u).ParseBoolSlice(key, opts...)
}

// BoolSlice is a one-shot convenience for New(u).BoolSlice(key, opts...).
func BoolSlice(u *url.URL, key string, opts ...SliceOption[bool]) []bool {
	return New(u).BoolSlice(key, opts...)
}

// GetBoolSlice is a one-shot convenience for New(u).GetBoolSlice(...).
func GetBoolSlice(
	u *url.URL,
	key string,
	opts ...SliceOption[bool],
) ([]bool, bool) {
	return New(u).GetBoolSlice(key, opts...)
}

// PullBoolSlice is a one-shot convenience for New(u).PullBoolSlice(...).
func PullBoolSlice(u *url.URL, key string, opts ...SliceOption[bool]) []bool {
	return New(u).PullBoolSlice(key, opts...)
}
