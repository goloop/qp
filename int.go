package qp

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/goloop/g"
)

// ParseInt parses an integer query parameter from the given URL.
//
// The function accepts a URL, a key, and an optional list of integers.
// The optional integers can specify a default value, a range (min and max),
// and additional valid values.
//
// If no optional integers are provided, the function simply attempts to
// parse the integer value from the query parameter. If the query parameter
// is absent or empty, the default value (zero) is returned.
//
// If one integer is provided, it is used as the default value. If the
// query parameter is absent or empty, this default value is returned.
//
// If two integers are provided, they specify a range (min and max), where
// the first integer is the default value. If the query parameter value
// is outside this range, the default value is returned.
//
// If more than two integers are provided, the first two integers specify
// the default value and the range (min and max), and any additional
// integers are treated as additional valid values. If the query parameter
// value is not within the range or among the additional valid values,
// the default value is returned.
//
// Example Usage:
//
//	// Simple call without default, min, max, or others.
//	result := ParseInt(u, "age")
//
//	// Call with default value.
//	// Default: 18
//	result := ParseInt(u, "age", 18)
//
//	// Call with default and min-max.
//	// Default: 18
//	// Range:   18-30
//	result := ParseInt(u, "age", 18, 30)
//
//	// Call with default and max-min (reversed order).
//	// Default: 30
//	// Range:   18-30
//	result := ParseInt(u, "age", 30, 18)
//
//	// Call with default, min-max, and additional valid values.
//	// Default:    18
//	// Range:      18-30
//	// Additional: 20, 25, 30
//	result := ParseInt(u, "age", 18, 30, 20, 25, 35)
//
//	// Call with default and additional valid values without min-max.
//	// Default:    10
//	// Additional: 10, 20, 30
//	result := ParseInt(u, "age", 10, 10, 20, 30)
func ParseInt(u *url.URL, key string, opt ...int) *Result[int] {
	result := &Result[int]{Key: key, Contains: true}
	data, ok := u.Query()[key]

	// Available values.
	if len(opt) == 1 {
		// Default value.
		result.Default = opt[0]
		result.Value = result.Default
	} else if len(opt) > 1 {
		// Range and default value.
		min, max := opt[0], opt[1]
		if min > max {
			min, max = max, min
		}

		result.Min = min
		result.Max = max
		result.Default = opt[0] // not min or max, but first value
		result.Value = result.Default

		// Set additional valid values.
		if len(opt) > 2 {
			result.Others = make([]int, 0, len(opt)-2)
			result.Others = append(result.Others, opt[2:]...)
		}
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
	value, err := strconv.Atoi(data[0])
	if err != nil {
		result.Error = fmt.Errorf("invalid value for key %s: %s", key, data[0])
		return result
	}

	if len(opt) < 2 {
		// No range or any available values.
		result.Value = value
	} else if value >= result.Min && value <= result.Max {
		// Check if the value is within the specified range.
		result.Value = value
	} else {
		// Check if the value is in the list of available values.
		if result.Others != nil && g.In(value, result.Others...) {
			result.Value = value
		} else {
			msg := "value out of range for key %s: %d"
			result.Error = fmt.Errorf(msg, key, value)
		}
	}

	return result
}

// GetInt parses an integer query parameter and return the value and a
// boolean indicating, true - if a value was passed in query params and
// successfully parsed.
//
// The function accepts a URL, a key, and an optional list of integers.
// The optional integers can specify a default value, a range (min and max),
// and additional valid values.
//
// If no optional integers are provided, the function simply attempts to
// parse the integer value from the query parameter. If the query parameter
// is absent or empty, the default value (zero) is returned.
//
// If one integer is provided, it is used as the default value. If the
// query parameter is absent or empty, this default value is returned.
//
// If two integers are provided, they specify a range (min and max), where
// the first integer is the default value. If the query parameter value
// is outside this range, the default value is returned.
//
// If more than two integers are provided, the first two integers specify
// the default value and the range (min and max), and any additional
// integers are treated as additional valid values. If the query parameter
// value is not within the range or among the additional valid values,
// the default value is returned.
//
// Example Usage:
//
//	// Simple call without default, min, max, or others.
//	result, ok := GetInt(u, "age")
//
//	// Call with default value.
//	// Default: 18
//	result, ok := GetInt(u, "age", 18)
//
//	// Call with default and min-max.
//	// Default: 18
//	// Range:   18-30
//	result, ok := GetInt(u, "age", 18, 30)
//
//	// Call with default and max-min (reversed order).
//	// Default: 30
//	// Range:   18-30
//	result, ok := GetInt(u, "age", 30, 18)
//
//	// Call with default, min-max, and additional valid values.
//	// Default:    18
//	// Range:      18-30
//	// Additional: 20, 25, 30
//	result, ok := GetInt(u, "age", 18, 30, 20, 25, 35)
//
//	// Call with default and additional valid values without min-max.
//	// Default:    10
//	// Additional: 10, 20, 30
//	result, ok := GetInt(u, "age", 10, 10, 20, 30)
func GetInt(u *url.URL, key string, opt ...int) (int, bool) {
	data := ParseInt(u, key, opt...)
	return data.Value, data.Contains && !data.Empty && data.Error == nil
}

// PullInt returns a pointer to the parsed integer query parameter value.
//
// The function accepts a URL, a key, and an optional list of integers.
// The optional integers can specify a default value, a range (min and max),
// and additional valid values.
//
// If no optional integers are provided, the function simply attempts to
// parse the integer value from the query parameter. If the query parameter
// is absent, nil is returned.
//
// If the parameter is specified, but it is empty or invalid, a pointer to
// the empty type is returned or default value.
//
// If one integer is provided, it is used as the default value. If the
// query parameter is absent or empty, this default value is returned
// as a pointer.
//
// If two integers are provided, they specify a range (min and max), where
// the first integer is the default value. If the query parameter value
// is outside this range, the default value is returned as a pointer.
//
// If more than two integers are provided, the first two integers specify
// the default value and the range (min and max), and any additional
// integers are treated as additional valid values. If the query parameter
// value is not within the range or among the additional valid values,
// the default value is returned as a pointer.
//
// Example Usage:
//
//	// Simple call without default, min, max, or others.
//	result := PullInt(u, "age")
//
//	// Call with default value.
//	// Default: 18
//	result := PullInt(u, "age", 18)
//
//	// Call with default and min-max.
//	// Default: 18
//	// Range:   18-30
//	result := PullInt(u, "age", 18, 30)
//
//	// Call with default and max-min (reversed order).
//	// Default: 30
//	// Range:   18-30
//	result := PullInt(u, "age", 30, 18)
//
//	// Call with default, min-max, and additional valid values.
//	// Default:    18
//	// Range:      18-30
//	// Additional: 20, 25, 35
//	result := PullInt(u, "age", 18, 30, 20, 25, 35)
//
//	// Call with default and additional valid values without min-max.
//	// Default:    10
//	// Additional: 10, 20, 30
//	result := PullInt(u, "age", 10, 10, 20, 30)
func PullInt(u *url.URL, key string, opt ...int) *int {
	data := ParseInt(u, key, opt...)
	if !data.Contains {
		return nil
	}

	return &data.Value
}

// ParseIntSlice parses an integer slice query parameter from the given URL.
//
// The function accepts a URL and a key. If the query parameter is absent,
// nil is returned. If the query parameter is present but empty, an empty
// slice is returned.
//
// The function supports query parameters specified as a single string
// (e.g., "?ids=1,2,3") or as multiple values (e.g., "?ids=1&ids=2&ids=3").
//
// Example Usage:
//
//	// Simple call.
//	result := ParseIntSlice(u, "ids")
//
//	// Handling the result.
//	if result.Contains && !result.Empty && result.Error == nil {
//	    fmt.Println("Parsed integers:", result.Value)
//	} else if result.Empty {
//	    fmt.Println("Query parameter is empty.")
//	} else if result.Error != nil {
//	    fmt.Println("Error parsing query parameter:", result.Error)
//	} else {
//	    fmt.Println("Query parameter is absent.")
//	}
func ParseIntSlice(u *url.URL, key string, opt ...[]int) *Result[[]int] {
	result := &Result[[]int]{Key: key, Contains: true}
	data, ok := u.Query()[key]

	// Default value.
	result.Default = []int{} // not nil
	result.Value = result.Default
	if len(opt) > 0 {
		result.Default = opt[0]
		result.Value = result.Default
	}

	// Check if the query parameter is empty or missing.
	if !ok {
		result.Empty = true
		result.Contains = false
		return result
	} else if data[0] == "" {
		result.Empty = true
		result.Contains = true
		return result
	}

	// An array can be specified as a single string "?ids=1,2,3" or
	// as multiple values "?ids=1&ids=2&ids=3".
	if len(data) > 1 {
		// Multiple values.
		result.Value = make([]int, 0, len(data))
		for _, str := range data {
			value, err := strconv.Atoi(str)
			if err != nil {
				msg := "invalid value for key %s: %s"
				result.Error = fmt.Errorf(msg, key, str)
				result.Value = []int{} // not nil
				return result
			}
			result.Value = append(result.Value, value)
		}
		return result
	}

	// Single value.
	result.Value = make([]int, 0)
	for _, str := range strings.Split(data[0], ",") {
		value, err := strconv.Atoi(str)
		if err != nil {
			msg := "invalid value for key %s: %s"
			result.Error = fmt.Errorf(msg, key, str)
			result.Value = []int{} // not nil
			return result
		}
		result.Value = append(result.Value, value)
	}

	return result
}

// GetIntSlice parses an integer slice query parameter from the given URL
// and returns the slice of values and a boolean indicating if the
// value is valid.
//
// The function accepts a URL and a key. If the query parameter is absent,
// nil is returned. If the query parameter is present but empty, an empty
// slice is returned.
//
// The function supports query parameters specified as a single string
// (e.g., "?ids=1,2,3") or as multiple values (e.g., "?ids=1&ids=2&ids=3").
//
// Example Usage:
//
//	// Simple call.
//	result, ok := GetIntSlice(u, "ids")
//
//	// Handling the result.
//	if ok {
//	    fmt.Println("Parsed integers:", result)
//	} else {
//	    fmt.Println("Query parameter is absent or invalid.")
//	}
//
//	// Call with default value.
//	// Default: []int{1, 2, 3}
//	result, ok := GetIntSlice(u, "ids", []int{1, 2, 3})
func GetIntSlice(u *url.URL, key string, opt ...[]int) ([]int, bool) {
	data := ParseIntSlice(u, key, opt...)
	return data.Value, data.Contains && !data.Empty && data.Error == nil
}

// PullIntSlice parses an integer slice query parameter
// and return the slice of values.
//
// The function accepts a URL and a key. If the query parameter is absent,
// nil is returned. Returns nil even if the default value is set.
//
// If the query parameter is present but empty, an empty slice is returned
// or the default value if it is set.
//
// Example Usage:
//
//	// Simple call.
//	values := PullIntSlice(u, "ids")
//
//	// Handling the result.
//	if values == nil {
//	    fmt.Println("Query parameter is absent.")
//	} else if len(values) == 0 {
//	    fmt.Println("Query parameter is empty.")
//	} else {
//	    fmt.Println("Parsed integers:", values)
//	}
func PullIntSlice(u *url.URL, key string, opt ...[]int) []int {
	data := ParseIntSlice(u, key, opt...)
	if !data.Contains {
		return nil // not default
	}

	return data.Value
}
