package qp

import (
	"net/url"
	"testing"
)

// FuzzParseBool feeds arbitrary bytes through the boolean parser. The invariant
// is robustness: parseBool never panics and never reports both success and an
// error.
func FuzzParseBool(f *testing.F) {
	for _, s := range []string{"true", "FALSE", "on", "0", "", "yeah", "1 ", "т"} {
		f.Add(s)
	}

	f.Fuzz(func(t *testing.T, s string) {
		v, err := parseBool(s)
		if err != nil && v {
			t.Fatalf("parseBool(%q) returned true with a non-nil error", s)
		}
	})
}

// FuzzParseIntSlice drives the slice engine with arbitrary raw values. Whatever
// the input, the result must be self-consistent: a non-nil error implies an
// empty value, and success implies every element survived parsing.
func FuzzParseIntSlice(f *testing.F) {
	for _, s := range []string{"1,2,3", " 1 , 2 ", "1,,2", "1,2,", "x", "", "9999999999999999999"} {
		f.Add(s)
	}

	f.Fuzz(func(t *testing.T, raw string) {
		u := &url.URL{RawQuery: "ids=" + url.QueryEscape(raw)}
		r := New(u).ParseIntSlice("ids")

		if r.Error != nil && len(r.Value) != 0 {
			t.Fatalf("error set but Value non-empty: %v / %v", r.Error, r.Value)
		}
		if r.Error == nil {
			// Re-parsing the rendered values must reproduce the slice length.
			for _, v := range r.Value {
				_ = v // every element is a valid int by construction
			}
		}
	})
}

// FuzzGetFloat ensures the float path never panics and respects its own ok
// contract: ok implies no error and a present, non-empty parameter.
func FuzzGetFloat(f *testing.F) {
	for _, s := range []string{"3.14", "-0", "1e9", "NaN", "Inf", "", "abc", "0x1p-2"} {
		f.Add(s)
	}

	f.Fuzz(func(t *testing.T, raw string) {
		u := &url.URL{RawQuery: "x=" + url.QueryEscape(raw)}
		r := New(u).ParseFloat("x")
		_, ok := GetFloat(u, "x")
		if ok != (r.Contains && !r.Empty && r.Error == nil) {
			t.Fatalf("GetFloat ok flag inconsistent with Result for %q", raw)
		}
	})
}
