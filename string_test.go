package qp

import (
	"errors"
	"net/url"
	"reflect"
	"testing"
)

// TestParseString tests the ParseString function.
func TestParseString(t *testing.T) {
	tests := []struct {
		name     string
		query    string
		opt      []string
		expected *Result[string]
	}{
		{
			name:  "Simple call",
			query: "name=john",
			opt:   nil,
			expected: &Result[string]{
				Key:     "name",
				Value:   "john",
				Default: "",
				Others:  nil,

				Empty:    false,
				Contains: true,
				Error:    nil,
			},
		},
		{
			name:  "Default value",
			query: "",
			opt:   nil,
			expected: &Result[string]{
				Key:     "name",
				Value:   "",
				Default: "",
				Others:  nil,

				Empty:    true,
				Contains: false,
				Error:    nil,
			},
		},
		{
			name:  "With default value",
			query: "name=john",
			opt:   []string{"doe"},
			expected: &Result[string]{
				Key:     "name",
				Value:   "john",
				Default: "doe",
				Others:  nil,

				Empty:    false,
				Contains: true,
				Error:    nil,
			},
		},
		{
			name:  "Default is part of valid values",
			query: "name=nik",
			opt:   []string{"nik", "bob", "john"},
			expected: &Result[string]{
				Key:     "name",
				Value:   "nik",
				Default: "nik",
				Others:  []string{"nik", "bob", "john"},

				Empty:    false,
				Contains: true,
				Error:    nil,
			},
		},
		{
			name:  "Out of range",
			query: "name=lisa",
			opt:   []string{"nik", "bob", "john"},
			expected: &Result[string]{
				Key:     "name",
				Value:   "nik",
				Default: "nik",
				Others:  []string{"nik", "bob", "john"},

				Empty:    false,
				Contains: true,
				Error:    errors.New("value out of range"),
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			u, _ := url.Parse("http://example.com?" + tc.query)
			got := ParseString(u, tc.expected.Key, tc.opt...)

			if got.Value != tc.expected.Value {
				t.Errorf("ParseString() .Value: got = %v, want %v",
					got.Value, tc.expected.Value)
			}

			if got.Default != tc.expected.Default {
				t.Errorf("ParseString() .Default: got = %v, want %v",
					got.Default, tc.expected.Default)
			}

			if !reflect.DeepEqual(got.Others, tc.expected.Others) {
				t.Errorf("ParseString() .Others: got = %v, want %v",
					got.Others, tc.expected.Others)
			}

			if got.Empty != tc.expected.Empty {
				t.Errorf("ParseString() .Empty: got = %v, want %v",
					got.Empty, tc.expected.Empty)
			}

			if got.Contains != tc.expected.Contains {
				t.Errorf("ParseString() .Contains: got = %v, want %v",
					got.Contains, tc.expected.Contains)
			}

			if (got.Error != nil && tc.expected.Error == nil) ||
				(got.Error == nil && tc.expected.Error != nil) {
				t.Errorf("ParseString() .Error: got = %v, want %v",
					got.Error, tc.expected.Error)
			}
		})
	}
}

func TestGetString(t *testing.T) {
	key := "name"
	tests := []struct {
		name     string
		query    string
		opt      []string
		expected string
		ok       bool
	}{
		{
			name:     "Simple call",
			query:    "name=john",
			opt:      nil,
			expected: "john",
			ok:       true,
		},
		{
			name:     "Default value",
			query:    "",
			opt:      nil,
			expected: "",
			ok:       false,
		},
		{
			name:     "With default value",
			query:    "name=",
			opt:      []string{"bob"},
			expected: "bob",
			ok:       false,
		},
		{
			name:     "Default is part of valid values",
			query:    "name=nik",
			opt:      []string{"nik", "bob", "john"},
			expected: "nik",
			ok:       true,
		},
		{
			name:     "Out of range",
			query:    "name=lisa",
			opt:      []string{"nik", "bob", "john"},
			expected: "nik",
			ok:       false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			u, _ := url.Parse("http://example.com?" + tc.query)
			got, ok := GetString(u, key, tc.opt...)

			if got != tc.expected {
				t.Errorf("GetString() .Value = %v, want %v", got, tc.expected)
			}

			if ok != tc.ok {
				t.Errorf("GetString() .Ok = %v, want %v", ok, tc.ok)
			}
		})
	}
}

// TestParseStringSlice tests the ParseStringSlice function.
func TestParseStringSlice(t *testing.T) {
	tests := []struct {
		name     string
		query    string
		expected *Result[[]string]
	}{
		{
			name:  "Simple call",
			query: "names=alice",
			expected: &Result[[]string]{
				Key:      "names",
				Value:    []string{"alice"},
				Empty:    false,
				Contains: true,
				Error:    nil,
			},
		},
		{
			name:  "Slice as single value",
			query: "names=alice,bob,charlie",
			expected: &Result[[]string]{
				Key:      "names",
				Value:    []string{"alice", "bob", "charlie"},
				Empty:    false,
				Contains: true,
				Error:    nil,
			},
		},
		{
			name:  "Slice as multiple values",
			query: "names=alice&names=bob&names=charlie",
			expected: &Result[[]string]{
				Key:      "names",
				Value:    []string{"alice", "bob", "charlie"},
				Empty:    false,
				Contains: true,
				Error:    nil,
			},
		},
		{
			name:  "Empty value",
			query: "names=",
			expected: &Result[[]string]{
				Key:      "names",
				Value:    []string{},
				Empty:    true,
				Contains: true,
				Error:    nil,
			},
		},
		{
			name:  "Declared only",
			query: "names",
			expected: &Result[[]string]{
				Key:      "names",
				Value:    []string{},
				Empty:    true,
				Contains: true,
				Error:    nil,
			},
		},
		{
			name:  "Without key",
			query: "age=18",
			expected: &Result[[]string]{
				Key:      "names",
				Value:    nil,
				Empty:    true,
				Contains: false,
				Error:    nil,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			u, _ := url.Parse("http://example.com?" + tc.query)
			got := ParseStringSlice(u, tc.expected.Key)

			if got.Value == nil && tc.expected.Value != nil ||
				got.Value != nil && tc.expected.Value == nil {
				t.Errorf("ParseStringSlice() .Value*: got = %v, want %v",
					got.Value, tc.expected.Value)
			} else if got.Value != nil && tc.expected.Value != nil {
				if !reflect.DeepEqual(got.Value, tc.expected.Value) {
					t.Errorf("ParseStringSlice() .Value: got = %v, want %v",
						got.Value, tc.expected.Value)
				}
			}

			if got.Empty != tc.expected.Empty {
				t.Errorf("ParseStringSlice() .Empty: got = %v, want %v",
					got.Empty, tc.expected.Empty)
			}

			if got.Contains != tc.expected.Contains {
				t.Errorf("ParseStringSlice() .Contains: got = %v, want %v",
					got.Contains, tc.expected.Contains)
			}

			if (got.Error != nil && tc.expected.Error == nil) ||
				(got.Error == nil && tc.expected.Error != nil) {
				t.Errorf("ParseStringSlice() .Error: got = %v, want %v",
					got.Error, tc.expected.Error)
			}
		})
	}
}

// TestGetStringSlice tests the GetStringSlice function.
func TestGetStringSlice(t *testing.T) {
	key := "names"
	tests := []struct {
		name     string
		query    string
		expected []string
	}{
		{
			name:     "Simple call",
			query:    "names=alice",
			expected: []string{"alice"},
		},
		{
			name:     "Slice as single value",
			query:    "names=alice,bob,charlie",
			expected: []string{"alice", "bob", "charlie"},
		},
		{
			name:     "Slice as multiple values",
			query:    "names=alice&names=bob&names=charlie",
			expected: []string{"alice", "bob", "charlie"},
		},
		{
			name:     "Empty value",
			query:    "names=",
			expected: []string{},
		},
		{
			name:     "Declared only",
			query:    "names",
			expected: []string{},
		},
		{
			name:     "Without key",
			query:    "age=18",
			expected: nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			u, _ := url.Parse("http://example.com?" + tc.query)
			got := GetStringSlice(u, key)

			if got == nil && tc.expected != nil {
				t.Errorf("GetStringSlice() = nil, want %v", tc.expected)
			} else if got != nil && tc.expected == nil {
				t.Errorf("GetStringSlice() = %v, want nil", got)
			} else if got != nil && tc.expected != nil {
				if !reflect.DeepEqual(got, tc.expected) {
					t.Errorf("GetStringSlice() = %v, want %v", got, tc.expected)
				}
			}
		})
	}
}
