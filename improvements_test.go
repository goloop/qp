package qp

import (
	"errors"
	"net/url"
	"reflect"
	"strings"
	"testing"
)

// TestSliceLeadingEmptyPosition guards BUG-01: an empty value among several
// repeated values must not discard the parameter, regardless of where it sits.
// The old engine bailed out whenever the first value was empty.
func TestSliceLeadingEmptyPosition(t *testing.T) {
	tests := []struct {
		query string
		want  []int
	}{
		{"ids=5&ids=", []int{5}},          // trailing empty
		{"ids=5&ids=&ids=6", []int{5, 6}}, // middle empty
		{"ids=&ids=5", []int{5}},          // leading empty (was the bug)
		{"ids=&ids=&ids=7", []int{7}},     // two leading empties (was the bug)
		{"ids=&ids=5&ids=", []int{5}},     // empties on both ends
	}

	for _, tc := range tests {
		t.Run(tc.query, func(t *testing.T) {
			got := IntSlice(mustURL(t, tc.query), "ids")
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("got %v, want %v", got, tc.want)
			}
		})
	}
}

// TestSinglePresentEmptySliceStillEmpty makes sure the BUG-01 fix did not change
// the canonical present-empty case: "?ids=" (a single empty value) is still
// reported Empty with the fallback, not run through the split policy.
func TestSinglePresentEmptySliceStillEmpty(t *testing.T) {
	r := ParseIntSlice(mustURL(t, "ids="), "ids")
	if !r.Contains || !r.Empty || r.Error != nil {
		t.Errorf("?ids= -> Contains=%v Empty=%v Error=%v, want true/true/nil",
			r.Contains, r.Empty, r.Error)
	}
	if len(r.Value) != 0 {
		t.Errorf("Value = %v, want empty slice", r.Value)
	}
}

// TestFromEmptyValueSliceIsEmpty guards the slice engine against a key whose
// url.Values entry is a zero-length slice (buildable via From): it must be
// reported as present-but-empty, matching the scalar engine and Query.Empty.
func TestFromEmptyValueSliceIsEmpty(t *testing.T) {
	q := From(url.Values{"k": {}})
	r := q.ParseIntSlice("k")
	if !r.Contains || !r.Empty || r.Error != nil {
		t.Errorf("From empty slice -> Contains=%v Empty=%v Error=%v, want true/true/nil",
			r.Contains, r.Empty, r.Error)
	}
	if !q.Empty("k") {
		t.Error("Query.Empty(k) = false, want true for a present-but-empty key")
	}
}

// TestEmptyOneOfRejectsAll guards BUG-02: an explicitly empty allowed set means
// "nothing is allowed", so every value is rejected — a dynamically-empty
// OneOf(allowed...) must not silently disable validation.
func TestEmptyOneOfRejectsAll(t *testing.T) {
	var noStrings []string
	if v, ok := GetString(mustURL(t, "role=anything"), "role",
		OneOf(noStrings...)); ok {
		t.Errorf("empty OneOf accepted %q (ok=%v), want reject", v, ok)
	}

	var noInts []int
	if _, ok := GetInt(mustURL(t, "n=5"), "n", OneOf(noInts...)); ok {
		t.Error("empty OneOf accepted an int, want reject")
	}

	// A non-empty OneOf still works as before.
	if _, ok := GetString(mustURL(t, "role=admin"), "role",
		OneOf("admin")); !ok {
		t.Error("non-empty OneOf must still accept a listed value")
	}
}

// TestFallbackIsCloned guards BUG-03: the fallback slice must not be aliased.
// Mutating one read's result must not corrupt the caller's slice or another
// read that uses the same fallback.
func TestFallbackIsCloned(t *testing.T) {
	fb := []int{1, 2, 3}
	q := New(mustURL(t, "other=1"))

	a := q.IntSlice("missing1", Fallback(fb)) // absent -> fallback
	b := q.IntSlice("missing2", Fallback(fb)) // absent -> fallback
	a[0] = 99

	if fb[0] != 1 {
		t.Errorf("caller's slice mutated: fb = %v", fb)
	}
	if b[0] != 1 {
		t.Errorf("second read shares backing array: b = %v", b)
	}

	// Same guarantee on the invalid path, where the fallback is also returned.
	c := IntSlice(mustURL(t, "ids=x"), "ids", Fallback(fb))
	c[0] = 77
	if fb[0] != 1 {
		t.Errorf("invalid-path fallback aliased caller's slice: fb = %v", fb)
	}
}

// TestValuesIsCopy guards the Values() footgun: the returned map is independent,
// so mutating it never changes what the Query reads.
func TestValuesIsCopy(t *testing.T) {
	q := New(mustURL(t, "a=1&list=x"))

	v := q.Values()
	v.Set("a", "999")
	v.Add("b", "injected")
	v["list"][0] = "tampered"

	if got := q.Int("a"); got != 1 {
		t.Errorf("a = %d after external mutation, want 1", got)
	}
	if q.Contains("b") {
		t.Error("Query gained key b from a mutated copy")
	}
	if got := q.String("list"); got != "x" {
		t.Errorf("list = %q after slice mutation, want \"x\"", got)
	}
}

// TestTypedErrors checks that every failure mode wraps the right sentinel and
// carries a ParseError with the offending key and value.
func TestTypedErrors(t *testing.T) {
	invalid := ParseInt(mustURL(t, "age=abc"), "age").Error
	if !errors.Is(invalid, ErrInvalid) {
		t.Errorf("parse failure should wrap ErrInvalid, got %v", invalid)
	}

	notAllowed := ParseInt(mustURL(t, "age=99"), "age", Between(0, 10)).Error
	if !errors.Is(notAllowed, ErrNotAllowed) {
		t.Errorf("range failure should wrap ErrNotAllowed, got %v", notAllowed)
	}

	emptyElem := ParseIntSlice(mustURL(t, "ids=1,,2"), "ids",
		ErrorOnEmpty[int]()).Error
	if !errors.Is(emptyElem, ErrEmptyElement) {
		t.Errorf("empty element should wrap ErrEmptyElement, got %v", emptyElem)
	}

	var pe *ParseError
	if !errors.As(invalid, &pe) {
		t.Fatal("error should unwrap to *ParseError")
	}
	if pe.Key != "age" || pe.Value != "abc" {
		t.Errorf("ParseError = {Key:%q Value:%q}, want {age abc}", pe.Key, pe.Value)
	}
}

// TestErrorMessages pins the human-readable message for each failure kind:
// key and value are quoted, and no format verb ever leaks (guards the old
// %!d(string=...) class of bug).
func TestErrorMessages(t *testing.T) {
	cases := []struct {
		name     string
		err      error
		contains []string
	}{
		{
			"invalid",
			ParseInt(mustURL(t, "age=abc"), "age").Error,
			[]string{"invalid value", `"age"`, `"abc"`},
		},
		{
			"not allowed",
			ParseString(mustURL(t, "role=hacker"), "role", OneOf("admin")).Error,
			[]string{"not allowed", `"role"`, `"hacker"`},
		},
		{
			"empty element",
			ParseIntSlice(mustURL(t, "ids=1,,2"), "ids", ErrorOnEmpty[int]()).Error,
			[]string{"empty element", `"ids"`},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.err == nil {
				t.Fatal("want non-nil error")
			}
			msg := tc.err.Error()
			if strings.Contains(msg, "%!") {
				t.Fatalf("message has a leaked format verb: %q", msg)
			}
			for _, want := range tc.contains {
				if !strings.Contains(msg, want) {
					t.Errorf("message %q missing %q", msg, want)
				}
			}
		})
	}
}
