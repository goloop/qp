package qp

import (
	"net/url"
	"testing"
)

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
