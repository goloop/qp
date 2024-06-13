package qp

import (
	"errors"
	"net/url"
	"reflect"
	"testing"

	"github.com/goloop/g"
)

// TestParseInt tests the ParseInt function.
func TestParseInt(t *testing.T) {
	tests := []struct {
		name     string
		query    string
		opt      []int
		expected *Result[int]
	}{
		{
			name:  "Simple call",
			query: "age=18",
			opt:   nil,
			expected: &Result[int]{
				Key:     "age",
				Value:   18,
				Default: 0,
				Min:     0,
				Max:     0,
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
			expected: &Result[int]{
				Key:     "age",
				Value:   0,
				Default: 0,
				Min:     0,
				Max:     0,
				Others:  nil,

				Empty:    true,
				Contains: false,
				Error:    nil,
			},
		},
		{
			name:  "With default value",
			query: "age=18",
			opt:   []int{21},
			expected: &Result[int]{
				Key:     "age",
				Value:   18,
				Default: 21,
				Min:     0,
				Max:     0,
				Others:  nil,

				Empty:    false,
				Contains: true,
				Error:    nil,
			},
		},
		{
			name:  "With default value and empty key",
			query: "age=",
			opt:   []int{21},
			expected: &Result[int]{
				Key:     "age",
				Value:   21,
				Default: 21,
				Min:     0,
				Max:     0,
				Others:  nil,

				Empty:    true,
				Contains: true,
				Error:    nil,
			},
		},
		{
			name:  "With range",
			query: "age=18",
			opt:   []int{18, 30},
			expected: &Result[int]{
				Key:     "age",
				Value:   18,
				Default: 18,
				Min:     18,
				Max:     30,
				Others:  nil,

				Empty:    false,
				Contains: true,
				Error:    nil,
			},
		},
		{
			name:  "With inverted range",
			query: "age=18",
			opt:   []int{30, 18},
			expected: &Result[int]{
				Key:     "age",
				Value:   18,
				Default: 30,
				Min:     18,
				Max:     30,
				Others:  nil,

				Empty:    false,
				Contains: true,
				Error:    nil,
			},
		},
		{
			name:  "With range and empty key",
			query: "age", // without =
			opt:   []int{18, 30},
			expected: &Result[int]{
				Key:     "age",
				Value:   18,
				Default: 18,
				Min:     18,
				Max:     30,
				Others:  nil,

				Empty:    true,
				Contains: true,
				Error:    nil,
			},
		},
		{
			name:  "Out of range",
			query: "age=55",
			opt:   []int{18, 30},
			expected: &Result[int]{
				Key:     "age",
				Value:   18,
				Default: 18,
				Min:     18,
				Max:     30,
				Others:  nil,

				Empty:    false,
				Contains: true,
				Error:    errors.New("value out of range"),
			},
		},
		{
			name:  "Incorrect value",
			query: "age=hello",
			opt:   []int{21},
			expected: &Result[int]{
				Key:     "age",
				Value:   21,
				Default: 21,
				Min:     0,
				Max:     0,
				Others:  nil,

				Empty:    false,
				Contains: true,
				Error:    errors.New("incorrect value"),
			},
		},
		{
			name:  "With other values",
			query: "age=70",
			opt:   []int{20, 20, 30, 50, 70},
			expected: &Result[int]{
				Key:     "age",
				Value:   70,
				Default: 20,
				Min:     20,
				Max:     20,
				Others:  []int{30, 50, 70},

				Empty:    false,
				Contains: true,
				Error:    nil,
			},
		},
		{
			name:  "With other values - out of range",
			query: "age=33",
			opt:   []int{20, 20, 30, 50, 70},
			expected: &Result[int]{
				Key:     "age",
				Value:   20,
				Default: 20,
				Min:     20,
				Max:     20,
				Others:  []int{30, 50, 70},

				Empty:    false,
				Contains: true,
				Error:    errors.New("value out of range"),
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			u, _ := url.Parse("http://example.com?" + tc.query)
			got := ParseInt(u, tc.expected.Key, tc.opt...)

			if got.Value != tc.expected.Value {
				t.Errorf("ParseInt() .Value: got = %v, want %v",
					got.Value, tc.expected.Value)
			}

			if got.Default != tc.expected.Default {
				t.Errorf("ParseInt() .Default: got = %v, want %v",
					got.Default, tc.expected.Default)
			}

			if got.Min != tc.expected.Min {
				t.Errorf("ParseInt() .Min: got = %v, want %v",
					got.Min, tc.expected.Min)
			}

			if got.Max != tc.expected.Max {
				t.Errorf("ParseInt() .Max: got = %v, want %v",
					got.Max, tc.expected.Max)
			}

			if !reflect.DeepEqual(got.Others, tc.expected.Others) {
				t.Errorf("ParseInt() .Others: got = %v, want %v",
					got.Others, tc.expected.Others)
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

// TestGetInt tests the GetInt function.
func TestGetInt(t *testing.T) {
	key := "age"
	tests := []struct {
		name     string
		query    string
		opt      []int
		expected int
		ok       bool
	}{
		{
			name:     "Simple call",
			query:    "age=18",
			opt:      nil,
			expected: 18,
			ok:       true,
		},
		{
			name:     "Default value",
			query:    "",
			opt:      nil,
			expected: 0,
			ok:       false,
		},
		{
			name:     "With default value",
			query:    "age=18",
			opt:      []int{21},
			expected: 18,
			ok:       true,
		},
		{
			name:     "With default value and empty key",
			query:    "age=",
			opt:      []int{21},
			expected: 21,
			ok:       false,
		},
		{
			name:     "With range",
			query:    "age=18",
			opt:      []int{18, 30},
			expected: 18,
			ok:       true,
		},
		{
			name:     "With inverted range",
			query:    "age=18",
			opt:      []int{30, 18},
			expected: 18,
			ok:       true,
		},
		{
			name:     "With range and empty key",
			query:    "age=",
			opt:      []int{18, 30},
			expected: 18,
			ok:       false,
		},
		{
			name:     "Out of range",
			query:    "age=55",
			opt:      []int{18, 30},
			expected: 18,
			ok:       false,
		},
		{
			name:     "Incorrect value",
			query:    "age=hello",
			opt:      []int{21},
			expected: 21,
			ok:       false,
		},
		{
			name:     "With other values",
			query:    "age=70",
			opt:      []int{20, 20, 30, 50, 70},
			expected: 70,
			ok:       true,
		},
		{
			name:     "With other values - out of range",
			query:    "age=33",
			opt:      []int{20, 20, 30, 50, 70},
			expected: 20,
			ok:       false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			u, _ := url.Parse("http://example.com?" + tc.query)
			got, ok := GetInt(u, key, tc.opt...)

			if got != tc.expected {
				t.Errorf("GetInt() .Value = %v, want %v", got, tc.expected)
			}

			if ok != tc.ok {
				t.Errorf("GetInt() .Ok = %v, want %v", ok, tc.ok)
			}
		})
	}
}

// TestPullInt tests the PullInt function.
func TestPullInt(t *testing.T) {
	key := "age"
	tests := []struct {
		name     string
		query    string
		opt      []int
		expected *int
	}{
		{
			name:     "Simple call",
			query:    "age=18",
			opt:      nil,
			expected: g.Ptr(18),
		},
		{
			name:     "Default value",
			query:    "",
			opt:      nil,
			expected: nil,
		},
		{
			name:     "With default value",
			query:    "age=18",
			opt:      []int{21},
			expected: g.Ptr(18),
		},
		{
			name:     "With default value and empty key",
			query:    "age=",
			opt:      []int{21},
			expected: g.Ptr(21),
		},
		{
			name:     "With range",
			query:    "age=18",
			opt:      []int{18, 30},
			expected: g.Ptr(18),
		},
		{
			name:     "With inverted range",
			query:    "age=18",
			opt:      []int{30, 18},
			expected: g.Ptr(18),
		},
		{
			name:     "With range and empty key",
			query:    "age=",
			opt:      []int{18, 30},
			expected: g.Ptr(18),
		},
		{
			name:     "Out of range",
			query:    "age=55",
			opt:      []int{18, 30},
			expected: g.Ptr(18),
		},
		{
			name:     "Incorrect value",
			query:    "age=hello",
			opt:      []int{21},
			expected: g.Ptr(21),
		},
		{
			name:     "With other values",
			query:    "age=70",
			opt:      []int{20, 20, 30, 50, 70},
			expected: g.Ptr(70),
		},
		{
			name:     "With other values - out of range",
			query:    "age=33",
			opt:      []int{20, 20, 30, 50, 70},
			expected: g.Ptr(20),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			u, _ := url.Parse("http://example.com?" + tc.query)
			got := PullInt(u, key, tc.opt...)

			if got == nil && tc.expected != nil {
				t.Errorf("PullInt() = nil, want %v", tc.expected)
			} else if got != nil && tc.expected == nil {
				t.Errorf("PullInt() = %v, want nil", got)
			} else if got != nil && tc.expected != nil {
				if *got != *tc.expected {
					t.Errorf("PullInt() = %v, want %v", got, tc.expected)
				}
			}
		})
	}
}

// TestParseIntSlice tests the ParseIntSlice function.
func TestParseIntSlice(t *testing.T) {
	tests := []struct {
		name     string
		query    string
		opt      [][]int
		expected *Result[[]int]
	}{
		{
			name:  "Simple call",
			query: "ids=18",
			opt:   nil,
			expected: &Result[[]int]{
				Key:      "ids",
				Value:    []int{18},
				Default:  []int{},
				Empty:    false,
				Contains: true,
				Error:    nil,
			},
		},
		{
			name:  "Slice as single value",
			query: "ids=18,19,20",
			opt:   nil,
			expected: &Result[[]int]{
				Key:      "ids",
				Value:    []int{18, 19, 20},
				Default:  []int{},
				Empty:    false,
				Contains: true,
				Error:    nil,
			},
		},
		{
			name:  "Slice as multiple values",
			query: "ids=18&ids=19&ids=20",
			opt:   nil,
			expected: &Result[[]int]{
				Key:      "ids",
				Value:    []int{18, 19, 20},
				Default:  []int{},
				Empty:    false,
				Contains: true,
				Error:    nil,
			},
		},
		{
			name:  "Empty value",
			query: "ids=",
			opt:   nil,
			expected: &Result[[]int]{
				Key:      "ids",
				Value:    []int{},
				Default:  []int{},
				Empty:    true,
				Contains: true,
				Error:    nil,
			},
		},
		{
			name:  "Declared only",
			query: "ids",
			opt:   nil,
			expected: &Result[[]int]{
				Key:      "ids",
				Value:    []int{},
				Default:  []int{},
				Empty:    true,
				Contains: true,
				Error:    nil,
			},
		},
		{
			name:  "Without key",
			query: "age=18",
			opt:   nil,
			expected: &Result[[]int]{
				Key:      "ids",
				Value:    []int{},
				Default:  []int{},
				Empty:    true,
				Contains: false,
				Error:    nil,
			},
		},
		{
			name:  "Incorrect as single value",
			query: "ids=18,string,20",
			opt:   nil,
			expected: &Result[[]int]{
				Key:      "ids",
				Value:    []int{},
				Default:  []int{},
				Empty:    false,
				Contains: true,
				Error:    errors.New("incorrect value"),
			},
		},
		{
			name:  "Incorrect as multiple values",
			query: "ids=18&ids=string&ids=20",
			opt:   nil,
			expected: &Result[[]int]{
				Key:      "ids",
				Value:    []int{},
				Default:  []int{},
				Empty:    false,
				Contains: true,
				Error:    errors.New("incorrect value"),
			},
		},
		{
			name:  "With default value",
			query: "ids=18&ids=19&ids=20",
			opt:   [][]int{{1, 2, 3}},
			expected: &Result[[]int]{
				Key:      "values",
				Value:    []int{1, 2, 3},
				Default:  []int{1, 2, 3},
				Empty:    true,
				Contains: false,
				Error:    nil,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			u, _ := url.Parse("http://example.com?" + tc.query)
			got := ParseIntSlice(u, tc.expected.Key, tc.opt...)

			if got.Value == nil && tc.expected.Value != nil ||
				got.Value != nil && tc.expected.Value == nil {
				t.Errorf("ParseIntSlice() .Value*: got = %v, want %v",
					got.Value, tc.expected.Value)
			} else if got.Value != nil && tc.expected.Value != nil {
				if !reflect.DeepEqual(got.Value, tc.expected.Value) {
					t.Errorf("ParseIntSlice() .Value: got = %v, want %v",
						got.Value, tc.expected.Value)
				}
			}

			if got.Empty != tc.expected.Empty {
				t.Errorf("ParseIntSlice() .Empty: got = %v, want %v",
					got.Empty, tc.expected.Empty)
			}

			if got.Contains != tc.expected.Contains {
				t.Errorf("ParseIntSlice() .Contains: got = %v, want %v",
					got.Contains, tc.expected.Contains)
			}

			if !reflect.DeepEqual(got.Default, tc.expected.Default) {
				t.Errorf("ParseIntSlice() .Default: got = %v, want %v",
					got.Default, tc.expected.Default)
			}

			if (got.Error != nil && tc.expected.Error == nil) ||
				(got.Error == nil && tc.expected.Error != nil) {
				t.Errorf("ParseIntSlice() .Error: got = %v, want %v",
					got.Error, tc.expected.Error)
			}
		})
	}
}

// TestGetIntSlice tests the GetIntSlice function.
func TestGetIntSlice(t *testing.T) {
	key := "ids"
	tests := []struct {
		name     string
		query    string
		opt      [][]int
		expected []int
		ok       bool
	}{
		{
			name:     "Simple call",
			query:    "ids=18",
			expected: []int{18},
			ok:       true,
		},
		{
			name:     "Slice as single value",
			query:    "ids=18,19,20",
			expected: []int{18, 19, 20},
			ok:       true,
		},
		{
			name:     "Slice as multiple values",
			query:    "ids=18&ids=19&ids=20",
			expected: []int{18, 19, 20},
			ok:       true,
		},
		{
			name:     "Empty value",
			query:    "ids=",
			expected: []int{},
			ok:       false,
		},
		{
			name:     "Default value",
			query:    "ids=",
			opt:      [][]int{{1, 2, 3}},
			expected: []int{1, 2, 3},
			ok:       false,
		},
		{
			name:     "Default value as nil",
			query:    "ids",
			opt:      [][]int{nil},
			expected: nil,
			ok:       false,
		},
		{
			name:     "Declared only",
			query:    "ids",
			expected: []int{},
			ok:       false,
		},
		{
			name:     "Without key",
			query:    "age=18",
			expected: []int{},
			ok:       false,
		},
		{
			name:     "Incorrect as single value",
			query:    "ids=18,string,20",
			expected: []int{},
			ok:       false,
		},
		{
			name:     "Incorrect as multiple values",
			query:    "ids=18&ids=string&ids=20",
			expected: []int{},
			ok:       false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			u, _ := url.Parse("http://example.com?" + tc.query)
			got, ok := GetIntSlice(u, key, tc.opt...)

			if !reflect.DeepEqual(got, tc.expected) {
				t.Errorf("GetIntSlice() = %v, want %v", got, tc.expected)
			}

			if ok != tc.ok {
				t.Errorf("GetIntSlice() .Ok = %v, want %v", ok, tc.ok)
			}
		})
	}
}

// TestPullIntSlice tests the PullIntSlice function.
func TestPullIntSlice(t *testing.T) {
	key := "ids"
	tests := []struct {
		name     string
		query    string
		opt      [][]int
		expected []int
	}{
		{
			name:     "Simple call",
			query:    "ids=18",
			expected: []int{18},
		},
		{
			name:     "Slice as single value",
			query:    "ids=18,19,20",
			expected: []int{18, 19, 20},
		},
		{
			name:     "Slice as multiple values",
			query:    "ids=18&ids=19&ids=20",
			expected: []int{18, 19, 20},
		},
		{
			name:     "Empty value",
			query:    "ids=",
			expected: []int{},
		},
		{
			name:     "Declared only",
			query:    "ids",
			expected: []int{},
		},
		{
			name:     "Without key",
			query:    "age=18",
			expected: nil,
		},
		{
			name:     "Empty value with default value",
			query:    "ids=",
			opt:      [][]int{{1, 2, 3}},
			expected: []int{1, 2, 3},
		},
		{
			name:     "Incorrect as single value",
			query:    "ids=18,string,20",
			expected: []int{},
		},
		{
			name:     "Incorrect as multiple values",
			query:    "ids=18&ids=string&ids=20",
			expected: []int{},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			u, _ := url.Parse("http://example.com?" + tc.query)
			got := PullIntSlice(u, key, tc.opt...)

			if got == nil && tc.expected != nil {
				t.Errorf("PullIntSlice() = nil, want %v", tc.expected)
			} else if got != nil && tc.expected == nil {
				t.Errorf("PullIntSlice() = %v, want nil", got)
			} else if got != nil && tc.expected != nil {
				if !reflect.DeepEqual(got, tc.expected) {
					t.Errorf("PullIntSlice() = %v, want %v", got, tc.expected)
				}
			}
		})
	}
}
