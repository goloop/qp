package qp

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/goloop/g"
)

// ParseFloat parses a float query parameter from the given URL.
//
// The function accepts a URL, a key, and an optional list of floats.
// The optional floats can specify a default value, a range (min and max),
// and additional valid values.
//
// If no optional floats are provided, the function simply attempts to
// parse the float value from the query parameter. If the query parameter
// is absent or empty, the default value (zero) is returned.
//
// If one float is provided, it is used as the default value. If the
// query parameter is absent or empty, this default value is returned.
//
// If two floats are provided, they specify a range (min and max), where
// the first float is the default value. If the query parameter value
// is outside this range, the default value is returned.
//
// If more than two floats are provided, the first two floats specify
// the default value and the range (min and max), and any additional
// floats are treated as additional valid values. If the query parameter
// value is not within the range or among the additional valid values,
// the default value is returned.
//
// Example Usage:
//
//	// Simple call without default, min, max, or others.
//	result := ParseFloat(u, "temperature")
//
//	// Call with default value.
//	// Default: 18.5
//	result := ParseFloat(u, "temperature", 18.5)
//
//	// Call with default and min-max.
//	// Default: 18.5
//	// Range:   18.5-30.0
//	result := ParseFloat(u, "temperature", 18.5, 30.0)
//
//	// Call with default and max-min (reversed order).
//	// Default: 30.0
//	// Range:   18.5-30.0
//	result := ParseFloat(u, "temperature", 30.0, 18.5)
//
//	// Call with default, min-max, and additional valid values.
//	// Default:    18.5
//	// Range:      18.5-30.0
//	// Additional: 20.0, 25.0, 30.0
//	result := ParseFloat(u, "temperature", 18.5, 30.0, 20.0, 25.0, 35.0)
//
//	// Call with default and additional valid values without min-max.
//	// Default:    10.5
//	// Additional: 10.5, 20.0, 30.0
//	result := ParseFloat(u, "temperature", 10.5, 10.5, 20.0, 30.0)
func ParseFloat(u *url.URL, key string, opt ...float64) *Result[float64] {
	result := &Result[float64]{Key: key, Contains: true}
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
			result.Others = make([]float64, 0, len(opt)-2)
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

	// Convert the result to a float.
	value, err := strconv.ParseFloat(data[0], 64)
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
			msg := "value out of range for key %s: %f"
			result.Error = fmt.Errorf(msg, key, value)
		}
	}

	return result
}

// GetFloat is the function to parse a float query parameter and return
// the value and a boolean indicating if the value is valid.
//
// The function accepts a URL, a key, and an optional list of floats.
// The optional floats can specify a default value, a range (min and max),
// and additional valid values.
//
// If no optional floats are provided, the function simply attempts to
// parse the float value from the query parameter. If the query parameter
// is absent or empty, the default value (zero) is returned.
//
// If one float is provided, it is used as the default value. If the
// query parameter is absent or empty, this default value is returned.
//
// If two floats are provided, they specify a range (min and max), where
// the first float is the default value. If the query parameter value
// is outside this range, the default value is returned.
//
// If more than two floats are provided, the first two floats specify
// the default value and the range (min and max), and any additional
// floats are treated as additional valid values. If the query parameter
// value is not within the range or among the additional valid values,
// the default value is returned.
//
// Example Usage:
//
//	// Simple call without default, min, max, or others.
//	result, ok := GetFloat(u, "temperature")
//
//	// Call with default value.
//	// Default: 18.5
//	result, ok := GetFloat(u, "temperature", 18.5)
//
//	// Call with default and min-max.
//	// Default: 18.5
//	// Range:   18.5-30.0
//	result, ok := GetFloat(u, "temperature", 18.5, 30.0)
//
//	// Call with default and max-min (reversed order).
//	// Default: 30.0
//	// Range:   18.5-30.0
//	result, ok := GetFloat(u, "temperature", 30.0, 18.5)
//
//	// Call with default, min-max, and additional valid values.
//	// Default:    18.5
//	// Range:      18.5-30.0
//	// Additional: 20.0, 25.0, 30.0
//	result, ok := GetFloat(u, "temperature", 18.5, 30.0, 20.0, 25.0, 35.0)
//
//	// Call with default and additional valid values without min-max.
//	// Default:    10.5
//	// Additional: 10.5, 20.0, 30.0
//	result, ok := GetFloat(u, "temperature", 10.5, 10.5, 20.0, 30.0)
func GetFloat(u *url.URL, key string, opt ...float64) (float64, bool) {
	data := ParseFloat(u, key, opt...)
	return data.Value, data.Contains && !data.Empty && data.Error == nil
}

// PullFloat returns a pointer to the parsed float64 query parameter value.
//
// The function accepts a URL, a key, and an optional list of floats.
// The optional floats can specify a default value, a range (min and max),
// and additional valid values.
//
// If no optional floats are provided, the function simply attempts to
// parse the integer value from the query parameter. If the query parameter
// is absent, nil is returned.
//
// If one integer is provided, it is used as the default value. If the
// query parameter is absent or empty, this default value is returned
// as a pointer.
//
// If two floats are provided, they specify a range (min and max), where
// the first integer is the default value. If the query parameter value
// is outside this range, the default value is returned as a pointer.
//
// If more than two floats are provided, the first two floats specify
// the default value and the range (min and max), and any additional
// floats are treated as additional valid values. If the query parameter
// value is not within the range or among the additional valid values,
// the default value is returned as a pointer.
//
// Example Usage:
//
//	// Simple call without default, min, max, or others.
//	result := PullFloat(u, "age")
//
//	// Call with default value.
//	// Default: 18.0
//	result := PullFloat(u, "age", 18.0)
//
//	// Call with default and min-max.
//	// Default: 18.0
//	// Range:   18.0-30.0
//	result := PullFloat(u, "age", 18.0, 30.0)
//
//	// Call with default and max-min (reversed order).
//	// Default: 30.0
//	// Range:   18.0-30.0
//	result := PullFloat(u, "age", 30.0, 18.0)
//
//	// Call with default, min-max, and additional valid values.
//	// Default:    18.0
//	// Range:      18.0-30.0
//	// Additional: 20.0, 25.0, 30.0
//	result := PullFloat(u, "age", 18.0, 30.0, 20.0, 25.0, 35.0)
//
//	// Call with default and additional valid values without min-max.
//	// Default:    10.0
//	// Additional: 10.0, 20.0, 30.0
//	result := PullFloat(u, "age", 10.0, 10.0, 20.0, 30.0)
func PullFloat(u *url.URL, key string, opt ...float64) *float64 {
	data := ParseFloat(u, key, opt...)
	if !data.Contains {
		return nil
	}

	return &data.Value
}

// ParseFloatSlice parses a float64 slice query parameter from the given URL.
//
// The function accepts a URL and a key. If the query parameter is absent,
// nil is returned. If the query parameter is present but empty, an empty
// slice is returned.
//
// The function supports query parameters specified as a single string
// (e.g., "?values=1.1,2.2,3.3") or as multiple values (e.g.,
// "?values=1.1&values=2.2&values=3.3").
//
// Example Usage:
//
//	// Simple call.
//	result := ParseFloatSlice(u, "values")
//
//	// Handling the result.
//	if result.Contains && !result.Empty && result.Error == nil {
//	    fmt.Println("Parsed floats:", result.Value)
//	} else if result.Empty {
//	    fmt.Println("Query parameter is empty.")
//	} else if result.Error != nil {
//	    fmt.Println("Error parsing query parameter:", result.Error)
//	} else {
//	    fmt.Println("Query parameter is absent.")
//	}
func ParseFloatSlice(
	u *url.URL,
	key string,
	opt ...[]float64,
) *Result[[]float64] {
	result := &Result[[]float64]{Key: key, Contains: true}
	data, ok := u.Query()[key]

	// Default value.
	result.Default = []float64{} // not nil
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

	// An array can be specified as a single string "?values=1.1,2.2,3.3" or
	// as multiple values "?values=1.1&values=2.2&values=3.3".
	if len(data) > 1 {
		// Multiple values.
		result.Value = make([]float64, 0, len(data))
		for _, str := range data {
			value, err := strconv.ParseFloat(str, 64)
			if err != nil {
				msg := "invalid value for key %s: %s"
				result.Error = fmt.Errorf(msg, key, str)
				result.Value = []float64{} // not nil
				return result
			}
			result.Value = append(result.Value, value)
		}
		return result
	}

	// Single value.
	result.Value = make([]float64, 0)
	for _, str := range strings.Split(data[0], ",") {
		value, err := strconv.ParseFloat(str, 64)
		if err != nil {
			msg := "invalid value for key %s: %s"
			result.Error = fmt.Errorf(msg, key, str)
			result.Value = []float64{} // not nil
			return result
		}
		result.Value = append(result.Value, value)
	}

	return result
}

// GetFloatSlice parses an float64 slice query parameter from the given URL
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
//	result, ok := GetFloatSlice(u, "ids")
//
//	// Handling the result.
//	if ok {
//	    fmt.Println("Parsed floats:", result)
//	} else {
//	    fmt.Println("Query parameter is absent or invalid.")
//	}
//
//	// Call with default value.
//	// Default: []int{1, 2, 3}
//	result, ok := GetFloatSlice(u, "ids", []int{1, 2, 3})
func GetFloatSlice(u *url.URL, key string, opt ...[]float64) ([]float64, bool) {
	data := ParseFloatSlice(u, key, opt...)
	return data.Value, data.Contains && !data.Empty && data.Error == nil
}

// PullFloatSlice is a convenience function to parse a float64 slice query
// parameter and return the slice of values.
//
// The function accepts a URL and a key. If the query parameter is absent,
// nil is returned. If the query parameter is present but empty, an empty
// slice is returned.
//
// Example Usage:
//
//	// Simple call.
//	values := PullFloatSlice(u, "values")
//
//	// Handling the result.
//	if values == nil {
//	    fmt.Println("Query parameter is absent.")
//	} else if len(values) == 0 {
//	    fmt.Println("Query parameter is empty.")
//	} else {
//	    fmt.Println("Parsed floats:", values)
//	}
func PullFloatSlice(u *url.URL, key string, opt ...[]float64) []float64 {
	data := ParseFloatSlice(u, key, opt...)
	if !data.Contains {
		return nil // not default
	}

	return data.Value
}
