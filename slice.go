package qp

// SliceOption configures how a slice-valued parameter of element type T is
// split, defaulted, and validated. It is passed to the slice reading methods
// (IntSlice, GetStringSlice, ParseFloatSlice, …).
type SliceOption[T any] func(*sliceConfig[T])

// emptyPolicy decides what happens to an empty element produced by splitting
// (for example the gap in "a,,b" or a trailing comma).
type emptyPolicy int

const (
	emptySkip  emptyPolicy = iota // drop empty elements (default)
	emptyKeep                     // keep them and parse "" like any other value
	emptyError                    // reject the whole slice on any empty element
)

// sliceConfig is the resolved set of slice options.
type sliceConfig[T any] struct {
	sep      string
	fallback []T // default slice when absent, empty, or invalid (empty if unset)
	empty    emptyPolicy
	each     config[T] // per-element validation
}

// newSliceConfig folds the given options into a sliceConfig. The default
// separator is a comma and empty elements are skipped.
func newSliceConfig[T any](opts []SliceOption[T]) sliceConfig[T] {
	sc := sliceConfig[T]{sep: ","}
	for _, o := range opts {
		if o != nil {
			o(&sc)
		}
	}

	return sc
}

// Sep sets the separator used to split a single comma-style value such as
// "?ids=1,2,3". An empty separator is ignored and the default comma is kept.
//
// Example usage:
//
//	tags := q.StringSlice("tags", qp.Sep(";")) // "?tags=a;b;c" -> [a b c]
func Sep[T any](sep string) SliceOption[T] {
	return func(sc *sliceConfig[T]) {
		if sep != "" {
			sc.sep = sep
		}
	}
}

// Fallback sets the slice returned when the parameter is absent, empty, or
// invalid. Without it that fallback is an empty (non-nil) slice.
func Fallback[T any](v []T) SliceOption[T] {
	return func(sc *sliceConfig[T]) {
		sc.fallback = v
	}
}

// Each applies scalar options to every element of the slice. It is the way to
// validate slice elements against a range or an allowed set; a single invalid
// element rejects the whole slice (Error is set, Value becomes the fallback).
//
// Example usage:
//
//	ids := q.IntSlice("ids", qp.Each(qp.Between(1, 100)))
func Each[T any](opts ...Option[T]) SliceOption[T] {
	return func(sc *sliceConfig[T]) {
		sc.each = newConfig(opts)
	}
}

// KeepEmpty keeps empty elements produced by splitting instead of dropping
// them. For strings this preserves "" entries ("a,,b" -> ["a", "", "b"]); for
// numeric and boolean slices an empty element then fails to parse and rejects
// the slice.
func KeepEmpty[T any]() SliceOption[T] {
	return func(sc *sliceConfig[T]) {
		sc.empty = emptyKeep
	}
}

// ErrorOnEmpty rejects the whole slice if splitting yields any empty element,
// turning otherwise-tolerated noise such as a trailing or doubled comma into an
// error.
func ErrorOnEmpty[T any]() SliceOption[T] {
	return func(sc *sliceConfig[T]) {
		sc.empty = emptyError
	}
}
