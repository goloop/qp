package qp

import (
	"strings"
	"testing"
)

// TestOutOfRangeContract pins the agreed contract (BUG-04): a present value that
// violates a constraint yields Value=default AND a non-nil Error; Get reports
// ok=false; Pull still returns a pointer to the default. The failure is never
// silently swallowed and never silently accepted.
func TestOutOfRangeContract(t *testing.T) {
	u := mustURL(t, "age=99")

	r := ParseInt(u, "age", Default(18), Between(18, 65))
	if r.Value != 18 {
		t.Errorf("Value = %d, want default 18", r.Value)
	}
	if r.Error == nil {
		t.Error("Error = nil, want non-nil for out-of-range value")
	}

	if v, ok := GetInt(u, "age", Default(18), Between(18, 65)); v != 18 || ok {
		t.Errorf("GetInt = (%d, %v), want (18, false)", v, ok)
	}

	p := PullInt(u, "age", Default(18), Between(18, 65))
	if p == nil || *p != 18 {
		t.Errorf("PullInt = %v, want pointer to 18", p)
	}
}

// TestNoDefaultFallsBackToZero documents that without Default an invalid value
// falls back to the zero value, not to some bound.
func TestNoDefaultFallsBackToZero(t *testing.T) {
	u := mustURL(t, "age=99")
	r := ParseInt(u, "age", Between(18, 65))
	if r.Value != 0 || r.Error == nil {
		t.Errorf("Value=%d Error=%v, want 0 and non-nil error", r.Value, r.Error)
	}
}

// TestRangeBoundsInclusive checks the inclusive edges and the open-ended Min/Max
// forms, which the old API could not express.
func TestRangeBoundsInclusive(t *testing.T) {
	tests := []struct {
		name  string
		query string
		opts  []Option[int]
		valid bool
	}{
		{"low edge in", "n=18", []Option[int]{Between(18, 65)}, true},
		{"high edge in", "n=65", []Option[int]{Between(18, 65)}, true},
		{"below", "n=17", []Option[int]{Between(18, 65)}, false},
		{"above", "n=66", []Option[int]{Between(18, 65)}, false},
		{"min only ok", "n=1000", []Option[int]{Min(1)}, true},
		{"min only bad", "n=0", []Option[int]{Min(1)}, false},
		{"max only ok", "n=-5", []Option[int]{Max(0)}, true},
		{"max only bad", "n=1", []Option[int]{Max(0)}, false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, ok := GetInt(mustURL(t, tc.query), "n", tc.opts...)
			if ok != tc.valid {
				t.Errorf("ok = %v, want %v", ok, tc.valid)
			}
		})
	}
}

// TestBetweenSwapsBounds verifies Between(hi, lo) behaves like Between(lo, hi).
func TestBetweenSwapsBounds(t *testing.T) {
	u := mustURL(t, "n=40")
	if _, ok := GetInt(u, "n", Between(65, 18)); !ok {
		t.Error("Between(65, 18) should accept 40 after swapping bounds")
	}
}

// TestOneOfAsRangeException locks the union semantics: with both a range and a
// set, a value outside the range is still accepted if it is in the set.
func TestOneOfAsRangeException(t *testing.T) {
	u := mustURL(t, "n=70")
	if _, ok := GetInt(u, "n", Between(18, 65), OneOf(16, 70)); !ok {
		t.Error("70 should be accepted as an explicit OneOf exception")
	}
	if _, ok := GetInt(mustURL(t, "n=99"), "n", Between(18, 65), OneOf(16, 70)); ok {
		t.Error("99 is neither in range nor in the set; want ok=false")
	}
}

// TestStringOneOf is the string-validation path (replaces the old fragile
// "all-args-are-allowed-values" rule). The error must name the value readably,
// guarding BUG-01 (no %!d(string=...) corruption).
func TestStringOneOf(t *testing.T) {
	u := mustURL(t, "role=hacker")

	r := ParseString(u, "role", Default("guest"), OneOf("admin", "user"))
	if r.Value != "guest" {
		t.Errorf("Value = %q, want default \"guest\"", r.Value)
	}
	if r.Error == nil {
		t.Fatal("Error = nil, want non-nil for disallowed value")
	}
	if msg := r.Error.Error(); !strings.Contains(msg, "hacker") ||
		strings.Contains(msg, "%!") {
		t.Errorf("error message %q must contain the value and no format glitch", msg)
	}

	if v, ok := GetString(mustURL(t, "role=admin"), "role",
		OneOf("admin", "user")); v != "admin" || !ok {
		t.Errorf("GetString allowed value = (%q, %v), want (admin, true)", v, ok)
	}
}

// TestParseErrorMessageQuoted ensures the parse-failure message is well-formed
// for every scalar type and never leaks a wrong format verb.
func TestParseErrorMessageQuoted(t *testing.T) {
	cases := []struct {
		name string
		err  error
	}{
		{"int", ParseInt(mustURL(t, "x=abc"), "x").Error},
		{"float", ParseFloat(mustURL(t, "x=abc"), "x").Error},
		{"bool", ParseBool(mustURL(t, "x=maybe"), "x").Error},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.err == nil {
				t.Fatal("want non-nil error")
			}
			if strings.Contains(tc.err.Error(), "%!") {
				t.Errorf("malformed error: %q", tc.err.Error())
			}
		})
	}
}
