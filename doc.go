/*
Package qp provides utilities for parsing query parameters from URLs.

This package includes functions for parsing various types of query parameters,
such as integers, floats, strings, and booleans. It supports parsing both single
values and slices of values. Additionally, it provides convenient functions for
checking the presence and emptiness of query parameters.

The functions in this package are designed to be robust and user-friendly,
returning detailed results that include the parsed value, default value, and
information about whether the query parameter was present and valid.

# Parsing Functions Details

The methods in this package are categorized into three types: Parse, Get,
and Pull.

- Parse methods return detailed results in a Result structure. This structure
includes technical data such as errors, presence, and emptiness of the key,
along with the parsed value, default value, and limits if any.

- Get methods always return an object (never nil, even for slices) and
a boolean indicating the presence and validity of the key. The boolean is
true only if the key was found and the value was correctly parsed and falls
within the specified limits.

- Pull methods return a pointer or nil:
  - nil if the key is absent;
  - a pointer to the default value if the key is present but empty
    or has invalid data;
  - a pointer to the parsed result if everything is parsed correctly.

# Example Usage

Here are some examples of how to use the functions provided by this package:

# Parsing Integer Query Parameters

Parse a single integer query parameter:

	u, _ := url.Parse("http://example.com?age=18")
	result := qp.ParseInt(u, "age")
	if result.Contains && !result.Empty && result.Error == nil {
	    fmt.Println("Parsed integer:", result.Value)
	} else {
	    fmt.Println("Is empty:", result.Empty)
	    fmt.Println("Is present:", result.Contains)
	    fmt.Println("Failed to parse integer:", result.Error)
	}

Get a single integer query parameter:

	u, _ := url.Parse("http://example.com?age=18")
	value, ok := qp.GetInt(u, "age")
	if ok {
	    fmt.Println("Parsed integer:", value)
	} else {
	    fmt.Println("Default value:", value)
	}

Pull a single integer query parameter:

	u, _ := url.Parse("http://example.com?age=18")
	age := qp.PullInt(u, "age")
	if age != nil {
	    fmt.Println("Parsed integer:", *age)
	} else {
	    fmt.Println("Key 'age' is absent")
	}

Parse a slice of integer query parameters:

	u, _ := url.Parse("http://example.com?ids=1,2,3")
	result := qp.ParseIntSlice(u, "ids")
	if result.Contains && !result.Empty && result.Error == nil {
	    fmt.Println("Parsed integers:", result.Value)
	} else {
	    fmt.Println("Is empty:", result.Empty)
	    fmt.Println("Is present:", result.Contains)
	    fmt.Println("Failed to parse integers:", result.Error)
	}

Get a slice of integer query parameters:

	u, _ := url.Parse("http://example.com?ids=1,2,3")
	ids, ok := qp.GetIntSlice(u, "ids")
	if ok {
	    fmt.Println("Parsed integers:", ids)
	} else {
	    fmt.Println("Default value:", ids)
	}

Pull a slice of integer query parameters:

	u, _ := url.Parse("http://example.com?ids=1,2,3") // or ?ids=1&ids=2&ids=3
	ids := qp.PullIntSlice(u, "ids")
	if ids != nil {
	    fmt.Println("Parsed integers:", ids)
	} else {
	    fmt.Println("Key 'ids' is absent")
	}

# Parsing Float Query Parameters

Parse a single float query parameter:

	u, _ := url.Parse("http://example.com?temperature=36.6")
	result := qp.ParseFloat(u, "temperature")
	if result.Contains && !result.Empty && result.Error == nil {
	    fmt.Println("Parsed float:", result.Value)
	} else {
	    fmt.Println("Is empty:", result.Empty)
	    fmt.Println("Is present:", result.Contains)
	    fmt.Println("Failed to parse float:", result.Error)
	}

Get a single float query parameter:

	u, _ := url.Parse("http://example.com?temperature=36.6")
	value, ok := qp.GetFloat(u, "temperature")
	if ok {
	    fmt.Println("Parsed float:", value)
	} else {
	    fmt.Println("Default value:", value)
	}

Pull a single float query parameter:

	u, _ := url.Parse("http://example.com?temperature=36.6")
	temperature := qp.PullFloat(u, "temperature")
	if temperature != nil {
	    fmt.Println("Parsed float:", *temperature)
	} else {
	    fmt.Println("Key 'temperature' is absent")
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
	ids, ok := qp.GetFloatSlice(u, "ids")
	if ok {
	    fmt.Println("Parsed floats:", ids)
	} else {
	    fmt.Println("Default value:", ids)
	}

Pull a slice of float query parameters:

	u, _ := url.Parse("http://example.com?ids=1.1,2.2,3.3")
	ids := qp.PullFloatSlice(u, "ids")
	if ids != nil {
	    fmt.Println("Parsed floats:", ids)
	} else {
	    fmt.Println("Key 'ids' is absent")
	}

# Parsing String Query Parameters

Parse a single string query parameter:

	u, _ := url.Parse("http://example.com?name=alice")
	result := qp.ParseString(u, "name")
	if result.Contains && !result.Empty && result.Error == nil {
	    fmt.Println("Parsed string:", result.Value)
	} else {
	    fmt.Println("Is empty:", result.Empty)
	    fmt.Println("Is present:", result.Contains)
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

Pull a single string query parameter:

	u, _ := url.Parse("http://example.com?name=alice")
	name := qp.PullString(u, "name")
	if name != nil {
	    fmt.Println("Parsed string:", *name)
	} else {
	    fmt.Println("Key 'name' is absent")
	}

Parse a slice of string query parameters:

	u, _ := url.Parse("http://example.com?names=alice,bob,charlie")
	result := qp.ParseStringSlice(u, "names")
	if result.Contains && !result.Empty && result.Error == nil {
	    fmt.Println("Parsed strings:", result.Value)
	} else {
	    fmt.Println("Is empty:", result.Empty)
	    fmt.Println("Is present:", result.Contains)
	    fmt.Println("Failed to parse strings:", result.Error)
	}

Get a slice of string query parameters:

	u, _ := url.Parse("http://example.com?names=alice,bob,charlie")
	ids, ok := qp.GetStringSlice(u, "names")
	if ok {
	    fmt.Println("Parsed strings:", ids)
	} else {
	    fmt.Println("Failed to parse strings")
	}

Pull a slice of string query parameters:

	u, _ := url.Parse("http://example.com?names=alice,bob,charlie")
	ids := qp.PullStringSlice(u, "names")
	if ids != nil {
	    fmt.Println("Parsed strings:", ids)
	} else {
	    fmt.Println("Key 'names' is absent")
	}

# Parsing Boolean Query Parameters

Parse a single boolean query parameter:

	u, _ := url.Parse("http://example.com?active=true")
	result := qp.ParseBool(u, "active")
	if result.Contains && !result.Empty && result.Error == nil {
	    fmt.Println("Parsed boolean:", result.Value)
	} else {
	    fmt.Println("Is empty:", result.Empty)
	    fmt.Println("Is present:", result.Contains)
	    fmt.Println("Failed to parse boolean:", result.Error)
	}

Get a single boolean query parameter:

	u, _ := url.Parse("http://example.com?active=true")
	value, ok := qp.GetBool(u, "active")
	if ok {
	    fmt.Println("Parsed boolean:", value)
	} else {
	    fmt.Println("Default value:", value)
	}

Pull a single boolean query parameter:

	u, _ := url.Parse("http://example.com?active=true")
	active := qp.PullBool(u, "active")
	if active != nil {
	    fmt.Println("Parsed boolean:", *active)
	} else {
	    fmt.Println("Key 'active' is absent")
	}

Parse a slice of boolean query parameters:

	u, _ := url.Parse("http://example.com?flags=true,false,yes,no")
	result := qp.ParseBoolSlice(u, "flags")
	if result.Contains && !result.Empty && result.Error == nil {
	    fmt.Println("Parsed booleans:", result.Value)
	} else {
	    fmt.Println("Is empty:", result.Empty)
	    fmt.Println("Is present:", result.Contains)
	    fmt.Println("Failed to parse booleans:", result.Error)
	}

Get a slice of boolean query parameters:

	u, _ := url.Parse("http://example.com?flags=true,false,yes,no")
	ids, ok := qp.GetBoolSlice(u, "flags")
	if ok {
	    fmt.Println("Parsed booleans:", ids)
	} else {
	    fmt.Println("Default value:", ids)
	}

Pull a slice of boolean query parameters:

	u, _ := url.Parse("http://example.com?flags=true,false,yes,no")
	ids := qp.PullBoolSlice(u, "flags")
	if ids != nil {
	    fmt.Println("Parsed booleans:", ids)
	} else {
	    fmt.Println("Key 'flags' is absent")
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
