package qp

import "net/url"

// ParseInt reads an integer parameter and returns the full Result. Value holds
// the parsed integer, or the default when the parameter is absent, empty, or
// rejected by validation; Error is set only when a present, non-empty value
// fails to parse or violates a Min/Max/Between/OneOf constraint.
//
// Example usage:
//
//	r := q.ParseInt("age", qp.Default(18), qp.Between(18, 65))
//	// "?age=99": r.Value == 18, r.Error != nil
func (q *Query) ParseInt(key string, opts ...Option[int]) *Result[int] {
	return getScalar(q.values, key, parseInt, opts)
}

// Int reads an integer parameter and returns its value, falling back to the
// default (zero unless Default is given) when the parameter is absent, empty,
// or invalid. Use it for the common "give me a usable number" case.
//
// Example usage:
//
//	page := q.Int("page", qp.Default(1), qp.Min(1))
func (q *Query) Int(key string, opts ...Option[int]) int {
	return q.ParseInt(key, opts...).Value
}

// GetInt reads an integer parameter and reports whether it was present,
// non-empty, and valid. On failure the returned value is the default.
//
// Example usage:
//
//	age, ok := q.GetInt("age", qp.Between(18, 65))
func (q *Query) GetInt(key string, opts ...Option[int]) (int, bool) {
	r := q.ParseInt(key, opts...)
	return r.Value, ok(r)
}

// PullInt reads an integer parameter and returns a pointer to its value, or nil
// when the parameter is absent. A present parameter always yields a non-nil
// pointer (to the parsed value, or to the default when invalid); use ParseInt
// when the validation error matters.
//
// Example usage:
//
//	if min := q.PullInt("min_age"); min != nil {
//	    // "min_age" was supplied
//	}
func (q *Query) PullInt(key string, opts ...Option[int]) *int {
	r := q.ParseInt(key, opts...)
	if !r.Contains {
		return nil
	}

	return &r.Value
}

// ParseIntSlice reads an integer slice parameter and returns the full Result.
// Both "?ids=1,2,3" and "?ids=1&ids=2&ids=3" are accepted and normalised the
// same way; a single unparseable or invalid element rejects the whole slice
// (Error is set, Value becomes the fallback).
//
// Example usage:
//
//	r := q.ParseIntSlice("ids", qp.Each(qp.Between(1, 100)))
func (q *Query) ParseIntSlice(
	key string,
	opts ...SliceOption[int],
) *Result[[]int] {
	return getSlice(q.values, key, parseInt, opts)
}

// IntSlice reads an integer slice parameter and returns its values, falling
// back to the fallback (an empty slice unless Fallback is given) when the
// parameter is absent, empty, or invalid.
//
// Example usage:
//
//	ids := q.IntSlice("ids", qp.Each(qp.Between(1, 100)))
func (q *Query) IntSlice(key string, opts ...SliceOption[int]) []int {
	return q.ParseIntSlice(key, opts...).Value
}

// GetIntSlice reads an integer slice parameter and reports whether it was
// present, non-empty, and valid.
func (q *Query) GetIntSlice(
	key string,
	opts ...SliceOption[int],
) ([]int, bool) {
	r := q.ParseIntSlice(key, opts...)
	return r.Value, ok(r)
}

// PullIntSlice reads an integer slice parameter and returns its values, or nil
// when the parameter is absent (even if a fallback is configured).
func (q *Query) PullIntSlice(key string, opts ...SliceOption[int]) []int {
	r := q.ParseIntSlice(key, opts...)
	if !r.Contains {
		return nil
	}

	return r.Value
}

// ParseInt is a one-shot convenience for New(u).ParseInt(key, opts...).
func ParseInt(u *url.URL, key string, opts ...Option[int]) *Result[int] {
	return New(u).ParseInt(key, opts...)
}

// Int is a one-shot convenience for New(u).Int(key, opts...).
func Int(u *url.URL, key string, opts ...Option[int]) int {
	return New(u).Int(key, opts...)
}

// GetInt is a one-shot convenience for New(u).GetInt(key, opts...).
func GetInt(u *url.URL, key string, opts ...Option[int]) (int, bool) {
	return New(u).GetInt(key, opts...)
}

// PullInt is a one-shot convenience for New(u).PullInt(key, opts...).
func PullInt(u *url.URL, key string, opts ...Option[int]) *int {
	return New(u).PullInt(key, opts...)
}

// ParseIntSlice is a one-shot convenience for New(u).ParseIntSlice(...).
func ParseIntSlice(
	u *url.URL,
	key string,
	opts ...SliceOption[int],
) *Result[[]int] {
	return New(u).ParseIntSlice(key, opts...)
}

// IntSlice is a one-shot convenience for New(u).IntSlice(key, opts...).
func IntSlice(u *url.URL, key string, opts ...SliceOption[int]) []int {
	return New(u).IntSlice(key, opts...)
}

// GetIntSlice is a one-shot convenience for New(u).GetIntSlice(...).
func GetIntSlice(
	u *url.URL,
	key string,
	opts ...SliceOption[int],
) ([]int, bool) {
	return New(u).GetIntSlice(key, opts...)
}

// PullIntSlice is a one-shot convenience for New(u).PullIntSlice(...).
func PullIntSlice(u *url.URL, key string, opts ...SliceOption[int]) []int {
	return New(u).PullIntSlice(key, opts...)
}
