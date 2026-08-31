package qp

import (
	"fmt"
	"net/url"
	"slices"
	"strconv"
	"strings"
)

// getScalar is the shared engine behind every scalar reader. parse converts the
// raw query value into T; the option set provides the default and validation.
func getScalar[T any](
	values url.Values,
	key string,
	parse func(string) (T, error),
	opts []Option[T],
) *Result[T] {
	c := newConfig(opts)
	r := &Result[T]{Key: key, Default: c.def, Value: c.def, Contains: true}

	data, ok := values[key]
	if !ok {
		// The parameter is absent.
		r.Contains = false
		r.Empty = true
		return r
	}
	if len(data) == 0 || data[0] == "" {
		// The parameter is present but carries no value.
		r.Empty = true
		return r
	}

	value, err := parse(data[0])
	if err != nil {
		r.Error = invalidErr(key, data[0])
		return r
	}
	// Clamp (when set) folds the value into range before validation, so a
	// clamped reader never rejects a parseable value — it returns the bound.
	value = c.clamp(value)
	if !c.valid(value) {
		r.Error = notAllowedErr(key, data[0])
		return r
	}

	r.Value = value
	return r
}

// getSlice is the shared engine behind every slice reader. Both comma-style
// values ("?ids=1,2,3") and repeated values ("?ids=1&ids=2") are normalised
// into one stream of elements: each raw value is split by the separator, every
// element is trimmed, and the empty-element policy is applied before parsing.
func getSlice[T any](
	values url.Values,
	key string,
	parse func(string) (T, error),
	opts []SliceOption[T],
) *Result[[]T] {
	sc := newSliceConfig(opts)
	r := &Result[[]T]{
		Key:      key,
		Default:  cloneOrEmpty(sc.fallback),
		Value:    cloneOrEmpty(sc.fallback),
		Contains: true,
	}

	data, ok := values[key]
	if !ok {
		// The parameter is absent.
		r.Contains = false
		r.Empty = true
		return r
	}
	if len(data) == 0 || (len(data) == 1 && data[0] == "") {
		// The parameter is present but empty: either a single empty value
		// ("?ids=") or a key whose value slice is empty (possible via From with
		// a hand-built url.Values). Both match the scalar engine and Query.Empty.
		// Note: a leading empty value among several ("?ids=&ids=5") is NOT
		// treated as empty here — it flows into the loop below so the
		// non-empty values survive, per the documented normalisation.
		r.Empty = true
		return r
	}

	out := make([]T, 0, len(data))
	for _, raw := range data {
		for part := range strings.SplitSeq(raw, sc.sep) {
			part = strings.TrimSpace(part)
			if part == "" {
				switch sc.empty {
				case emptyError:
					r.Error = emptyElementErr(key)
					r.Value = cloneOrEmpty(sc.fallback)
					return r
				case emptySkip:
					continue
				}
				// emptyKeep falls through and parses "".
			}

			value, err := parse(part)
			if err != nil {
				r.Error = invalidErr(key, part)
				r.Value = cloneOrEmpty(sc.fallback)
				return r
			}
			if !sc.each.valid(value) {
				r.Error = notAllowedErr(key, part)
				r.Value = cloneOrEmpty(sc.fallback)
				return r
			}

			out = append(out, value)
		}
	}

	r.Value = out
	return r
}

// cloneOrEmpty returns an independent copy of the fallback slice, or a fresh
// empty (non-nil) slice when no fallback was configured. Cloning keeps each read
// independent: a returned default can never be aliased by the caller's slice or
// shared between two reads.
func cloneOrEmpty[T any](fallback []T) []T {
	if fallback == nil {
		return []T{}
	}

	return slices.Clone(fallback)
}

// ok reports the success flag used by the Get* family: the parameter is
// present, non-empty, and survived parsing and validation.
func ok[T any](r *Result[T]) bool {
	return r.Contains && !r.Empty && r.Error == nil
}

// parseString is the identity parser; a raw string is always a valid string.
func parseString(s string) (string, error) {
	return s, nil
}

// parseInt parses a base-10 integer.
func parseInt(s string) (int, error) {
	return strconv.Atoi(s)
}

// parseFloat parses a 64-bit floating-point number.
func parseFloat(s string) (float64, error) {
	return strconv.ParseFloat(s, 64)
}

// parseBool parses a boolean, accepting true/false, yes/no, on/off and 1/0
// without regard to case.
func parseBool(s string) (bool, error) {
	switch strings.ToLower(s) {
	case "1", "true", "yes", "on":
		return true, nil
	case "0", "false", "no", "off":
		return false, nil
	default:
		return false, fmt.Errorf("invalid bool value: %q", s)
	}
}
