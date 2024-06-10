package qp

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

// ParseBool parses a boolean query parameter from the given URL.
//
// The function accepts a URL, a key, and an optional default boolean value.
// If the query parameter is absent or empty, the default value (false)
// is returned.
//
// If one boolean is provided, it is used as the default value. If the
// query parameter is absent or empty, this default value is returned.
//
// The function supports the following values as valid booleans:
//   - true/false
//   - yes/no
//   - on/off
//   - 1/0
//
// Example Usage:
//
//	// Simple call without default value.
//	result := ParseBool(u, "active")
//
//	// Call with default value.
//	// Default: true
//	result := ParseBool(u, "enabled", true)
func ParseBool(u *url.URL, key string, opt ...bool) *Result[bool] {
	result := &Result[bool]{Key: key, Contains: true}
	data, ok := u.Query()[key]

	// Default value.
	if len(opt) >= 1 {
		result.Default = opt[0]
		result.Value = result.Default
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

	// Convert the result to an integer.
	raw := strings.ToLower(data[0])
	if raw == "1" || raw == "0" {
		result.Value = raw == "1"
	} else if raw == "yes" || raw == "no" {
		result.Value = raw == "yes"
	} else if raw == "on" || raw == "off" {
		result.Value = raw == "on"
	} else {
		value, err := strconv.ParseBool(data[0])
		if err != nil {
			msg := "invalid value for key %s: %s"
			result.Error = fmt.Errorf(msg, key, data[0])
			return result
		}

		result.Value = value
	}

	return result
}

// GetBool is the function to parse a boolean query parameter and return
// the value and a boolean indicating if the value is valid.
//
// The function accepts a URL, a key, and an optional default boolean value.
// If the query parameter is absent or empty, the default value (false)
// is returned.
//
// If one boolean is provided, it is used as the default value. If the
// query parameter is absent or empty, this default value is returned.
//
// The function supports the following values as valid booleans:
//   - true/false
//   - yes/no
//   - on/off
//   - 1/0
//
// Example Usage:
//
//	// Simple call without default value.
//	result, ok := GetBool(u, "active")
//
//	// Call with default value.
//	// Default: true
//	result, ok := GetBool(u, "active", true)
func GetBool(u *url.URL, key string, opt ...bool) (bool, bool) {
	data := ParseBool(u, key, opt...)
	return data.Value, data.Contains && !data.Empty && data.Error == nil
}

// ParseBoolSlice parses a boolean slice query parameter from the given URL.
//
// The function accepts a URL and a key. If the query parameter is absent,
// nil is returned. If the query parameter is present but empty, an empty
// slice is returned.
//
// The function supports the following values as valid booleans:
//   - true/false
//   - yes/no
//   - on/off
//   - 1/0
//
// The function supports query parameters specified as a single string
// (e.g., "?flags=true,false,yes,no") or as multiple values (e.g., "?flags=true&flags=false&flags=yes&flags=no").
//
// Example Usage:
//
//	// Simple call.
//	result := ParseBoolSlice(u, "flags")
//
//	// Handling the result.
//	if result.Contains && !result.Empty && result.Error == nil {
//	    fmt.Println("Parsed booleans:", result.Value)
//	} else if result.Empty {
//	    fmt.Println("Query parameter is empty.")
//	} else if result.Error != nil {
//	    fmt.Println("Error parsing query parameter:", result.Error)
//	} else {
//	    fmt.Println("Query parameter is absent.")
//	}
func ParseBoolSlice(u *url.URL, key string) *Result[[]bool] {
	result := &Result[[]bool]{Key: key, Contains: true}
	data, ok := u.Query()[key]

	if !ok {
		result.Value = nil
		result.Empty = true
		result.Contains = false
		return result
	} else if data[0] == "" {
		result.Value = []bool{} // not nil
		result.Empty = true
		result.Contains = true
		return result
	}

	// An array can be specified as a single string "?flags=true,false,yes,no" or
	// as multiple values "?flags=true&flags=false&flags=yes&flags=no".
	if len(data) > 1 {
		// Multiple values.
		result.Value = make([]bool, 0, len(data))
		for _, str := range data {
			value, err := parseBoolValue(str)
			if err != nil {
				msg := "invalid value for key %s: %s"
				result.Error = fmt.Errorf(msg, key, str)
				result.Value = []bool{} // not nil
				return result
			}
			result.Value = append(result.Value, value)
		}
		return result
	}

	// Single value.
	result.Value = make([]bool, 0)
	for _, str := range strings.Split(data[0], ",") {
		value, err := parseBoolValue(str)
		if err != nil {
			msg := "invalid value for key %s: %s"
			result.Error = fmt.Errorf(msg, key, str)
			result.Value = []bool{} // not nil
			return result
		}
		result.Value = append(result.Value, value)
	}

	return result
}

// parseBoolValue parses a string and returns its boolean value if valid.
func parseBoolValue(str string) (bool, error) {
	raw := strings.ToLower(str)
	switch raw {
	case "1", "true", "yes", "on":
		return true, nil
	case "0", "false", "no", "off":
		return false, nil
	default:
		return strconv.ParseBool(raw)
	}
}

// GetBoolSlice is a convenience function to parse a boolean slice query
// parameter and return the slice of values.
//
// The function accepts a URL and a key. If the query parameter is absent,
// nil is returned. If the query parameter is present but empty, an empty
// slice is returned.
//
// Example Usage:
//
//	// Simple call.
//	values := GetBoolSlice(u, "flags")
//
//	// Handling the result.
//	if values == nil {
//	    fmt.Println("Query parameter is absent.")
//	} else if len(values) == 0 {
//	    fmt.Println("Query parameter is empty.")
//	} else {
//	    fmt.Println("Parsed booleans:", values)
//	}
func GetBoolSlice(u *url.URL, key string) []bool {
	return ParseBoolSlice(u, key).Value
}
