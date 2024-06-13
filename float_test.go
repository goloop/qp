package qp

import (
	"errors"
	"net/url"
	"reflect"
	"testing"

	"github.com/goloop/g"
)

// TestParseFloat tests the ParseFloat function.
func TestParseFloat(t *testing.T) {
	tests := []struct {
		name     string
		query    string
		opt      []float64
		expected *Result[float64]
	}{
		{
			name:  "Simple call",
			query: "price=12.34",
			opt:   nil,
			expected: &Result[float64]{
				Key:     "price",
				Value:   12.34,
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
			expected: &Result[float64]{
				Key:     "price",
				Value:   0.0,
				Default: 0.0,
				Min:     0.0,
				Max:     0.0,
				Others:  nil,

				Empty:    true,
				Contains: false,
				Error:    nil,
			},
		},
		{
			name:  "With default value",
			query: "price=12.34",
			opt:   []float64{21.0},
			expected: &Result[float64]{
				Key:     "price",
				Value:   12.34,
				Default: 21.0,
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
			query: "price=",
			opt:   []float64{21.1},
			expected: &Result[float64]{
				Key:     "price",
				Value:   21.1,
				Default: 21.1,
				Min:     0.0,
				Max:     0.0,
				Others:  nil,

				Empty:    true,
				Contains: true,
				Error:    nil,
			},
		},
		{
			name:  "With range",
			query: "price=21.0",
			opt:   []float64{18.0, 30.0},
			expected: &Result[float64]{
				Key:     "price",
				Value:   21.0,
				Default: 18.0,
				Min:     18.0,
				Max:     30.0,
				Others:  nil,

				Empty:    false,
				Contains: true,
				Error:    nil,
			},
		},
		{
			name:  "With inverted range",
			query: "price=21.0",
			opt:   []float64{30.0, 18.0},
			expected: &Result[float64]{
				Key:     "price",
				Value:   21.0,
				Default: 30.0,
				Min:     18.0,
				Max:     30.0,
				Others:  nil,

				Empty:    false,
				Contains: true,
				Error:    nil,
			},
		},
		{
			name:  "With range and empty key",
			query: "price", // without =
			opt:   []float64{18.0, 30.0},
			expected: &Result[float64]{
				Key:     "price",
				Value:   18.0,
				Default: 18.0,
				Min:     18.0,
				Max:     30.0,
				Others:  nil,

				Empty:    true,
				Contains: true,
				Error:    nil,
			},
		},
		{
			name:  "Out of range",
			query: "price=55",
			opt:   []float64{18.0, 30.0},
			expected: &Result[float64]{
				Key:     "price",
				Value:   18.0,
				Default: 18.0,
				Min:     18.0,
				Max:     30.0,
				Others:  nil,

				Empty:    false,
				Contains: true,
				Error:    errors.New("value out of range"),
			},
		},
		{
			name:  "Incorrect value",
			query: "price=hello",
			opt:   []float64{21},
			expected: &Result[float64]{
				Key:     "price",
				Value:   21.0,
				Default: 21.0,
				Min:     0.0,
				Max:     0.0,
				Others:  nil,

				Empty:    false,
				Contains: true,
				Error:    errors.New("incorrect value"),
			},
		},
		{
			name:  "With other values",
			query: "price=70",
			opt:   []float64{20.0, 20.0, 30.0, 50.0, 70.0},
			expected: &Result[float64]{
				Key:     "price",
				Value:   70.0,
				Default: 20.0,
				Min:     20.0,
				Max:     20.0,
				Others:  []float64{30.0, 50.0, 70.0},

				Empty:    false,
				Contains: true,
				Error:    nil,
			},
		},
		{
			name:  "With other values - out of range",
			query: "price=33",
			opt:   []float64{20.0, 20.0, 30.0, 50.0, 70.0},
			expected: &Result[float64]{
				Key:     "price",
				Value:   20,
				Default: 20,
				Min:     20,
				Max:     20,
				Others:  []float64{30.0, 50.0, 70.0},

				Empty:    false,
				Contains: true,
				Error:    errors.New("value out of range"),
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			u, _ := url.Parse("http://example.com?" + tc.query)
			got := ParseFloat(u, tc.expected.Key, tc.opt...)

			if got.Value != tc.expected.Value {
				t.Errorf("ParseFloat() .Value: got = %v, want %v",
					got.Value, tc.expected.Value)
			}

			if got.Default != tc.expected.Default {
				t.Errorf("ParseFloat() .Default: got = %v, want %v",
					got.Default, tc.expected.Default)
			}

			if got.Min != tc.expected.Min {
				t.Errorf("ParseFloat() .Min: got = %v, want %v",
					got.Min, tc.expected.Min)
			}

			if got.Max != tc.expected.Max {
				t.Errorf("ParseFloat() .Max: got = %v, want %v",
					got.Max, tc.expected.Max)
			}

			if !reflect.DeepEqual(got.Others, tc.expected.Others) {
				t.Errorf("ParseFloat() .Others: got = %v, want %v",
					got.Others, tc.expected.Others)
			}

			if got.Empty != tc.expected.Empty {
				t.Errorf("ParseFloat() .Empty: got = %v, want %v",
					got.Empty, tc.expected.Empty)
			}

			if got.Contains != tc.expected.Contains {
				t.Errorf("ParseFloat() .Contains: got = %v, want %v",
					got.Contains, tc.expected.Contains)
			}

			if (got.Error != nil && tc.expected.Error == nil) ||
				(got.Error == nil && tc.expected.Error != nil) {
				t.Errorf("ParseFloat() .Error: got = %v, want %v",
					got.Error, tc.expected.Error)
			}
		})
	}
}

// TestGetFloat tests the GetFloat function.
func TestGetFloat(t *testing.T) {
	key := "price"
	tests := []struct {
		name     string
		query    string
		opt      []float64
		expected float64
		ok       bool
	}{
		{
			name:     "Simple call",
			query:    "price=18",
			opt:      nil,
			expected: 18.0,
			ok:       true,
		},
		{
			name:     "Default value",
			query:    "",
			opt:      nil,
			expected: 0.0,
			ok:       false,
		},
		{
			name:     "With default value",
			query:    "price=18.0",
			opt:      []float64{21.0},
			expected: 18.0,
			ok:       true,
		},
		{
			name:     "With default value and empty key",
			query:    "price=",
			opt:      []float64{21.0},
			expected: 21.0,
			ok:       false,
		},
		{
			name:     "With range",
			query:    "price=21.0",
			opt:      []float64{18.0, 30.0},
			expected: 21.0,
			ok:       true,
		},
		{
			name:     "With inverted range",
			query:    "price=21.0",
			opt:      []float64{30.0, 18.0},
			expected: 21.0,
			ok:       true,
		},
		{
			name:     "With range and empty key",
			query:    "price",
			opt:      []float64{18.0, 30.0},
			expected: 18.0,
			ok:       false,
		},
		{
			name:     "Out of range",
			query:    "price=55.0",
			opt:      []float64{18.0, 30.0},
			expected: 18.0,
			ok:       false,
		},
		{
			name:     "Incorrect value",
			query:    "price=hello",
			opt:      []float64{21.0},
			expected: 21.0,
			ok:       false,
		},
		{
			name:     "With other values",
			query:    "price=70.0",
			opt:      []float64{20.0, 20.0, 30.0, 50.0, 70.0},
			expected: 70.0,
			ok:       true,
		},
		{
			name:     "With other values - out of range",
			query:    "price=33.0",
			opt:      []float64{20.0, 20.0, 30.0, 50.0, 70.0},
			expected: 20.0,
			ok:       false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			u, _ := url.Parse("http://example.com?" + tc.query)
			got, ok := GetFloat(u, key, tc.opt...)

			if got != tc.expected {
				t.Errorf("GetFloat() .Value = %v, want %v", got, tc.expected)
			}

			if ok != tc.ok {
				t.Errorf("GetFloat() .Ok = %v, want %v", ok, tc.ok)
			}
		})
	}
}

// TestPullFloat tests the PullFloat function.
func TestPullFloat(t *testing.T) {
	key := "price"
	tests := []struct {
		name     string
		query    string
		opt      []float64
		expected *float64
	}{
		{
			name:     "Simple call",
			query:    "price=18.0",
			opt:      nil,
			expected: g.Ptr(18.0),
		},
		{
			name:     "Default value",
			query:    "",
			opt:      nil,
			expected: nil,
		},
		{
			name:     "With default value",
			query:    "price=18.0",
			opt:      []float64{21.0},
			expected: g.Ptr(18.0),
		},
		{
			name:     "With default value and empty key",
			query:    "price=",
			opt:      []float64{21.0},
			expected: g.Ptr(21.0),
		},
		{
			name:     "With range",
			query:    "price=18.5",
			opt:      []float64{18.0, 30.0},
			expected: g.Ptr(18.5),
		},
		{
			name:     "With inverted range",
			query:    "price=18.3",
			opt:      []float64{30.0, 18.0},
			expected: g.Ptr(18.3),
		},
		{
			name:     "With range and empty key",
			query:    "price=",
			opt:      []float64{18.0, 30.0},
			expected: g.Ptr(18.0),
		},
		{
			name:     "Out of range",
			query:    "price=55",
			opt:      []float64{18.0, 30.0},
			expected: g.Ptr(18.0),
		},
		{
			name:     "Incorrect value",
			query:    "price=hello",
			opt:      []float64{21.0},
			expected: g.Ptr(21.0),
		},
		{
			name:     "With other values",
			query:    "price=70",
			opt:      []float64{20.0, 20.0, 30.0, 50.0, 70.0},
			expected: g.Ptr(70.0),
		},
		{
			name:     "With other values - out of range",
			query:    "price=33",
			opt:      []float64{20.0, 20.0, 30.0, 50.0, 70.0},
			expected: g.Ptr(20.0),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			u, _ := url.Parse("http://example.com?" + tc.query)
			got := PullFloat(u, key, tc.opt...)

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

// TestParseFloatSlice tests the ParseFloatSlice function.
func TestParseFloatSlice(t *testing.T) {
	tests := []struct {
		name     string
		query    string
		opt      [][]float64
		expected *Result[[]float64]
	}{
		{
			name:  "Simple call",
			query: "values=18.5",
			expected: &Result[[]float64]{
				Key:      "values",
				Value:    []float64{18.5},
				Default:  []float64{},
				Empty:    false,
				Contains: true,
				Error:    nil,
			},
		},
		{
			name:  "Slice as single value",
			query: "values=18.5,19.5,20.5",
			expected: &Result[[]float64]{
				Key:      "values",
				Value:    []float64{18.5, 19.5, 20.5},
				Default:  []float64{},
				Empty:    false,
				Contains: true,
				Error:    nil,
			},
		},
		{
			name:  "Slice as multiple values",
			query: "values=18.5&values=19.5&values=20.5",
			expected: &Result[[]float64]{
				Key:      "values",
				Value:    []float64{18.5, 19.5, 20.5},
				Default:  []float64{},
				Contains: true,
				Error:    nil,
			},
		},
		{
			name:  "Empty value",
			query: "values=",
			expected: &Result[[]float64]{
				Key:      "values",
				Value:    []float64{},
				Default:  []float64{},
				Empty:    true,
				Contains: true,
				Error:    nil,
			},
		},
		{
			name:  "Declared only",
			query: "values",
			expected: &Result[[]float64]{
				Key:      "values",
				Value:    []float64{},
				Default:  []float64{},
				Empty:    true,
				Contains: true,
				Error:    nil,
			},
		},
		{
			name:  "Without key",
			query: "age=18.5",
			expected: &Result[[]float64]{
				Key:      "values",
				Value:    []float64{},
				Default:  []float64{},
				Empty:    true,
				Contains: false,
				Error:    nil,
			},
		},
		{
			name:  "Without key with default value",
			query: "age=18.5",
			opt:   [][]float64{{1.0, 2.0, 3.0}},
			expected: &Result[[]float64]{
				Key:      "values",
				Value:    []float64{1.0, 2.0, 3.0},
				Default:  []float64{1.0, 2.0, 3.0},
				Empty:    true,
				Contains: false,
				Error:    nil,
			},
		},
		{
			name:  "Incorrect as single value",
			query: "ids=18.0,string,20.0",
			expected: &Result[[]float64]{
				Key:      "ids",
				Value:    []float64{},
				Default:  []float64{},
				Empty:    false,
				Contains: true,
				Error:    errors.New("incorrect value"),
			},
		},
		{
			name:  "Incorrect as multiple values",
			query: "ids=18.0&ids=string&ids=20.0",
			expected: &Result[[]float64]{
				Key:      "ids",
				Value:    []float64{},
				Default:  []float64{},
				Empty:    false,
				Contains: true,
				Error:    errors.New("incorrect value"),
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			u, _ := url.Parse("http://example.com?" + tc.query)
			got := ParseFloatSlice(u, tc.expected.Key, tc.opt...)

			if got.Value == nil && tc.expected.Value != nil ||
				got.Value != nil && tc.expected.Value == nil {
				t.Errorf("ParseFloatSlice()* .Value*: got = %v, want %v",
					got.Value, tc.expected.Value)
			} else if got.Value != nil && tc.expected.Value != nil {
				if !reflect.DeepEqual(got.Value, tc.expected.Value) {
					t.Errorf("ParseFloatSlice() .Value: got = %v, want %v",
						got.Value, tc.expected.Value)
				}
			}

			if !reflect.DeepEqual(got.Default, tc.expected.Default) {
				t.Errorf("ParseFloatSlice() .Default: got = %v, want %v",
					got.Default, tc.expected.Default)
			}

			if got.Empty != tc.expected.Empty {
				t.Errorf("ParseFloatSlice() .Empty: got = %v, want %v",
					got.Empty, tc.expected.Empty)
			}

			if got.Contains != tc.expected.Contains {
				t.Errorf("ParseFloatSlice() .Contains: got = %v, want %v",
					got.Contains, tc.expected.Contains)
			}

			if (got.Error != nil && tc.expected.Error == nil) ||
				(got.Error == nil && tc.expected.Error != nil) {
				t.Errorf("ParseFloatSlice() .Error: got = %v, want %v",
					got.Error, tc.expected.Error)
			}
		})
	}
}

// TestGetFloatSlice tests the GetIntSlice function.
func TestGetFloatSlice(t *testing.T) {
	key := "ids"
	tests := []struct {
		name     string
		query    string
		opt      [][]float64
		expected []float64
		ok       bool
	}{
		{
			name:     "Simple call",
			query:    "ids=18",
			expected: []float64{18.0},
			ok:       true,
		},
		{
			name:     "Slice as single value",
			query:    "ids=18,19,20",
			expected: []float64{18.0, 19.0, 20.0},
			ok:       true,
		},
		{
			name:     "Slice as multiple values",
			query:    "ids=18&ids=19&ids=20",
			expected: []float64{18.0, 19.0, 20.0},
			ok:       true,
		},
		{
			name:     "Empty value",
			query:    "ids=",
			expected: []float64{},
			ok:       false,
		},
		{
			name:     "Declared only",
			query:    "ids",
			expected: []float64{},
			ok:       false,
		},
		{
			name:     "Without key",
			query:    "age=18",
			expected: []float64{},
			ok:       false,
		},
		{
			name:     "Without key with default value",
			query:    "age=18",
			opt:      [][]float64{{1.0, 2.0, 3.0}},
			expected: []float64{1.0, 2.0, 3.0},
			ok:       false,
		},
		{
			name:     "Incorrect as single value",
			query:    "ids=18,string,20",
			expected: []float64{},
			ok:       false,
		},
		{
			name:     "Incorrect as multiple values",
			query:    "ids=18&ids=string&ids=20",
			expected: []float64{},
			ok:       false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			u, _ := url.Parse("http://example.com?" + tc.query)
			got, ok := GetFloatSlice(u, key, tc.opt...)

			if !reflect.DeepEqual(got, tc.expected) {
				t.Errorf("GetFloatSlice() = %v, want %v", got, tc.expected)
			}

			if ok != tc.ok {
				t.Errorf("GetFloatSlice() .Ok = %v, want %v", ok, tc.ok)
			}
		})
	}
}

// TestPullFloatSlice tests the PullFloatSlice function.
func TestPullFloatSlice(t *testing.T) {
	key := "ids"
	tests := []struct {
		name     string
		query    string
		expected []float64
	}{
		{
			name:     "Simple call",
			query:    "ids=18.5",
			expected: []float64{18.5},
		},
		{
			name:     "Slice as single value",
			query:    "ids=18.5,19.5,20.5",
			expected: []float64{18.5, 19.5, 20.5},
		},
		{
			name:     "Slice as multiple ids",
			query:    "ids=18.5&ids=19.5&ids=20.5",
			expected: []float64{18.5, 19.5, 20.5},
		},
		{
			name:     "Empty value",
			query:    "ids=",
			expected: []float64{},
		},
		{
			name:     "Declared only",
			query:    "ids",
			expected: []float64{},
		},
		{
			name:     "Without key",
			query:    "age=18.5",
			expected: nil,
		},
		{
			name:     "Incorrect as single value",
			query:    "ids=18.0,string,20.0",
			expected: []float64{},
		},
		{
			name:     "Incorrect as multiple values",
			query:    "ids=18.0&ids=string&ids=20.0",
			expected: []float64{},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			u, _ := url.Parse("http://example.com?" + tc.query)
			got := PullFloatSlice(u, key)

			if got == nil && tc.expected != nil {
				t.Errorf("PullFloatSlice() = nil, want %v", tc.expected)
			} else if got != nil && tc.expected == nil {
				t.Errorf("PullFloatSlice() = %v, want nil", got)
			} else if got != nil && tc.expected != nil {
				if !reflect.DeepEqual(got, tc.expected) {
					t.Errorf("PullFloatSlice() = %v, want %v", got, tc.expected)
				}
			}
		})
	}
}
