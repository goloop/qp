package qp

import (
	"math"
	"reflect"
	"testing"
)

// TestBoolForms exercises every accepted boolean spelling, case-insensitivity,
// and rejection of near-misses.
func TestBoolForms(t *testing.T) {
	tests := []struct {
		raw   string
		want  bool
		valid bool
	}{
		{"true", true, true},
		{"TRUE", true, true},
		{"True", true, true},
		{"yes", true, true},
		{"on", true, true},
		{"1", true, true},
		{"false", false, true},
		{"NO", false, true},
		{"off", false, true},
		{"0", false, true},
		{"t", false, false},    // strconv would accept this; we do not
		{"2", false, false},    // not a recognised bool
		{"yeah", false, false}, // near-miss
	}

	for _, tc := range tests {
		t.Run(tc.raw, func(t *testing.T) {
			got, ok := GetBool(mustURL(t, "b="+tc.raw), "b")
			if ok != tc.valid {
				t.Errorf("%q: ok=%v, want %v", tc.raw, ok, tc.valid)
			}
			if tc.valid && got != tc.want {
				t.Errorf("%q: got %v, want %v", tc.raw, got, tc.want)
			}
		})
	}
}

// TestScalarDoesNotTrim pins that scalar reads keep surrounding whitespace
// (only slice elements are trimmed). "%20true" is a leading-space value and is
// therefore not a recognised boolean.
func TestScalarDoesNotTrim(t *testing.T) {
	if _, ok := GetBool(mustURL(t, "b=%20true"), "b"); ok {
		t.Error("scalar bool must not trim; \" true\" is not a valid bool")
	}
	if v := String(mustURL(t, "s=%20hi%20"), "s"); v != " hi " {
		t.Errorf("scalar string = %q, want \" hi \" (untrimmed)", v)
	}
}

// TestIntOverflow makes sure a value that overflows int is reported as invalid
// rather than wrapping silently.
func TestIntOverflow(t *testing.T) {
	r := ParseInt(mustURL(t, "n=99999999999999999999999999"), "n")
	if r.Error == nil {
		t.Error("overflowing integer must be reported as an error")
	}
	if r.Value != 0 {
		t.Errorf("Value = %d, want 0 fallback", r.Value)
	}
}

// TestFloatSpecials checks negatives, exponents, and that NaN/Inf round-trip as
// values strconv accepts (documenting the contract, not hiding it).
func TestFloatSpecials(t *testing.T) {
	if v, ok := GetFloat(mustURL(t, "x=-2.5e3"), "x"); !ok || v != -2500 {
		t.Errorf("scientific notation: got %v ok=%v, want -2500", v, ok)
	}
	if v, ok := GetFloat(mustURL(t, "x=NaN"), "x"); !ok || !math.IsNaN(v) {
		t.Errorf("NaN: got %v ok=%v, want NaN/true", v, ok)
	}
}

// TestUnicodeAndEscaping confirms percent-encoding and "+"-as-space are decoded
// by net/url before qp ever sees the value, for keys and values alike.
func TestUnicodeAndEscaping(t *testing.T) {
	// net/url decodes %2C to a real comma before qp ever sees the value, so an
	// encoded comma becomes an ordinary separator. There is therefore no way to
	// smuggle a literal comma into a comma-separated element; this pins that
	// contract so it can never regress silently.
	got := StringSlice(mustURL(t, "n=a%2Cb,c"), "n")
	if want := []string{"a", "b", "c"}; !reflect.DeepEqual(got, want) {
		t.Errorf("encoded comma: got %v, want %v", got, want)
	}

	// "+" decodes to a space; UTF-8 keys and values survive intact.
	if v := String(mustURL(t, "greeting=Pryvit+світ"), "greeting"); v != "Pryvit світ" {
		t.Errorf("plus/utf8 value = %q, want \"Pryvit світ\"", v)
	}
	if v := String(mustURL(t, "ім%27я=Тарас"), "ім'я"); v != "Тарас" {
		t.Errorf("utf8 key value = %q, want \"Тарас\"", v)
	}
}

// TestFourFamiliesAgree checks that the bare/Get/Pull/Parse families return
// consistent values for the same input, so users can switch families freely.
func TestFourFamiliesAgree(t *testing.T) {
	u := mustURL(t, "age=25")

	bare := Int(u, "age")
	get, ok := GetInt(u, "age")
	pull := PullInt(u, "age")
	parse := ParseInt(u, "age")

	if bare != 25 || get != 25 || pull == nil || *pull != 25 || parse.Value != 25 {
		t.Errorf("families disagree: bare=%d get=%d pull=%v parse=%d",
			bare, get, *pull, parse.Value)
	}
	if !ok || !parse.Contains || parse.Error != nil {
		t.Error("valid value should be reported valid by all families")
	}
}

// TestDefaultAppliedConsistently makes sure the default reaches every family on
// an absent parameter.
func TestDefaultAppliedConsistently(t *testing.T) {
	u := mustURL(t, "other=1")

	if got := Int(u, "age", Default(18)); got != 18 {
		t.Errorf("bare default = %d, want 18", got)
	}
	if got, ok := GetInt(u, "age", Default(18)); got != 18 || ok {
		t.Errorf("Get default = (%d, %v), want (18, false)", got, ok)
	}
	if PullInt(u, "age", Default(18)) != nil {
		t.Error("Pull on absent must be nil regardless of default")
	}
	if r := ParseInt(u, "age", Default(18)); r.Value != 18 || r.Contains {
		t.Errorf("Parse default = %d contains=%v, want 18/false", r.Value, r.Contains)
	}
}
