package qp

import (
	"net/url"
	"testing"
)

func TestClampFoldsIntoRange(t *testing.T) {
	q := New(&url.URL{RawQuery: "page=999999999"})
	// Without a ceiling this value would overflow arithmetic downstream; Clamp
	// returns the top bound instead.
	if got := q.Int("page", Default(1), Clamp(1, 1000)); got != 1000 {
		t.Errorf("page = %d, want 1000 (clamped to hi)", got)
	}
}

func TestClampReturnsLowBoundBelowRange(t *testing.T) {
	q := New(&url.URL{RawQuery: "n=-50"})
	if got := q.Int("n", Default(7), Clamp(1, 100)); got != 1 {
		t.Errorf("n = %d, want 1 (clamped to lo)", got)
	}
}

func TestClampLeavesInRangeValueAlone(t *testing.T) {
	q := New(&url.URL{RawQuery: "page=42"})
	if got := q.Int("page", Default(1), Clamp(1, 1000)); got != 42 {
		t.Errorf("page = %d, want 42 (unchanged)", got)
	}
}

func TestClampSwapsReversedBounds(t *testing.T) {
	q := New(&url.URL{RawQuery: "page=999999999"})
	if got := q.Int("page", Default(1), Clamp(1000, 1)); got != 1000 {
		t.Errorf("page = %d, want 1000 (bounds swapped)", got)
	}
}

func TestClampAbsentUsesDefault(t *testing.T) {
	q := New(&url.URL{RawQuery: ""})
	if got := q.Int("page", Default(1), Clamp(1, 1000)); got != 1 {
		t.Errorf("page = %d, want the default of 1", got)
	}
}
