package qp

import (
	"net/url"
	"strings"
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

// FuzzParseIntSlice drives the slice engine with two raw values of the same key
// (the repeated form, which is where BUG-01 lived). Whatever the input, the
// result must be self-consistent: a non-nil error implies an empty value, and a
// successful parse must drop empty elements rather than lose surrounding data.
func FuzzParseIntSlice(f *testing.F) {
	seeds := [][2]string{
		{"1,2,3", ""}, {"", "5"}, {"", ""}, {" 1 , 2 ", "3"},
		{"1,,2", "4"}, {"x", "1"}, {"9999999999999999999", "1"},
	}
	for _, s := range seeds {
		f.Add(s[0], s[1])
	}

	f.Fuzz(func(t *testing.T, a, b string) {
		u := &url.URL{RawQuery: "ids=" + url.QueryEscape(a) + "&ids=" + url.QueryEscape(b)}
		r := New(u).ParseIntSlice("ids")

		if r.Error != nil && len(r.Value) != 0 {
			t.Fatalf("error set but Value non-empty: %v / %v", r.Error, r.Value)
		}
		// A leading/trailing empty value must never wipe a parseable one: if
		// exactly one of the two values is a single valid integer, it must
		// appear in the result.
		if r.Error == nil {
			if n, ok := loneInt(a, b); ok && (len(r.Value) != 1 || r.Value[0] != n) {
				t.Fatalf("lone value %d lost: a=%q b=%q -> %v", n, a, b, r.Value)
			}
		}
	})
}

// loneInt returns the single integer value when exactly one of a/b is a bare
// integer and the other is empty/blank, matching the BUG-01 reproduction shape.
func loneInt(a, b string) (int, bool) {
	pa, ea := parseInt(a)
	pb, eb := parseInt(b)
	switch {
	case ea == nil && strings.TrimSpace(b) == "":
		return pa, true
	case eb == nil && strings.TrimSpace(a) == "":
		return pb, true
	}

	return 0, false
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
