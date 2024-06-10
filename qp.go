package qp

import (
	"net/url"
)

// Value defines an interface for types supported by the query parsing
// functions. It restricts type parameters to basic types and their slices
// that can be parsed from query parameters in an HTTP request. This includes
// integers, floating-point numbers, strings, and booleans, along with their
// slice counterparts.
//
// This type constraint ensures that the parsing functions can safely operate
// on the provided type parameters without runtime type errors, leveraging
// Go's type safety features in a generic programming context.
//
// Example types supported:
//
//   - int, float64, string, bool
//   - []int, []float64, []string, []bool
type Value interface {
	~int | ~float64 | ~string | ~bool |
		~[]int | ~[]float64 | ~[]string | ~[]bool
}

// Result is a generic type to hold parsed query parameter values.
type Result[T Value] struct {
	Key   string // the query parameter name
	Value T      // the parsed query parameter value

	Default T   // the default value for the query parameter
	Min     T   // the minimum value for the query parameter
	Max     T   // the maximum value for the query parameter
	Others  []T // additional valid values for the query parameter

	Empty    bool  // indicates if the query parameter is empty
	Contains bool  // indicates if the query parameter is present
	Error    error // the error encountered during parsing
}

// Contains checks if a specified query parameter is present in the request.
// It returns true if the parameter is present, regardless of whether it has
// a value or not.
//
// Example usage:
//
//	if qp.Contains(r.URL, "id") {
//	    fmt.Println("ID parameter is present")
//	}
func Contains(u *url.URL, key string) bool {
	_, present := u.Query()[key]
	return present
}

// Empty checks if a specified query parameter is present in the request and
// has an empty value. It returns true if the parameter is present and has no
// value.
//
// Example usage:
//
//	if qp.Empty(r.URL, "id") {
//	    fmt.Println("ID parameter is empty")
//	}
func Empty(u *url.URL, key string) bool {
	return u.Query().Get(key) == ""
}
