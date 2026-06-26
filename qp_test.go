package qp

import (
	"net/url"
	"testing"
)

// mustURL builds a *url.URL from a raw query string ("a=1&b=2"). A nil result
// is returned for the sentinel "<nil>" so tests can exercise the nil-safe path.
func mustURL(t *testing.T, query string) *url.URL {
	t.Helper()
	if query == "<nil>" {
		return nil
	}

	u, err := url.Parse("http://example.com/?" + query)
	if err != nil {
		t.Fatalf("bad test URL %q: %v", query, err)
	}

	return u
}

// TestNilURLNeverPanics guards BUG-02: a nil *url.URL must not crash any entry
// point. Every family is exercised, scalar and slice, holder and top-level.
func TestNilURLNeverPanics(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("panicked on nil URL: %v", r)
		}
	}()

	q := New(nil)

	// Holder reads of every kind.
	_ = q.Int("x", qp_default(7))
	_, _ = q.GetInt("x")
	_ = q.PullInt("x")
	_ = q.ParseInt("x")
	_ = q.IntSlice("x")
	_ = q.Float("x")
	_ = q.String("x", Default("d"))
	_ = q.Bool("x")
	_ = q.StringSlice("x")
	_ = q.Contains("x")
	_ = q.Empty("x")

	// Top-level shortcuts on a nil URL.
	if got := Int(nil, "x", Default(42)); got != 42 {
		t.Fatalf("Int(nil) = %d, want default 42", got)
	}
	if p := PullInt(nil, "x"); p != nil {
		t.Fatalf("PullInt(nil) = %v, want nil", p)
	}
	if Contains(nil, "x") || Empty(nil, "x") {
		t.Fatal("Contains/Empty on nil URL must be false")
	}
}

// qp_default is a tiny alias so the test above reads naturally; it documents
// that Default infers its type from the argument.
func qp_default(v int) Option[int] { return Default(v) }

// TestPresenceMatrix pins the four observable states of a parameter so the
// absent/empty/valid/invalid distinctions never blur again (BUG-03, BUG-04).
func TestPresenceMatrix(t *testing.T) {
	tests := []struct {
		name     string
		query    string
		contains bool
		empty    bool
		hasErr   bool
		value    int
	}{
		{"absent", "other=1", false, true, false, 0},
		{"present empty", "age=", true, true, false, 0},
		{"present valid", "age=25", true, false, false, 25},
		{"present invalid", "age=oops", true, false, true, 0},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			r := New(mustURL(t, tc.query)).ParseInt("age")
			if r.Contains != tc.contains {
				t.Errorf("Contains = %v, want %v", r.Contains, tc.contains)
			}
			if r.Empty != tc.empty {
				t.Errorf("Empty = %v, want %v", r.Empty, tc.empty)
			}
			if (r.Error != nil) != tc.hasErr {
				t.Errorf("Error = %v, want hasErr=%v", r.Error, tc.hasErr)
			}
			if r.Value != tc.value {
				t.Errorf("Value = %d, want %d", r.Value, tc.value)
			}
		})
	}
}

// TestEmptyVsAbsent locks in the fix for BUG-03: Empty distinguishes an absent
// parameter from a present-but-empty one, where the old code conflated them.
func TestEmptyVsAbsent(t *testing.T) {
	tests := []struct {
		name     string
		query    string
		key      string
		contains bool
		empty    bool
	}{
		{"present with value", "foo=bar", "foo", true, false},
		{"present empty", "foo=", "foo", true, true},
		{"absent", "foo=bar", "baz", false, false},
		{"no query at all", "", "foo", false, false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			q := New(mustURL(t, tc.query))
			if got := q.Contains(tc.key); got != tc.contains {
				t.Errorf("Contains = %v, want %v", got, tc.contains)
			}
			if got := q.Empty(tc.key); got != tc.empty {
				t.Errorf("Empty = %v, want %v", got, tc.empty)
			}
		})
	}
}

// TestFromMatchesNew checks that From(u.Query()) and New(u) read identically,
// and that From(nil) is safe.
func TestFromMatchesNew(t *testing.T) {
	u := mustURL(t, "a=1&a=2&b=hello")
	fromNew := New(u)
	fromVals := From(u.Query())

	if a1, a2 := fromNew.IntSlice("a"), fromVals.IntSlice("a"); len(a1) != len(a2) {
		t.Fatalf("IntSlice differs: %v vs %v", a1, a2)
	}
	if New(u).String("b") != From(u.Query()).String("b") {
		t.Fatal("String differs between New and From")
	}
	if From(nil).Contains("b") {
		t.Fatal("From(nil) must have no parameters")
	}
}

// TestQuerySingleDecode is a behavioural proxy for the one-parse design: a
// Query reads from its own snapshot, so mutating the source URL afterwards does
// not change what the holder sees.
func TestQuerySingleDecode(t *testing.T) {
	u := mustURL(t, "n=1")
	q := New(u)
	u.RawQuery = "n=999" // change the source after construction

	if got := q.Int("n"); got != 1 {
		t.Fatalf("Int = %d, want 1 (holder must use its own snapshot)", got)
	}
}
