package qp

import (
	"net/url"
	"testing"
)

// TestContains tests the Contains function.
func TestContains(t *testing.T) {
	tests := []struct {
		name     string
		query    string
		param    string
		expected bool
	}{
		{"Parameter present", "foo=bar", "foo", true},
		{"Parameter absent", "foo=bar", "baz", false},
		{"Empty parameter present", "foo=", "foo", true},
		{"No parameters", "", "foo", false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			reqURL := "http://example.com?" + tc.query
			u, _ := url.Parse(reqURL)

			result := Contains(u, tc.param)
			if result != tc.expected {
				t.Errorf("expected %v, got %v", tc.expected, result)
			}
		})
	}
}

// TestEmpty tests the Empty function.
func TestEmpty(t *testing.T) {
	tests := []struct {
		name     string
		query    string
		param    string
		expected bool
	}{
		{"Empty parameter", "foo=", "foo", true},
		{"Non-empty parameter", "foo=bar", "foo", false},
		{"Parameter absent", "", "foo", true},
		{"Multiple parameters with empty", "foo=bar&baz=", "baz", true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			reqURL := "http://example.com?" + tc.query
			u, _ := url.Parse(reqURL)

			result := Empty(u, tc.param)
			if result != tc.expected {
				t.Errorf("expected %v, got %v", tc.expected, result)
			}
		})
	}
}

// TestParseBoolValue tests the parseBoolValue function.
func TestParseBoolValue(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
		hasError bool
	}{
		{"1", true, false},
		{"true", true, false},
		{"yes", true, false},
		{"on", true, false},
		{"0", false, false},
		{"false", false, false},
		{"no", false, false},
		{"off", false, false},
		{"True", true, false},    // case insensitive
		{"FALSE", false, false},  // case insensitive
		{"invalid", false, true}, // invalid input
		{"2", false, true},       // invalid input
		{"", false, true},        // empty input
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result, err := parseBoolValue(tt.input)
			if (err != nil) != tt.hasError {
				t.Errorf("expected error: %v, got: %v", tt.hasError, err != nil)
			}
			if result != tt.expected {
				t.Errorf("expected: %v, got: %v", tt.expected, result)
			}
		})
	}
}
