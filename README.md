[![Go Report Card](https://goreportcard.com/badge/github.com/goloop/qp)](https://goreportcard.com/report/github.com/goloop/qp) [![License](https://img.shields.io/badge/license-MIT-brightgreen)](https://github.com/goloop/qp/blob/master/LICENSE) [![License](https://img.shields.io/badge/godoc-YES-green)](https://godoc.org/github.com/goloop/qp) [![Stay with Ukraine](https://img.shields.io/static/v1?label=Stay%20with&message=Ukraine%20♥&color=ffD700&labelColor=0057B8&style=flat)](https://u24.gov.ua/)


# qp

The **qp** package provides robust and user-friendly utilities for parsing query parameters from URLs in Go applications.


## Installation

```sh
go get github.com/goloop/qp
```

## Features

- Parse query parameters into various types (int, float64, string, bool)
- Support for slices of values
- Range validation for numeric types
- Value validation against allowed sets
- Multiple boolean formats support (true/false, yes/no, on/off, 1/0)
- Detailed error reporting
- Zero dependencies

## Method Types

The package provides three types of methods for each data type:

1. **Parse methods** - Return detailed results in a Result structure:
   - Technical data like errors and presence flags
   - Parsed values with validation status
   - Default values and constraints

2. **Get methods** - Return value and validity boolean:
   - Value is never nil (even for slices)
   - Boolean is true only if parameter exists and parsed correctly

3. **Pull methods** - Return pointer:
   - nil if parameter is absent
   - Points to default value if parameter present but invalid
   - Points to parsed value if successful

## Usage Examples

### Boolean Parsing

```go
// Parse single boolean.
u, _ := url.Parse("http://example.com?active=true")

// Using Parse method.
result := qp.ParseBool(u, "active")
if result.Contains && !result.Empty && result.Error == nil {
    fmt.Println("Active:", result.Value)
}

// Using Get method.
value, ok := qp.GetBool(u, "active")

// Using Pull method.
active := qp.PullBool(u, "active")

// Parse boolean slice.
u, _ := url.Parse("http://example.com?flags=true,false,yes,no")
flags := qp.PullBoolSlice(u, "flags")
```

### Integer Parsing

```go
// Parse single integer.
u, _ := url.Parse("http://example.com?age=25")

// Basic parsing.
age := qp.PullInt(u, "age")

// With default value.
age := qp.ParseInt(u, "age", 18)

// With range validation.
age := qp.ParseInt(u, "age", 18, 65) // default: 18, range: 18-65

// With additional valid values.
age := qp.ParseInt(u, "age", 18, 65, 16, 70) // allows 16 and 70 besides range

// Parse integer slice.
u, _ := url.Parse("http://example.com?ids=1,2,3")
ids := qp.PullIntSlice(u, "ids")
```

### Float Parsing

```go
// Parse single float.
u, _ := url.Parse("http://example.com?temperature=36.6")

// Basic parsing.
temp := qp.PullFloat(u, "temperature")

// With range validation.
temp := qp.ParseFloat(u, "temperature", 36.0, 42.0)

// Parse float slice.
u, _ := url.Parse("http://example.com?temps=36.6,37.2,36.9")
temps := qp.PullFloatSlice(u, "temps")
```

### String Parsing

```go
// Parse single string.
u, _ := url.Parse("http://example.com?name=alice")

// Basic parsing.
name := qp.PullString(u, "name")

// With default value.
name := qp.ParseString(u, "name", "guest")

// With valid values validation.
role := qp.ParseString(u, "role", "guest", "admin", "user")

// Parse string slice.
u, _ := url.Parse("http://example.com?names=alice,bob,charlie")
names := qp.PullStringSlice(u, "names")
```

### Practical Example: SQL WHERE Clause

```go
func where(isActive, isStaff, isSuperuser *bool) string {
    check := [...]struct {
        name  string
        value *bool
    }{
        {name: "is_active", value: isActive},
        {name: "is_staff", value: isStaff},
        {name: "is_superuser", value: isSuperuser},
    }

    and := make([]string, 0, len(check))
    for _, m := range check {
        if m.value != nil {
            and = append(and, fmt.Sprintf("%s=%t", m.name, *m.value))
        }
    }

    if len(and) != 0 {
        return "WHERE " + strings.Join(and, ", ")
    }
    return ""
}

u, _ := url.Parse("http://example.com?is_active=true&is_staff=false")
isActive := qp.PullBool(u, "is_active")
isStaff := qp.PullBool(u, "is_staff")
isSuperuser := qp.PullBool(u, "is_superuser") // nil
query := where(isActive, isStaff, isSuperuser)
// Result: WHERE is_active=true, is_staff=false
```

### Utility Functions

```go
// Check parameter presence.
if qp.Contains(u, "age") {
    fmt.Println("Age parameter exists")
}

// Check if parameter is empty.
if qp.Empty(u, "age") {
    fmt.Println("Age parameter is empty")
}
```

## Notes

- Boolean parsing supports multiple formats: true/false, yes/no, on/off, 1/0
- Slice parameters can be specified in two ways:
  - Comma-separated: `?ids=1,2,3`
  - Multiple parameters: `?ids=1&ids=2&ids=3`
- Numeric parsers support range validation and additional valid values
- String parsers support validation against a list of valid values


## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
