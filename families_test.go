package qp

import (
	"reflect"
	"testing"
)

// TestAllScalarFamilies exercises every scalar family method — bare, Get, Pull,
// Parse — through both the holder and the top-level shortcut, for all four
// types. Beyond coverage, it asserts that each surface reports the same value,
// so the thin wrappers can never silently diverge from the engine.
func TestAllScalarFamilies(t *testing.T) {
	u := mustURL(t, "i=7&f=2.5&s=hello&b=yes")
	q := New(u)

	// int.
	if q.Int("i") != 7 || Int(u, "i") != 7 {
		t.Error("int bare mismatch")
	}
	if v, ok := q.GetInt("i"); v != 7 || !ok {
		t.Error("int holder Get mismatch")
	}
	if v, ok := GetInt(u, "i"); v != 7 || !ok {
		t.Error("int top Get mismatch")
	}
	if p := q.PullInt("i"); p == nil || *p != 7 {
		t.Error("int holder Pull mismatch")
	}
	if p := PullInt(u, "i"); p == nil || *p != 7 {
		t.Error("int top Pull mismatch")
	}
	if q.ParseInt("i").Value != 7 || ParseInt(u, "i").Value != 7 {
		t.Error("int Parse mismatch")
	}

	// float64.
	if q.Float("f") != 2.5 || Float(u, "f") != 2.5 {
		t.Error("float bare mismatch")
	}
	if v, ok := q.GetFloat("f"); v != 2.5 || !ok {
		t.Error("float holder Get mismatch")
	}
	if v, ok := GetFloat(u, "f"); v != 2.5 || !ok {
		t.Error("float top Get mismatch")
	}
	if p := q.PullFloat("f"); p == nil || *p != 2.5 {
		t.Error("float holder Pull mismatch")
	}
	if p := PullFloat(u, "f"); p == nil || *p != 2.5 {
		t.Error("float top Pull mismatch")
	}
	if q.ParseFloat("f").Value != 2.5 || ParseFloat(u, "f").Value != 2.5 {
		t.Error("float Parse mismatch")
	}

	// string.
	if q.String("s") != "hello" || String(u, "s") != "hello" {
		t.Error("string bare mismatch")
	}
	if v, ok := q.GetString("s"); v != "hello" || !ok {
		t.Error("string holder Get mismatch")
	}
	if v, ok := GetString(u, "s"); v != "hello" || !ok {
		t.Error("string top Get mismatch")
	}
	if p := q.PullString("s"); p == nil || *p != "hello" {
		t.Error("string holder Pull mismatch")
	}
	if p := PullString(u, "s"); p == nil || *p != "hello" {
		t.Error("string top Pull mismatch")
	}
	if q.ParseString("s").Value != "hello" || ParseString(u, "s").Value != "hello" {
		t.Error("string Parse mismatch")
	}

	// bool.
	if !q.Bool("b") || !Bool(u, "b") {
		t.Error("bool bare mismatch")
	}
	if v, ok := q.GetBool("b"); !v || !ok {
		t.Error("bool holder Get mismatch")
	}
	if v, ok := GetBool(u, "b"); !v || !ok {
		t.Error("bool top Get mismatch")
	}
	if p := q.PullBool("b"); p == nil || !*p {
		t.Error("bool holder Pull mismatch")
	}
	if p := PullBool(u, "b"); p == nil || !*p {
		t.Error("bool top Pull mismatch")
	}
	if !q.ParseBool("b").Value || !ParseBool(u, "b").Value {
		t.Error("bool Parse mismatch")
	}
}

// TestAllSliceFamilies does the same for every slice family method across all
// four element types.
func TestAllSliceFamilies(t *testing.T) {
	u := mustURL(t, "i=1,2,3&f=1.5,2.5&s=a,b,c&b=true,false,yes")
	q := New(u)

	if !reflect.DeepEqual(q.IntSlice("i"), []int{1, 2, 3}) ||
		!reflect.DeepEqual(IntSlice(u, "i"), []int{1, 2, 3}) {
		t.Error("int slice bare mismatch")
	}
	if v, ok := q.GetIntSlice("i"); !ok || len(v) != 3 {
		t.Error("int slice holder Get mismatch")
	}
	if v, ok := GetIntSlice(u, "i"); !ok || len(v) != 3 {
		t.Error("int slice top Get mismatch")
	}
	if len(q.PullIntSlice("i")) != 3 || len(PullIntSlice(u, "i")) != 3 {
		t.Error("int slice Pull mismatch")
	}
	if q.ParseIntSlice("i").Error != nil || ParseIntSlice(u, "i").Error != nil {
		t.Error("int slice Parse mismatch")
	}

	if !reflect.DeepEqual(q.FloatSlice("f"), []float64{1.5, 2.5}) ||
		!reflect.DeepEqual(FloatSlice(u, "f"), []float64{1.5, 2.5}) {
		t.Error("float slice bare mismatch")
	}
	if v, ok := q.GetFloatSlice("f"); !ok || len(v) != 2 {
		t.Error("float slice holder Get mismatch")
	}
	if v, ok := GetFloatSlice(u, "f"); !ok || len(v) != 2 {
		t.Error("float slice top Get mismatch")
	}
	if len(q.PullFloatSlice("f")) != 2 || len(PullFloatSlice(u, "f")) != 2 {
		t.Error("float slice Pull mismatch")
	}
	if q.ParseFloatSlice("f").Error != nil || ParseFloatSlice(u, "f").Error != nil {
		t.Error("float slice Parse mismatch")
	}

	if !reflect.DeepEqual(q.StringSlice("s"), []string{"a", "b", "c"}) ||
		!reflect.DeepEqual(StringSlice(u, "s"), []string{"a", "b", "c"}) {
		t.Error("string slice bare mismatch")
	}
	if v, ok := q.GetStringSlice("s"); !ok || len(v) != 3 {
		t.Error("string slice holder Get mismatch")
	}
	if v, ok := GetStringSlice(u, "s"); !ok || len(v) != 3 {
		t.Error("string slice top Get mismatch")
	}
	if len(q.PullStringSlice("s")) != 3 || len(PullStringSlice(u, "s")) != 3 {
		t.Error("string slice Pull mismatch")
	}
	if q.ParseStringSlice("s").Error != nil || ParseStringSlice(u, "s").Error != nil {
		t.Error("string slice Parse mismatch")
	}

	if !reflect.DeepEqual(q.BoolSlice("b"), []bool{true, false, true}) ||
		!reflect.DeepEqual(BoolSlice(u, "b"), []bool{true, false, true}) {
		t.Error("bool slice bare mismatch")
	}
	if v, ok := q.GetBoolSlice("b"); !ok || len(v) != 3 {
		t.Error("bool slice holder Get mismatch")
	}
	if v, ok := GetBoolSlice(u, "b"); !ok || len(v) != 3 {
		t.Error("bool slice top Get mismatch")
	}
	if len(q.PullBoolSlice("b")) != 3 || len(PullBoolSlice(u, "b")) != 3 {
		t.Error("bool slice Pull mismatch")
	}
	if q.ParseBoolSlice("b").Error != nil || ParseBoolSlice(u, "b").Error != nil {
		t.Error("bool slice Parse mismatch")
	}
}

// TestAllPullAbsentReturnsNil checks that every Pull family returns nil for an
// absent parameter — the "this filter was not supplied" signal — across all
// types and both scalar and slice forms.
func TestAllPullAbsentReturnsNil(t *testing.T) {
	q := New(mustURL(t, "other=1"))

	if q.PullInt("x") != nil || PullInt(nil, "x") != nil {
		t.Error("PullInt absent must be nil")
	}
	if q.PullFloat("x") != nil || PullFloat(nil, "x") != nil {
		t.Error("PullFloat absent must be nil")
	}
	if q.PullString("x") != nil || PullString(nil, "x") != nil {
		t.Error("PullString absent must be nil")
	}
	if q.PullBool("x") != nil || PullBool(nil, "x") != nil {
		t.Error("PullBool absent must be nil")
	}
	if q.PullIntSlice("x") != nil || PullIntSlice(nil, "x") != nil {
		t.Error("PullIntSlice absent must be nil")
	}
	if q.PullFloatSlice("x") != nil || PullFloatSlice(nil, "x") != nil {
		t.Error("PullFloatSlice absent must be nil")
	}
	if q.PullStringSlice("x") != nil || PullStringSlice(nil, "x") != nil {
		t.Error("PullStringSlice absent must be nil")
	}
	if q.PullBoolSlice("x") != nil || PullBoolSlice(nil, "x") != nil {
		t.Error("PullBoolSlice absent must be nil")
	}
}

// TestValuesAccessor covers the Values escape hatch.
func TestValuesAccessor(t *testing.T) {
	q := New(mustURL(t, "a=1&b=2"))
	if got := q.Values(); len(got) != 2 || got.Get("a") != "1" {
		t.Errorf("Values() = %v, want the decoded query map", got)
	}
}
