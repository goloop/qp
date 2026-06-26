package qp

import "net/url"

// ParseFloat reads a float64 parameter and returns the full Result. Value holds
// the parsed number, or the default when the parameter is absent, empty, or
// rejected by validation; Error is set only when a present, non-empty value
// fails to parse or violates a Min/Max/Between/OneOf constraint.
//
// Example usage:
//
//	r := q.ParseFloat("temperature", qp.Between(36.0, 42.0))
func (q *Query) ParseFloat(key string, opts ...Option[float64]) *Result[float64] {
	return getScalar(q.values, key, parseFloat, opts)
}

// Float reads a float64 parameter and returns its value, falling back to the
// default (zero unless Default is given) when the parameter is absent, empty,
// or invalid.
//
// Example usage:
//
//	ratio := q.Float("ratio", qp.Default(1.0), qp.Between(0, 1))
func (q *Query) Float(key string, opts ...Option[float64]) float64 {
	return q.ParseFloat(key, opts...).Value
}

// GetFloat reads a float64 parameter and reports whether it was present,
// non-empty, and valid. On failure the returned value is the default.
func (q *Query) GetFloat(key string, opts ...Option[float64]) (float64, bool) {
	r := q.ParseFloat(key, opts...)
	return r.Value, ok(r)
}

// PullFloat reads a float64 parameter and returns a pointer to its value, or
// nil when the parameter is absent. Use ParseFloat when the validation error
// matters.
func (q *Query) PullFloat(key string, opts ...Option[float64]) *float64 {
	r := q.ParseFloat(key, opts...)
	if !r.Contains {
		return nil
	}

	return &r.Value
}

// ParseFloatSlice reads a float64 slice parameter and returns the full Result.
// Both "?xs=1.1,2.2" and "?xs=1.1&xs=2.2" are accepted; a single unparseable or
// invalid element rejects the whole slice.
func (q *Query) ParseFloatSlice(
	key string,
	opts ...SliceOption[float64],
) *Result[[]float64] {
	return getSlice(q.values, key, parseFloat, opts)
}

// FloatSlice reads a float64 slice parameter and returns its values, falling
// back to the fallback (an empty slice unless Fallback is given) when the
// parameter is absent, empty, or invalid.
func (q *Query) FloatSlice(key string, opts ...SliceOption[float64]) []float64 {
	return q.ParseFloatSlice(key, opts...).Value
}

// GetFloatSlice reads a float64 slice parameter and reports whether it was
// present, non-empty, and valid.
func (q *Query) GetFloatSlice(
	key string,
	opts ...SliceOption[float64],
) ([]float64, bool) {
	r := q.ParseFloatSlice(key, opts...)
	return r.Value, ok(r)
}

// PullFloatSlice reads a float64 slice parameter and returns its values, or nil
// when the parameter is absent (even if a fallback is configured).
func (q *Query) PullFloatSlice(
	key string,
	opts ...SliceOption[float64],
) []float64 {
	r := q.ParseFloatSlice(key, opts...)
	if !r.Contains {
		return nil
	}

	return r.Value
}

// ParseFloat is a one-shot convenience for New(u).ParseFloat(key, opts...).
func ParseFloat(u *url.URL, key string, opts ...Option[float64]) *Result[float64] {
	return New(u).ParseFloat(key, opts...)
}

// Float is a one-shot convenience for New(u).Float(key, opts...).
func Float(u *url.URL, key string, opts ...Option[float64]) float64 {
	return New(u).Float(key, opts...)
}

// GetFloat is a one-shot convenience for New(u).GetFloat(key, opts...).
func GetFloat(u *url.URL, key string, opts ...Option[float64]) (float64, bool) {
	return New(u).GetFloat(key, opts...)
}

// PullFloat is a one-shot convenience for New(u).PullFloat(key, opts...).
func PullFloat(u *url.URL, key string, opts ...Option[float64]) *float64 {
	return New(u).PullFloat(key, opts...)
}

// ParseFloatSlice is a one-shot convenience for New(u).ParseFloatSlice(...).
func ParseFloatSlice(
	u *url.URL,
	key string,
	opts ...SliceOption[float64],
) *Result[[]float64] {
	return New(u).ParseFloatSlice(key, opts...)
}

// FloatSlice is a one-shot convenience for New(u).FloatSlice(key, opts...).
func FloatSlice(u *url.URL, key string, opts ...SliceOption[float64]) []float64 {
	return New(u).FloatSlice(key, opts...)
}

// GetFloatSlice is a one-shot convenience for New(u).GetFloatSlice(...).
func GetFloatSlice(
	u *url.URL,
	key string,
	opts ...SliceOption[float64],
) ([]float64, bool) {
	return New(u).GetFloatSlice(key, opts...)
}

// PullFloatSlice is a one-shot convenience for New(u).PullFloatSlice(...).
func PullFloatSlice(
	u *url.URL,
	key string,
	opts ...SliceOption[float64],
) []float64 {
	return New(u).PullFloatSlice(key, opts...)
}
