// Package qp provides utilities for parsing query parameters from URLs.
//
// The package offers robust and user-friendly functions for parsing query
// parameters into various types (integers, floats, strings, booleans) and their
// slice variants. Each parsing function returns detailed results including the
// parsed value, validation status, and error information if any.
//
// # Types of Methods
//
// The package provides three types of methods for each data type:
//
//   - Parse methods: Return detailed results in a Result structure, including
//     technical data such as errors, presence flags, and parsed values.
//
//   - Get methods: Return a value and a boolean indicating validity. The value
//     is never nil (even for slices), and the boolean is true only if the
//     parameter exists and was parsed correctly.
//
//   - Pull methods: Return a pointer that is nil if the parameter is absent,
//     points to the default value if the parameter is present but invalid,
//     or points to the parsed value if successful.
//
// # Result Structure
//
//	type Result[T Value] struct {
//	    Key      string // Parameter name
//	    Value    T      // Parsed value
//	    Default  T      // Default value
//	    Min      T      // Minimum value (for numeric types)
//	    Max      T      // Maximum value (for numeric types)
//	    Others   []T    // Additional valid values
//	    Empty    bool   // True if parameter is empty
//	    Contains bool   // True if parameter exists
//	    Error    error  // Parsing error, if any
//	}
//
// # Examples
//
// Using Pull methods with SQL WHERE clause generation:
//
//	func where(isActive, isStaff, isSuperuser *bool) string {
//	    check := [...]struct {
//	        name  string
//	        value *bool
//	    }{
//	        {name: "is_active", value: isActive},
//	        {name: "is_staff", value: isStaff},
//	        {name: "is_superuser", value: isSuperuser},
//	    }
//
//	    and := make([]string, 0, len(check))
//	    for _, m := range check {
//	        if m.value != nil {
//	            and = append(and, fmt.Sprintf("%s=%t", m.name, *m.value))
//	        }
//	    }
//
//	    if len(and) != 0 {
//	        return "WHERE " + strings.Join(and, ", ")
//	    }
//
//	    return ""
//	}
//
//	u, _ := url.Parse("http://example.com?is_active=true&is_staff=false")
//	isActive := qp.PullBool(u, "is_active")
//	isStaff := qp.PullBool(u, "is_staff")
//	isSuperuser := qp.PullBool(u, "is_superuser") // nil
//	query := where(isActive, isStaff, isSuperuser)
//	// Result: WHERE is_active=true, is_staff=false
//
// # Integer Parsing
//
// Parse a single integer:
//
//	u, _ := url.Parse("http://example.com?age=18")
//
//	// Using Parse method
//	result := qp.ParseInt(u, "age")
//	if result.Contains && !result.Empty && result.Error == nil {
//	    fmt.Println("Age:", result.Value)
//	}
//
//	// Using Get method
//	value, ok := qp.GetInt(u, "age")
//
//	// Using Pull method
//	age := qp.PullInt(u, "age")
//
// Parse integer slice:
//
//	u, _ := url.Parse("http://example.com?ids=1,2,3")
//	result := qp.ParseIntSlice(u, "ids")
//	ids, ok := qp.GetIntSlice(u, "ids")
//	ids = qp.PullIntSlice(u, "ids")
//
// # Float Parsing
//
// Parse a single float:
//
//	u, _ := url.Parse("http://example.com?temperature=36.6")
//
//	// Using Parse method
//	result := qp.ParseFloat(u, "temperature")
//	if result.Contains && !result.Empty && result.Error == nil {
//	    fmt.Println("Temperature:", result.Value)
//	}
//
//	// Using Get method
//	value, ok := qp.GetFloat(u, "temperature")
//
//	// Using Pull method
//	temp := qp.PullFloat(u, "temperature")
//
// Parse float slice:
//
//	u, _ := url.Parse("http://example.com?temps=36.6,37.2,36.9")
//	result := qp.ParseFloatSlice(u, "temps")
//	temps, ok := qp.GetFloatSlice(u, "temps")
//	temps = qp.PullFloatSlice(u, "temps")
//
// # String Parsing
//
// Parse a single string:
//
//	u, _ := url.Parse("http://example.com?name=alice")
//
//	// Using Parse method
//	result := qp.ParseString(u, "name")
//	if result.Contains && !result.Empty && result.Error == nil {
//	    fmt.Println("Name:", result.Value)
//	}
//
//	// Using Get method
//	value, ok := qp.GetString(u, "name")
//
//	// Using Pull method
//	name := qp.PullString(u, "name")
//
// Parse string slice:
//
//	u, _ := url.Parse("http://example.com?names=alice,bob,charlie")
//	result := qp.ParseStringSlice(u, "names")
//	names, ok := qp.GetStringSlice(u, "names")
//	names = qp.PullStringSlice(u, "names")
//
// # Boolean Parsing
//
// Parse a single boolean:
//
//	u, _ := url.Parse("http://example.com?active=true")
//
//	// Using Parse method
//	result := qp.ParseBool(u, "active")
//	if result.Contains && !result.Empty && result.Error == nil {
//	    fmt.Println("Active:", result.Value)
//	}
//
//	// Using Get method
//	value, ok := qp.GetBool(u, "active")
//
//	// Using Pull method
//	active := qp.PullBool(u, "active")
//
// Parse boolean slice:
//
//	u, _ := url.Parse("http://example.com?flags=true,false,yes,no")
//	result := qp.ParseBoolSlice(u, "flags")
//	flags, ok := qp.GetBoolSlice(u, "flags")
//	flags = qp.PullBoolSlice(u, "flags")
//
// # Utility Functions
//
// Check parameter presence:
//
//	u, _ := url.Parse("http://example.com?age=18")
//	if qp.Contains(u, "age") {
//	    fmt.Println("Age parameter exists")
//	}
//
// Check if parameter is empty:
//
//	u, _ := url.Parse("http://example.com?age=")
//	if qp.Empty(u, "age") {
//	    fmt.Println("Age parameter is empty")
//	}
//
// # Notes
//
//   - Boolean parsing supports multiple formats: true/false, yes/no, on/off, 1/0
//   - Slice parameters can be specified either as comma-separated values
//     (?ids=1,2,3) or as multiple parameters (?ids=1&ids=2&ids=3)
//   - All numeric parsers support range validation and additional valid values
//   - String parsers support validation against a list of valid values
package qp
