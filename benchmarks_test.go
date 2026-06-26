package qp

import (
	"net/url"
	"testing"
)

var benchURL, _ = url.Parse(
	"http://x/?page=2&limit=50&sort=name&active=true&ids=1,2,3,4,5",
)

// BenchmarkScalarTopLevel measures a single top-level read, which decodes the
// query once per call.
func BenchmarkScalarTopLevel(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = Int(benchURL, "page", Default(1), Min(1))
	}
}

// BenchmarkScalarHolder measures the same read through a holder, isolating the
// per-read cost once the query is decoded.
func BenchmarkScalarHolder(b *testing.B) {
	q := New(benchURL)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = q.Int("page", Default(1), Min(1))
	}
}

// BenchmarkFiveKeysTopLevel reads five parameters the per-call way: five decodes
// of the same query string.
func BenchmarkFiveKeysTopLevel(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = Int(benchURL, "page", Default(1), Min(1))
		_ = Int(benchURL, "limit", Default(20), Between(1, 100))
		_ = String(benchURL, "sort", Default("created_at"))
		_ = Bool(benchURL, "active")
		_ = IntSlice(benchURL, "ids")
	}
}

// BenchmarkFiveKeysHolder reads the same five parameters from one holder: a
// single decode amortised across all reads. This is the OPT-01 win.
func BenchmarkFiveKeysHolder(b *testing.B) {
	for i := 0; i < b.N; i++ {
		q := New(benchURL)
		_ = q.Int("page", Default(1), Min(1))
		_ = q.Int("limit", Default(20), Between(1, 100))
		_ = q.String("sort", Default("created_at"))
		_ = q.Bool("active")
		_ = q.IntSlice("ids")
	}
}

// BenchmarkIntSlice measures slice parsing with element validation.
func BenchmarkIntSlice(b *testing.B) {
	q := New(benchURL)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = q.IntSlice("ids", Each(Between(1, 100)))
	}
}
