package qp

import (
	"fmt"
	"net/url"
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
		r.Error = fmt.Errorf("invalid value for key %q: %q", key, data[0])
		return r
	}
	if !c.valid(value) {
		r.Error = fmt.Errorf("value not allowed for key %q: %q", key, data[0])
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

	def := sc.fallback
	if def == nil {
		def = []T{}
	}
	r := &Result[[]T]{Key: key, Default: def, Value: def, Contains: true}

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

	out := make([]T, 0, len(data))
	for _, raw := range data {
		for _, part := range strings.Split(raw, sc.sep) {
			part = strings.TrimSpace(part)
			if part == "" {
				switch sc.empty {
				case emptyError:
					r.Error = fmt.Errorf("empty element for key %q", key)
					r.Value = def
					return r
				case emptySkip:
					continue
				}
				// emptyKeep falls through and parses "".
			}

			value, err := parse(part)
			if err != nil {
				r.Error = fmt.Errorf("invalid value for key %q: %q", key, part)
				r.Value = def
				return r
			}
			if !sc.each.valid(value) {
				r.Error = fmt.Errorf("value not allowed for key %q: %q", key, part)
				r.Value = def
				return r
			}

			out = append(out, value)
		}
	}

	r.Value = out
	return r
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
