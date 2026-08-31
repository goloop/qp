package qp

import "cmp"

// Option configures how a single scalar value of type T is defaulted and
// validated. Options are passed to the scalar reading methods (Int, GetInt,
// ParseInt, …) and, via Each, to the element validation of the slice methods.
type Option[T any] func(*config[T])

// config is the resolved set of scalar options.
type config[T any] struct {
	def T // value used when absent, empty, or invalid (zero value if unset)

	// Range bounds. compare is installed by Min/Max/Between (which require an
	// ordered T); it stays nil when no range option is used, so the generic
	// engine never needs T to be ordered. selfEq is installed alongside and is
	// false only for a value that does not equal itself (a NaN float), so a NaN
	// is never treated as within range.
	compare func(a, b T) int
	selfEq  func(v T) bool
	hasMin  bool
	min     T
	hasMax  bool
	max     T

	// Allowed set. eq is installed by OneOf (which requires a comparable T).
	// hasOneOf records that OneOf was applied at all, so an explicitly empty set
	// rejects everything instead of silently disabling validation.
	eq       func(a, b T) bool
	hasOneOf bool
	oneOf    []T

	// Clamp bounds. Unlike a range, which rejects an out-of-range value (falling
	// back to the default), a clamp folds it into [clampLo, clampHi]. Installed
	// by Clamp, which also installs compare/selfEq.
	hasClamp bool
	clampLo  T
	clampHi  T
}

// newConfig folds the given options into a config.
func newConfig[T any](opts []Option[T]) config[T] {
	var c config[T]
	for _, o := range opts {
		if o != nil {
			o(&c)
		}
	}

	return c
}

// clamp folds v into [clampLo, clampHi] when Clamp was set, and returns v
// unchanged otherwise. A value that does not equal itself (a NaN float) folds
// to the low bound rather than escaping the range.
func (c *config[T]) clamp(v T) T {
	if !c.hasClamp {
		return v
	}
	if c.selfEq != nil && !c.selfEq(v) {
		return c.clampLo
	}
	if c.compare(v, c.clampLo) < 0 {
		return c.clampLo
	}
	if c.compare(v, c.clampHi) > 0 {
		return c.clampHi
	}
	return v
}

// valid reports whether v satisfies the configured constraints. With no range
// and no allowed set, every value is valid. A range and a set combine as a
// union: a value is valid if it falls within the range OR appears in the set
// (the set acts as explicit exceptions to the range).
func (c *config[T]) valid(v T) bool {
	hasRange := c.hasMin || c.hasMax
	hasSet := c.hasOneOf
	if !hasRange && !hasSet {
		return true
	}

	if hasRange {
		ok := true
		if c.selfEq != nil && !c.selfEq(v) {
			ok = false // a NaN is not within any range
		}
		if c.hasMin && c.compare(v, c.min) < 0 {
			ok = false
		}
		if c.hasMax && c.compare(v, c.max) > 0 {
			ok = false
		}
		if ok {
			return true
		}
	}

	if hasSet {
		for _, o := range c.oneOf {
			if c.eq(v, o) {
				return true
			}
		}
	}

	return false
}

// Default sets the value returned when the parameter is absent, empty, or
// invalid. Without it, that fallback is the zero value of T.
func Default[T any](v T) Option[T] {
	return func(c *config[T]) {
		c.def = v
	}
}

// Min rejects values less than v. It can be combined with Max or used alone for
// an open-ended lower bound.
//
// Note the open end: Min(1) floors the value but does not cap it, so a parsed
// value can still be arbitrarily large. When the value then feeds arithmetic —
// a page number multiplied into a SQL OFFSET, say — that unbounded top can
// overflow. For a bounded range use Between; to fold an out-of-range value to
// the nearest bound instead of rejecting it, use Clamp.
func Min[T cmp.Ordered](v T) Option[T] {
	return func(c *config[T]) {
		c.hasMin = true
		c.min = v
		c.compare = cmp.Compare[T]
		c.selfEq = func(v T) bool { return v == v }
	}
}

// Max rejects values greater than v. It can be combined with Min or used alone
// for an open-ended upper bound.
func Max[T cmp.Ordered](v T) Option[T] {
	return func(c *config[T]) {
		c.hasMax = true
		c.max = v
		c.compare = cmp.Compare[T]
		c.selfEq = func(v T) bool { return v == v }
	}
}

// Between rejects values outside the inclusive range [lo, hi]. The bounds are
// swapped automatically when lo > hi, so Between(65, 18) and Between(18, 65)
// behave identically.
func Between[T cmp.Ordered](lo, hi T) Option[T] {
	if lo > hi {
		lo, hi = hi, lo
	}

	return func(c *config[T]) {
		c.hasMin, c.min = true, lo
		c.hasMax, c.max = true, hi
		c.compare = cmp.Compare[T]
		c.selfEq = func(v T) bool { return v == v }
	}
}

// Clamp folds an out-of-range value into the inclusive range [lo, hi] instead
// of rejecting it. Where Min/Max/Between reject a value outside the range — so
// ?page=999999999 with Between(1, N) falls back to the default, and Min(1)
// alone leaves the value unbounded above — Clamp keeps the value usable: it
// returns the nearest bound. For a page number that means the last page rather
// than a fallback or an overflow, which is almost always what a paginated
// endpoint wants. The bounds are swapped automatically when lo > hi.
func Clamp[T cmp.Ordered](lo, hi T) Option[T] {
	if lo > hi {
		lo, hi = hi, lo
	}
	return func(c *config[T]) {
		c.hasClamp = true
		c.clampLo = lo
		c.clampHi = hi
		c.compare = cmp.Compare[T]
		c.selfEq = func(v T) bool { return v == v }
	}
}

// OneOf restricts the value to the given set. When combined with a range the
// two form a union: a value is accepted if it is within the range OR in the
// set, so the set names explicit exceptions to the range. For strings this is
// the primary way to validate against a list of allowed values.
//
// An explicitly empty OneOf() rejects every value (an empty allowed set means
// "nothing is allowed", not "no constraint") — pass no OneOf option at all to
// skip set validation.
func OneOf[T comparable](vs ...T) Option[T] {
	return func(c *config[T]) {
		c.hasOneOf = true
		c.oneOf = append(c.oneOf, vs...)
		c.eq = func(a, b T) bool { return a == b }
	}
}

// OneOfFunc restricts the value to the given set, comparing values with eq. It
// is the equality-customising form of OneOf for types whose natural equality is
// not ==, most notably time.Time (pass time.Time.Equal), where two values can
// denote the same instant yet differ bit for bit. It combines with a range as a
// union, exactly like OneOf, and an explicitly empty set rejects every value.
func OneOfFunc[T any](eq func(a, b T) bool, vs ...T) Option[T] {
	return func(c *config[T]) {
		c.hasOneOf = true
		c.oneOf = append(c.oneOf, vs...)
		c.eq = eq
	}
}
