package qp

// ParseAs reads a parameter of an arbitrary type T using a caller-supplied
// parse function, returning the full Result. It is the generic escape hatch
// behind the built-in readers — use it for types the package does not cover
// (enums, net/netip addresses, custom value types). Methods cannot take type
// parameters in Go, so the As family is a set of free functions that take the
// Query explicitly.
//
// The options (Default, OneOf, and — for ordered T — Min/Max/Between) apply
// exactly as they do for the built-in readers.
//
// Example usage:
//
//	import "net/netip"
//
//	r := qp.ParseAs(q, "ip", netip.ParseAddr)
func ParseAs[T any](
	q *Query,
	key string,
	parse func(string) (T, error),
	opts ...Option[T],
) *Result[T] {
	return getScalar(q.values, key, parse, opts)
}

// GetAs reads a parameter of type T with a custom parser and reports whether it
// was present, non-empty, and valid.
//
// Example usage:
//
//	import "time"
//
//	ttl, ok := qp.GetAs(q, "ttl", time.ParseDuration, qp.Default(time.Minute))
func GetAs[T any](
	q *Query,
	key string,
	parse func(string) (T, error),
	opts ...Option[T],
) (T, bool) {
	r := getScalar(q.values, key, parse, opts)
	return r.Value, ok(r)
}

// PullAs reads a parameter of type T with a custom parser and returns a pointer
// to its value, or nil when the parameter is absent.
func PullAs[T any](
	q *Query,
	key string,
	parse func(string) (T, error),
	opts ...Option[T],
) *T {
	r := getScalar(q.values, key, parse, opts)
	if !r.Contains {
		return nil
	}

	return &r.Value
}

// ParseSliceAs reads a slice parameter of element type T with a custom
// element parser, returning the full Result. Splitting, trimming, the
// empty-element policy, and Each validation work as for the built-in slice
// readers.
func ParseSliceAs[T any](
	q *Query,
	key string,
	parse func(string) (T, error),
	opts ...SliceOption[T],
) *Result[[]T] {
	return getSlice(q.values, key, parse, opts)
}

// GetSliceAs reads a slice parameter of element type T with a custom element
// parser and reports whether it was present, non-empty, and valid.
func GetSliceAs[T any](
	q *Query,
	key string,
	parse func(string) (T, error),
	opts ...SliceOption[T],
) ([]T, bool) {
	r := getSlice(q.values, key, parse, opts)
	return r.Value, ok(r)
}

// PullSliceAs reads a slice parameter of element type T with a custom element
// parser and returns its values, or nil when the parameter is absent.
func PullSliceAs[T any](
	q *Query,
	key string,
	parse func(string) (T, error),
	opts ...SliceOption[T],
) []T {
	r := getSlice(q.values, key, parse, opts)
	if !r.Contains {
		return nil
	}

	return r.Value
}
