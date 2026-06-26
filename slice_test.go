package qp

import (
	"reflect"
	"testing"
)

// TestSliceCommaSplitDefault covers the agreed default policy (BUG-06): every
// element is trimmed and empty elements are dropped, so real-world noise from
// JS array joins and hand-typed URLs parses predictably.
func TestSliceCommaSplitDefault(t *testing.T) {
	tests := []struct {
		name  string
		query string
		want  []int
	}{
		{"plain", "ids=1,2,3", []int{1, 2, 3}},
		{"spaces around commas", "ids=%201,%202%20,3", []int{1, 2, 3}},
		{"doubled comma", "ids=1,,2", []int{1, 2}},
		{"leading comma", "ids=,1,2", []int{1, 2}},
		{"trailing comma", "ids=1,2,", []int{1, 2}},
		{"only commas", "ids=,,,", []int{}},
		{"single value", "ids=42", []int{42}},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, ok := GetIntSlice(mustURL(t, tc.query), "ids")
			if !ok {
				t.Fatalf("ok=false for %q", tc.query)
			}
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("got %v, want %v", got, tc.want)
			}
		})
	}
}

// TestSliceMixedAndRepeatedForms guards BUG-09: the comma form, the repeated
// form, and a mixture of the two must all yield the same slice. The old code
// split only single values, so "?ids=1,2&ids=3" used to error out.
func TestSliceMixedAndRepeatedForms(t *testing.T) {
	want := []int{1, 2, 3}
	queries := []string{
		"ids=1,2,3",         // single comma value
		"ids=1&ids=2&ids=3", // repeated values
		"ids=1,2&ids=3",     // mixed: a comma value plus a repeated one
		"ids=1&ids=2,3",     // mixed the other way
	}

	for _, query := range queries {
		got, ok := GetIntSlice(mustURL(t, query), "ids")
		if !ok || !reflect.DeepEqual(got, want) {
			t.Errorf("%q: got %v ok=%v, want %v", query, got, ok, want)
		}
	}
}

// TestStringSliceMixedForms is the string analogue of BUG-09: "?n=a,b&n=c" must
// split the comma element too, instead of silently keeping "a,b" whole.
func TestStringSliceMixedForms(t *testing.T) {
	got := StringSlice(mustURL(t, "n=a,b&n=c"), "n")
	want := []string{"a", "b", "c"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

// TestSliceKeepEmpty checks the opt-in KeepEmpty policy: strings preserve the
// empty slots, while numeric slices then fail to parse them.
func TestSliceKeepEmpty(t *testing.T) {
	names := StringSlice(mustURL(t, "n=a,,b"), "n", KeepEmpty[string]())
	if want := []string{"a", "", "b"}; !reflect.DeepEqual(names, want) {
		t.Errorf("KeepEmpty strings = %v, want %v", names, want)
	}

	r := ParseIntSlice(mustURL(t, "ids=1,,2"), "ids", KeepEmpty[int]())
	if r.Error == nil {
		t.Error("KeepEmpty int with empty element should error")
	}
}

// TestSliceErrorOnEmpty checks the strict policy: any empty element rejects the
// whole slice, even a harmless trailing comma.
func TestSliceErrorOnEmpty(t *testing.T) {
	for _, query := range []string{"ids=1,,2", "ids=1,2,", "ids=,1"} {
		r := ParseIntSlice(mustURL(t, query), "ids", ErrorOnEmpty[int]())
		if r.Error == nil {
			t.Errorf("%q: want error under ErrorOnEmpty", query)
		}
		if len(r.Value) != 0 {
			t.Errorf("%q: Value = %v, want empty fallback on error", query, r.Value)
		}
	}
}

// TestSliceEachValidation guards BUG-05: element validation is now expressible.
// One out-of-range element rejects the whole slice and yields the fallback.
func TestSliceEachValidation(t *testing.T) {
	ok := IntSlice(mustURL(t, "ids=1,2,50"), "ids", Each(Between(1, 100)))
	if want := []int{1, 2, 50}; !reflect.DeepEqual(ok, want) {
		t.Errorf("in-range slice = %v, want %v", ok, want)
	}

	r := ParseIntSlice(mustURL(t, "ids=1,2,999"), "ids", Each(Between(1, 100)))
	if r.Error == nil {
		t.Error("element 999 outside [1,100] should reject the slice")
	}
	if len(r.Value) != 0 {
		t.Errorf("rejected slice Value = %v, want empty fallback", r.Value)
	}

	fb := IntSlice(mustURL(t, "ids=1,2,999"), "ids",
		Each(Between(1, 100)), Fallback([]int{7}))
	if want := []int{7}; !reflect.DeepEqual(fb, want) {
		t.Errorf("rejected slice with Fallback = %v, want %v", fb, want)
	}
}

// TestSliceCustomSeparator covers Sep, including a multi-byte separator. The
// separators are percent-encoded in the raw query so they reach qp as literal
// characters of the decoded value (a bare ';' would be eaten by net/url).
func TestSliceCustomSeparator(t *testing.T) {
	got := StringSlice(mustURL(t, "tags=a%3Bb%3Bc"), "tags", Sep[string](";"))
	if want := []string{"a", "b", "c"}; !reflect.DeepEqual(got, want) {
		t.Errorf("semicolon Sep = %v, want %v", got, want)
	}

	got = StringSlice(mustURL(t, "tags=a%3A%3Ab%3A%3Ac"), "tags", Sep[string]("::"))
	if want := []string{"a", "b", "c"}; !reflect.DeepEqual(got, want) {
		t.Errorf("multi-byte Sep = %v, want %v", got, want)
	}
}

// TestSliceAbsentVsEmpty mirrors the scalar presence matrix for slices: absent
// pulls nil, empty pulls an empty slice, and the bare form is always non-nil.
func TestSliceAbsentVsEmpty(t *testing.T) {
	absent := New(mustURL(t, "other=1"))
	if absent.PullIntSlice("ids") != nil {
		t.Error("PullIntSlice on absent param must be nil")
	}
	if got := absent.IntSlice("ids"); got == nil || len(got) != 0 {
		t.Errorf("bare IntSlice on absent = %v, want non-nil empty", got)
	}

	empty := New(mustURL(t, "ids="))
	if p := empty.PullIntSlice("ids"); p == nil || len(p) != 0 {
		t.Errorf("PullIntSlice on empty = %v, want non-nil empty", p)
	}
	if !empty.Empty("ids") || !empty.Contains("ids") {
		t.Error("empty slice param must be Contains && Empty")
	}
}
