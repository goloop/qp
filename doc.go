/*
Package qp provides utilities for parsing query parameters from URLs.

This package includes functions for parsing various types of query parameters,
such as integers, floats, strings, and booleans. It supports parsing both single
ids and slices of ids. Additionally, it provides convenient functions for
checking the presence and emptiness of query parameters.

The functions in this package are designed to be robust and user-friendly,
returning detailed results that include the parsed value, default value, and
information about whether the query parameter was present and valid.

# Example Usage

Here are some examples of how to use the functions provided by this package:

# Parsing Integer Query Parameters

Parse a single integer query parameter:

	u, _ := url.Parse("http://example.com?age=18")
	result := qp.ParseInt(u, "age")
	if result.Contains && !result.Empty && result.Error == nil {
	    fmt.Println("Parsed integer:", result.Value)
	} else {
	    fmt.Println("Failed to parse integer:", result.Error)
	}

Get a single integer query parameter:

	u, _ := url.Parse("http://example.com?age=18")
	value, ok := qp.GetInt(u, "age")
	if ok {
	    fmt.Println("Parsed integer:", value)
	} else {
	    fmt.Println("Failed to parse integer")
	}

Parse a slice of integer query parameters:

	u, _ := url.Parse("http://example.com?ids=1,2,3")
	result := qp.ParseIntSlice(u, "ids")
	if result.Contains && !result.Empty && result.Error == nil {
	    fmt.Println("Parsed integers:", result.Value)
	} else {
	    fmt.Println("Failed to parse integers:", result.Error)
	}

Get a slice of integer query parameters:

	u, _ := url.Parse("http://example.com?ids=1,2,3") // or ?ids=1&ids=2&ids=3
	ids := qp.GetIntSlice(u, "ids")
	if ids != nil {
	    fmt.Println("Parsed integers:", ids)
	} else {
	    fmt.Println("Failed to parse integers")
	}

# Parsing Float Query Parameters

Parse a single float query parameter:

	u, _ := url.Parse("http://example.com?temperature=36.6")
	result := qp.ParseFloat(u, "temperature")
	if result.Contains && !result.Empty && result.Error == nil {
	    fmt.Println("Parsed float:", result.Value)
	} else {
	    fmt.Println("Failed to parse float:", result.Error)
	}

Get a single float query parameter:

	u, _ := url.Parse("http://example.com?temperature=36.6")
	value, ok := qp.GetFloat(u, "temperature")
	if ok {
	    fmt.Println("Parsed float:", value)
	} else {
	    fmt.Println("Failed to parse float")
	}

Parse a slice of float query parameters:

	u, _ := url.Parse("http://example.com?ids=1.1,2.2,3.3")
	result := qp.ParseFloatSlice(u, "ids")
	if result.Contains && !result.Empty && result.Error == nil {
	    fmt.Println("Parsed floats:", result.Value)
	} else {
	    fmt.Println("Failed to parse floats:", result.Error)
	}

Get a slice of float query parameters:

	u, _ := url.Parse("http://example.com?ids=1.1,2.2,3.3")
	ids := qp.GetFloatSlice(u, "ids")
	if ids != nil {
	    fmt.Println("Parsed floats:", ids)
	} else {
	    fmt.Println("Failed to parse floats")
	}

# Parsing String Query Parameters

Parse a single string query parameter:

	u, _ := url.Parse("http://example.com?name=alice")
	result := qp.ParseString(u, "name")
	if result.Contains && !result.Empty && result.Error == nil {
	    fmt.Println("Parsed string:", result.Value)
	} else {
	    fmt.Println("Failed to parse string:", result.Error)
	}

Get a single string query parameter:

	u, _ := url.Parse("http://example.com?name=alice")
	value, ok := qp.GetString(u, "name")
	if ok {
	    fmt.Println("Parsed string:", value)
	} else {
	    fmt.Println("Failed to parse string")
	}

Parse a slice of string query parameters:

	u, _ := url.Parse("http://example.com?names=alice,bob,charlie")
	result := qp.ParseStringSlice(u, "names")
	if result.Contains && !result.Empty && result.Error == nil {
	    fmt.Println("Parsed strings:", result.Value)
	} else {
	    fmt.Println("Failed to parse strings:", result.Error)
	}

Get a slice of string query parameters:

	u, _ := url.Parse("http://example.com?names=alice,bob,charlie")
	ids := qp.GetStringSlice(u, "names")
	if ids != nil {
	    fmt.Println("Parsed strings:", ids)
	} else {
	    fmt.Println("Failed to parse strings")
	}

# Parsing Boolean Query Parameters

Parse a single boolean query parameter:

	u, _ := url.Parse("http://example.com?active=true")
	result := qp.ParseBool(u, "active")
	if result.Contains && !result.Empty && result.Error == nil {
	    fmt.Println("Parsed boolean:", result.Value)
	} else {
	    fmt.Println("Failed to parse boolean:", result.Error)
	}

Get a single boolean query parameter:

	u, _ := url.Parse("http://example.com?active=true")
	value, ok := qp.GetBool(u, "active")
	if ok {
	    fmt.Println("Parsed boolean:", value)
	} else {
	    fmt.Println("Failed to parse boolean")
	}

Parse a slice of boolean query parameters:

	u, _ := url.Parse("http://example.com?flags=true,false,yes,no")
	result := qp.ParseBoolSlice(u, "flags")
	if result.Contains && !result.Empty && result.Error == nil {
	    fmt.Println("Parsed booleans:", result.Value)
	} else {
	    fmt.Println("Failed to parse booleans:", result.Error)
	}

Get a slice of boolean query parameters:

	u, _ := url.Parse("http://example.com?flags=true,false,yes,no")
	ids := qp.GetBoolSlice(u, "flags")
	if ids != nil {
	    fmt.Println("Parsed booleans:", ids)
	} else {
	    fmt.Println("Failed to parse booleans")
	}

# Checking for Presence and Emptiness of Query Parameters

Check if a query parameter is present:

	u, _ := url.Parse("http://example.com?age=18")
	if qp.Contains(u, "age") {
	    fmt.Println("Query parameter 'age' is present")
	} else {
	    fmt.Println("Query parameter 'age' is absent")
	}

Check if a query parameter is empty:

	u, _ := url.Parse("http://example.com?age=")
	if qp.Empty(u, "age") {
	    fmt.Println("Query parameter 'age' is empty")
	} else {
	    fmt.Println("Query parameter 'age' is not empty")
	}
*/
package qp
