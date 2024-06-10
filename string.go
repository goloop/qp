package qp

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/goloop/g"
)

// ParseString parses a string query parameter from the given URL.
//
// The function accepts a URL, a key, and an optional list of strings.
// The optional strings can specify a default value and a list of valid values.
//
// If no optional strings are provided, the function simply attempts to
// parse the string value from the query parameter. If the query parameter
// is absent or empty, the default value (an empty string) is returned.
//
// If one string is provided, it is used as the default value. If the
// query parameter is absent or empty, this default value is returned.
//
// If more than one string is provided, the first string is used as the
// default value.
// Any additional strings specify valid values. If the query parameter
// value is not among the valid values, the default value is returned.
//
// Example Usage:
//
//	// Simple call without default value or valid values.
//	result := ParseString(u, "name")
//
//	// Call with default value.
//	// Default: "guest"
//	result := ParseString(u, "name", "guest")
//
//	// Call with default value and valid values.
//	// Default: "guest"
//	// Valid values: "guest", "admin", "user"
//	result := ParseString(u, "name", "guest", "admin", "user")
func ParseString(u *url.URL, key string, opt ...string) *Result[string] {
	result := &Result[string]{Key: key, Contains: true}
	data, ok := u.Query()[key]

	// Available values.
	if len(opt) == 1 {
		// Default value.
		result.Default = opt[0]
		result.Value = result.Default
	} else if len(opt) > 1 {
		result.Default = opt[0]
		result.Value = result.Default

		// Additional valid values.
		// Default value is part of the valid values.
		result.Others = make([]string, 0, len(opt))
		result.Others = append(result.Others, opt...)
	}

	// Check if the query parameter is empty or missing.
	if !ok {
		// The query parameter is missing.
		result.Empty = true
		result.Contains = false
		return result
	} else if data[0] == "" {
		// The query parameter is empty.
		result.Empty = true
		result.Contains = true
		return result
	}

	// Check if the value is in the list of valid values.
	value := data[0]
	if len(opt) < 2 {
		result.Value = value
	} else {
		if result.Others != nil && g.In(value, result.Others...) {
			result.Value = value
		} else {
			msg := "value out of range for key %s: %d"
			result.Error = fmt.Errorf(msg, key, value)
		}
	}

	return result
}

// GetString is the function to parse a string query parameter
// and return the value and a boolean indicating if the value is valid.
//
// The function accepts a URL, a key, and an optional list of strings.
// The optional strings can specify a default value and a list of valid values.
//
// If no optional strings are provided, the function simply attempts to
// parse the string value from the query parameter. If the query parameter
// is absent or empty, the default value (an empty string) is returned.
//
// If one string is provided, it is used as the default value. If the
// query parameter is absent or empty, this default value is returned.
//
// If more than one string is provided, the first string is used as the
// default value.
// Any additional strings specify valid values. If the query parameter
// value is not among the valid values, the default value is returned.
//
// Example Usage:
//
//	// Simple call without default value or valid values.
//	result, ok := GetString(u, "name")
//
//	// Call with default value.
//	// Default: "guest"
//	result, ok := GetString(u, "name", "guest")
//
//	// Call with default value and valid values.
//	// Default: "guest"
//	// Valid values: "guest", "admin", "user"
//	result, ok := GetString(u, "name", "guest", "admin", "user")
func GetString(u *url.URL, key string, opt ...string) (string, bool) {
	data := ParseString(u, key, opt...)
	return data.Value, data.Contains && !data.Empty && data.Error == nil
}

// ParseStringSlice parses a string slice query parameter from the given URL.
//
// The function accepts a URL and a key. If the query parameter is absent,
// nil is returned. If the query parameter is present but empty, an empty
// slice is returned.
//
// The function supports query parameters specified as a single string
// (e.g., "?names=alice,bob,charlie") or as multiple values (e.g., "?names=alice&names=bob&names=charlie").
//
// Example Usage:
//
//	// Simple call.
//	result := ParseStringSlice(u, "names")
//
//	// Handling the result.
//	if result.Contains && !result.Empty && result.Error == nil {
//	    fmt.Println("Parsed strings:", result.Value)
//	} else if result.Empty {
//	    fmt.Println("Query parameter is empty.")
//	} else if result.Error != nil {
//	    fmt.Println("Error parsing query parameter:", result.Error)
//	} else {
//	    fmt.Println("Query parameter is absent.")
//	}
func ParseStringSlice(u *url.URL, key string) *Result[[]string] {
	result := &Result[[]string]{Key: key, Contains: true}
	data, ok := u.Query()[key]

	if !ok {
		result.Value = nil
		result.Empty = true
		result.Contains = false
		return result
	} else if data[0] == "" {
		result.Value = []string{} // not nil
		result.Empty = true
		result.Contains = true
		return result
	}

	// An array can be specified as a single string "?names=alice,bob,charlie" or
	// as multiple values "?names=alice&names=bob&names=charlie".
	if len(data) > 1 {
		// Multiple values.
		result.Value = make([]string, 0, len(data))
		for _, str := range data {
			result.Value = append(result.Value, str)
		}
		return result
	}

	// Single value.
	result.Value = strings.Split(data[0], ",")
	return result
}

// GetStringSlice is a convenience function to parse a string slice query
// parameter and return the slice of values.
//
// The function accepts a URL and a key. If the query parameter is absent,
// nil is returned. If the query parameter is present but empty, an empty
// slice is returned.
//
// Example Usage:
//
//	// Simple call.
//	values := GetStringSlice(u, "names")
//
//	// Handling the result.
//	if values == nil {
//	    fmt.Println("Query parameter is absent.")
//	} else if len(values) == 0 {
//	    fmt.Println("Query parameter is empty.")
//	} else {
//	    fmt.Println("Parsed strings:", values)
//	}
func GetStringSlice(u *url.URL, key string) []string {
	return ParseStringSlice(u, key).Value
}
