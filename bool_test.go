package qp

import (
	"errors"
	"net/url"
	"reflect"
	"testing"

	"github.com/goloop/g"
)

// TestParseBool tests the ParseBool function.
func TestParseBool(t *testing.T) {
	tests := []struct {
		name     string
		query    string
		opt      []bool
		expected *Result[bool]
	}{
		{
			name:  "Constants: true or false",
			query: "is_active=true",
			opt:   nil,
			expected: &Result[bool]{
				Key:     "is_active",
				Value:   true,
				Default: false,

				Empty:    false,
				Contains: true,
				Error:    nil,
			},
		},
		{
			name:  "Constants: yes or no",
			query: "is_active=yes",
			opt:   nil,
			expected: &Result[bool]{
				Key:     "is_active",
				Value:   true,
				Default: false,

				Empty:    false,
				Contains: true,
				Error:    nil,
			},
		},
		{
			name:  "Constants: 1 or 0",
			query: "is_active=1",
			opt:   nil,
			expected: &Result[bool]{
				Key:     "is_active",
				Value:   true,
				Default: false,

				Empty:    false,
				Contains: true,
				Error:    nil,
			},
		},
		{
			name:  "Constants: on or off",
			query: "is_active=on",
			opt:   nil,
			expected: &Result[bool]{
				Key:     "is_active",
				Value:   true,
				Default: false,

				Empty:    false,
				Contains: true,
				Error:    nil,
			},
		},
		{
			name:  "Empty value",
			query: "is_active=",
			opt:   nil,
			expected: &Result[bool]{
				Key:     "is_active",
				Value:   false,
				Default: false,

				Empty:    true,
				Contains: true,
				Error:    nil,
			},
		},
		{
			name:  "Marked key",
			query: "is_active",
			opt:   nil,
			expected: &Result[bool]{
				Key:     "is_active",
				Value:   false,
				Default: false,

				Empty:    true,
				Contains: true,
				Error:    nil,
			},
		},
		{
			name:  "With default value",
			query: "is_active=hello",
			opt:   []bool{true},
			expected: &Result[bool]{
				Key:     "is_active",
				Value:   true,
				Default: true,

				Empty:    false,
				Contains: true,
				Error:    errors.New("incorrect value"),
			},
		},
		{
			name:  "Without key",
			query: "message=hello",
			opt:   []bool{true, false, true},
			expected: &Result[bool]{
				Key:     "is_active",
				Value:   true,
				Default: true,

				Empty:    true,
				Contains: false,
				Error:    nil,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			u, _ := url.Parse("http://example.com?" + tc.query)
			got := ParseBool(u, tc.expected.Key, tc.opt...)

			if got.Value != tc.expected.Value {
				t.Errorf("ParseInt() .Value: got = %v, want %v",
					got.Value, tc.expected.Value)
			}

			if got.Default != tc.expected.Default {
				t.Errorf("ParseInt() .Default: got = %v, want %v",
					got.Default, tc.expected.Default)
			}

			if got.Empty != tc.expected.Empty {
				t.Errorf("ParseInt() .Empty: got = %v, want %v",
					got.Empty, tc.expected.Empty)
			}

			if got.Contains != tc.expected.Contains {
				t.Errorf("ParseInt() .Contains: got = %v, want %v",
					got.Contains, tc.expected.Contains)
			}

			if (got.Error != nil && tc.expected.Error == nil) ||
				(got.Error == nil && tc.expected.Error != nil) {
				t.Errorf("ParseInt() .Error: got = %v, want %v",
					got.Error, tc.expected.Error)
			}
		})
	}
}

// TestGetBool tests the ParseBool function.
func TestGetBool(t *testing.T) {
	key := "is_active"
	tests := []struct {
		name     string
		query    string
		opt      []bool
		expected bool
		ok       bool
	}{
		{
			name:     "Constants: true or false",
			query:    "is_active=true",
			opt:      nil,
			expected: true,
			ok:       true,
		},
		{
			name:     "Constants: yes or no",
			query:    "is_active=yes",
			opt:      nil,
			expected: true,
			ok:       true,
		},
		{
			name:     "Constants: 1 or 0",
			query:    "is_active=1",
			opt:      nil,
			expected: true,
			ok:       true,
		},
		{
			name:     "Constants: on or off",
			query:    "is_active=on",
			opt:      nil,
			expected: true,
			ok:       true,
		},
		{
			name:     "Empty value",
			query:    "is_active=",
			opt:      nil,
			expected: false,
			ok:       false,
		},
		{
			name:     "Marked key",
			query:    "is_active",
			opt:      nil,
			expected: false,
			ok:       false,
		},
		{
			name:     "With default value",
			query:    "is_active=hello",
			opt:      []bool{true},
			expected: true,
			ok:       false,
		},
		{
			name:     "Without key",
			query:    "message=hello",
			opt:      []bool{true, false, true},
			expected: true,
			ok:       false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			u, _ := url.Parse("http://example.com?" + tc.query)
			got, ok := GetBool(u, key, tc.opt...)

			if got != tc.expected {
				t.Errorf("GetBool() .Value = %v, want %v", got, tc.expected)
			}

			if ok != tc.ok {
				t.Errorf("GetBool() .Ok = %v, want %v", ok, tc.ok)
			}
		})
	}
}

// TestPullBool tests the PullBool function.
func TestPullBool(t *testing.T) {
	key := "is_active"
	tests := []struct {
		name     string
		query    string
		opt      []bool
		expected *bool
	}{
		{
			name:     "Constants: true or false",
			query:    "is_active=true",
			opt:      nil,
			expected: g.Ptr(true),
		},
		{
			name:     "Constants: yes or no",
			query:    "is_active=yes",
			opt:      nil,
			expected: g.Ptr(true),
		},
		{
			name:     "Constants: 1 or 0",
			query:    "is_active=1",
			opt:      nil,
			expected: g.Ptr(true),
		},
		{
			name:     "Constants: on or off",
			query:    "is_active=on",
			opt:      nil,
			expected: g.Ptr(true),
		},
		{
			name:     "Empty value",
			query:    "is_active=",
			opt:      nil,
			expected: g.Ptr(false),
		},
		{
			name:     "Marked key",
			query:    "is_active",
			opt:      nil,
			expected: g.Ptr(false),
		},
		{
			name:     "With default value",
			query:    "is_active=hello",
			opt:      []bool{true},
			expected: g.Ptr(true),
		},
		{
			name:     "Without key",
			query:    "message=hello",
			opt:      []bool{true, false, true},
			expected: nil, // nil, but not a pointer
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			u, _ := url.Parse("http://example.com?" + tc.query)
			got := PullBool(u, key, tc.opt...)

			if (got == nil && tc.expected != nil) ||
				(got != nil && tc.expected == nil) {
				t.Errorf("PullBool() = %v, want %v", got, tc.expected)
			} else if got != nil && tc.expected != nil {
				if *got != *tc.expected {
					t.Errorf("PullBool() = %v, want %v", *got, *tc.expected)
				}
			}
		})
	}
}

// TestParseBoolSlice tests the ParseBoolSlice function.
func TestParseBoolSlice(t *testing.T) {
	tests := []struct {
		name     string
		query    string
		expected *Result[[]bool]
	}{
		{
			name:  "Simple call",
			query: "flags=true",
			expected: &Result[[]bool]{
				Key:      "flags",
				Value:    []bool{true},
				Empty:    false,
				Contains: true,
				Error:    nil,
			},
		},
		{
			name:  "Slice as single value",
			query: "flags=true,false,yes,no",
			expected: &Result[[]bool]{
				Key:      "flags",
				Value:    []bool{true, false, true, false},
				Empty:    false,
				Contains: true,
				Error:    nil,
			},
		},
		{
			name:  "Slice as multiple values",
			query: "flags=true&flags=false&flags=yes&flags=no",
			expected: &Result[[]bool]{
				Key:      "flags",
				Value:    []bool{true, false, true, false},
				Empty:    false,
				Contains: true,
				Error:    nil,
			},
		},
		{
			name:  "Empty value",
			query: "flags=",
			expected: &Result[[]bool]{
				Key:      "flags",
				Value:    []bool{},
				Empty:    true,
				Contains: true,
				Error:    nil,
			},
		},
		{
			name:  "Declared only",
			query: "flags",
			expected: &Result[[]bool]{
				Key:      "flags",
				Value:    []bool{},
				Empty:    true,
				Contains: true,
				Error:    nil,
			},
		},
		{
			name:  "Without key",
			query: "age=18",
			expected: &Result[[]bool]{
				Key:      "flags",
				Value:    []bool{},
				Empty:    true,
				Contains: false,
				Error:    nil,
			},
		},
		{
			name:  "Incorrect as single value",
			query: "flags=true,string,false",
			expected: &Result[[]bool]{
				Key:      "flags",
				Value:    []bool{},
				Empty:    false,
				Contains: true,
				Error:    errors.New("incorrect value"),
			},
		},
		{
			name:  "Incorrect as multiple values",
			query: "flags=true&flags=string&flags=false",
			expected: &Result[[]bool]{
				Key:      "flags",
				Value:    []bool{},
				Empty:    false,
				Contains: true,
				Error:    errors.New("incorrect value"),
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			u, _ := url.Parse("http://example.com?" + tc.query)
			got := ParseBoolSlice(u, tc.expected.Key)

			if got.Value == nil && tc.expected.Value != nil ||
				got.Value != nil && tc.expected.Value == nil {
				t.Errorf("ParseBoolSlice() .Value*: got = %v, want %v",
					got.Value, tc.expected.Value)
			} else if got.Value != nil && tc.expected.Value != nil {
				if !reflect.DeepEqual(got.Value, tc.expected.Value) {
					t.Errorf("ParseBoolSlice() .Value: got = %v, want %v",
						got.Value, tc.expected.Value)
				}
			}

			if got.Empty != tc.expected.Empty {
				t.Errorf("ParseBoolSlice() .Empty: got = %v, want %v",
					got.Empty, tc.expected.Empty)
			}

			if got.Contains != tc.expected.Contains {
				t.Errorf("ParseBoolSlice() .Contains: got = %v, want %v",
					got.Contains, tc.expected.Contains)
			}

			if (got.Error != nil && tc.expected.Error == nil) ||
				(got.Error == nil && tc.expected.Error != nil) {
				t.Errorf("ParseBoolSlice() .Error: got = %v, want %v",
					got.Error, tc.expected.Error)
			}
		})
	}
}

// TestGetBoolSlice tests the GetBoolSlice function.
func TestGetBoolSlice(t *testing.T) {
	key := "flags"
	tests := []struct {
		name     string
		query    string
		opt      [][]bool
		expected []bool
		ok       bool
	}{
		{
			name:     "Simple call",
			query:    "flags=true",
			expected: []bool{true},
			ok:       true,
		},
		{
			name:     "Slice as single value",
			query:    "flags=true,false,yes,no",
			expected: []bool{true, false, true, false},
			ok:       true,
		},
		{
			name:     "Slice as multiple values",
			query:    "flags=true&flags=false&flags=yes&flags=no",
			expected: []bool{true, false, true, false},
			ok:       true,
		},
		{
			name:     "Empty value",
			query:    "flags=",
			expected: []bool{},
			ok:       false,
		},
		{
			name:     "Default value",
			query:    "flags=",
			opt:      [][]bool{{true, true, true}},
			expected: []bool{true, true, true},
			ok:       false,
		},
		{
			name:     "Declared only",
			query:    "flags",
			expected: []bool{},
			ok:       false,
		},
		{
			name:     "Without key",
			query:    "age=18",
			expected: []bool{},
			ok:       false,
		},
		{
			name:     "Incorrect as single value",
			query:    "flags=true,string,false",
			expected: []bool{},
			ok:       false,
		},
		{
			name:     "Incorrect as multiple values",
			query:    "flags=true&flags=string&flags=false",
			expected: []bool{},
			ok:       false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			u, _ := url.Parse("http://example.com?" + tc.query)
			got, ok := GetBoolSlice(u, key, tc.opt...)

			if !reflect.DeepEqual(got, tc.expected) {
				t.Errorf("GetBoolSlice() .Value = %v, want %v", got, tc.expected)
			}

			if ok != tc.ok {
				t.Errorf("GetBoolSlice() .Ok = %v, want %v", ok, tc.ok)
			}
		})
	}
}

// TestPullBoolSlice tests the PullBoolSlice function.
func TestPullBoolSlice(t *testing.T) {
	key := "flags"
	tests := []struct {
		name     string
		query    string
		expected []bool
	}{
		{
			name:     "Simple call",
			query:    "flags=true",
			expected: []bool{true},
		},
		{
			name:     "Slice as single value",
			query:    "flags=true,false,yes,no",
			expected: []bool{true, false, true, false},
		},
		{
			name:     "Slice as multiple values",
			query:    "flags=true&flags=false&flags=yes&flags=no",
			expected: []bool{true, false, true, false},
		},
		{
			name:     "Empty value",
			query:    "flags=",
			expected: []bool{},
		},
		{
			name:     "Declared only",
			query:    "flags",
			expected: []bool{},
		},
		{
			name:     "Without key",
			query:    "age=18",
			expected: nil,
		},
		{
			name:     "Incorrect as single value",
			query:    "flags=true,string,false",
			expected: []bool{},
		},
		{
			name:     "Incorrect as multiple values",
			query:    "flags=true&flags=string&flags=false",
			expected: []bool{},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			u, _ := url.Parse("http://example.com?" + tc.query)
			got := PullBoolSlice(u, key)

			if got == nil && tc.expected != nil {
				t.Errorf("PullBoolSlice() = nil, want %v", tc.expected)
			} else if got != nil && tc.expected == nil {
				t.Errorf("PullBoolSlice() = %v, want nil", got)
			} else if got != nil && tc.expected != nil {
				if !reflect.DeepEqual(got, tc.expected) {
					t.Errorf("PullBoolSlice() = %v, want %v", got, tc.expected)
				}
			}
		})
	}
}
