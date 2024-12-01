package qp

import (
	"net/url"
	"testing"
)

// BenchmarkBooleanParsing benchmarks boolean parsing functions
func BenchmarkBooleanParsing(b *testing.B) {
	urls := map[string]*url.URL{
		"empty":     mustParseURL("http://example.com"),
		"true":      mustParseURL("http://example.com?active=true"),
		"false":     mustParseURL("http://example.com?active=false"),
		"yes":       mustParseURL("http://example.com?active=yes"),
		"no":        mustParseURL("http://example.com?active=no"),
		"on":        mustParseURL("http://example.com?active=on"),
		"off":       mustParseURL("http://example.com?active=off"),
		"1":         mustParseURL("http://example.com?active=1"),
		"0":         mustParseURL("http://example.com?active=0"),
		"invalid":   mustParseURL("http://example.com?active=invalid"),
		"multiple":  mustParseURL("http://example.com?active=true&active=false"),
		"withEmpty": mustParseURL("http://example.com?active="),
	}

	b.Run("ParseBool/empty", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = ParseBool(urls["empty"], "active")
		}
	})

	b.Run("ParseBool/true", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = ParseBool(urls["true"], "active")
		}
	})

	b.Run("ParseBool/withDefault", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = ParseBool(urls["empty"], "active", true)
		}
	})

	b.Run("GetBool/valid", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, _ = GetBool(urls["true"], "active")
		}
	})

	b.Run("PullBool/valid", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = PullBool(urls["true"], "active")
		}
	})

	// Benchmark bool slice parsing
	sliceUrls := map[string]*url.URL{
		"empty":     mustParseURL("http://example.com"),
		"single":    mustParseURL("http://example.com?flags=true"),
		"multiple":  mustParseURL("http://example.com?flags=true,false,yes"),
		"separate":  mustParseURL("http://example.com?flags=true&flags=false&flags=yes"),
		"invalid":   mustParseURL("http://example.com?flags=true,invalid,yes"),
		"withEmpty": mustParseURL("http://example.com?flags="),
	}

	b.Run("ParseBoolSlice/empty", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = ParseBoolSlice(sliceUrls["empty"], "flags")
		}
	})

	b.Run("ParseBoolSlice/single", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = ParseBoolSlice(sliceUrls["single"], "flags")
		}
	})

	b.Run("ParseBoolSlice/multiple", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = ParseBoolSlice(sliceUrls["multiple"], "flags")
		}
	})

	b.Run("ParseBoolSlice/separate", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = ParseBoolSlice(sliceUrls["separate"], "flags")
		}
	})
}

// BenchmarkFloatParsing benchmarks float parsing functions
func BenchmarkFloatParsing(b *testing.B) {
	urls := map[string]*url.URL{
		"empty":     mustParseURL("http://example.com"),
		"zero":      mustParseURL("http://example.com?temp=0"),
		"positive":  mustParseURL("http://example.com?temp=23.5"),
		"negative":  mustParseURL("http://example.com?temp=-10.7"),
		"invalid":   mustParseURL("http://example.com?temp=invalid"),
		"multiple":  mustParseURL("http://example.com?temp=23.5&temp=24.6"),
		"withEmpty": mustParseURL("http://example.com?temp="),
	}

	b.Run("ParseFloat/empty", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = ParseFloat(urls["empty"], "temp")
		}
	})

	b.Run("ParseFloat/valid", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = ParseFloat(urls["positive"], "temp")
		}
	})

	b.Run("ParseFloat/withRange", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = ParseFloat(urls["positive"], "temp", 20.0, 30.0)
		}
	})

	b.Run("GetFloat/valid", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, _ = GetFloat(urls["positive"], "temp")
		}
	})

	b.Run("PullFloat/valid", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = PullFloat(urls["positive"], "temp")
		}
	})

	// Benchmark float slice parsing
	sliceUrls := map[string]*url.URL{
		"empty":     mustParseURL("http://example.com"),
		"single":    mustParseURL("http://example.com?temps=23.5"),
		"multiple":  mustParseURL("http://example.com?temps=23.5,24.6,25.7"),
		"separate":  mustParseURL("http://example.com?temps=23.5&temps=24.6&temps=25.7"),
		"invalid":   mustParseURL("http://example.com?temps=23.5,invalid,25.7"),
		"withEmpty": mustParseURL("http://example.com?temps="),
	}

	b.Run("ParseFloatSlice/empty", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = ParseFloatSlice(sliceUrls["empty"], "temps")
		}
	})

	b.Run("ParseFloatSlice/multiple", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = ParseFloatSlice(sliceUrls["multiple"], "temps")
		}
	})

	b.Run("ParseFloatSlice/separate", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = ParseFloatSlice(sliceUrls["separate"], "temps")
		}
	})
}

// BenchmarkIntParsing benchmarks integer parsing functions
func BenchmarkIntParsing(b *testing.B) {
	urls := map[string]*url.URL{
		"empty":     mustParseURL("http://example.com"),
		"zero":      mustParseURL("http://example.com?age=0"),
		"positive":  mustParseURL("http://example.com?age=25"),
		"negative":  mustParseURL("http://example.com?age=-10"),
		"invalid":   mustParseURL("http://example.com?age=invalid"),
		"multiple":  mustParseURL("http://example.com?age=25&age=30"),
		"withEmpty": mustParseURL("http://example.com?age="),
	}

	b.Run("ParseInt/empty", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = ParseInt(urls["empty"], "age")
		}
	})

	b.Run("ParseInt/valid", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = ParseInt(urls["positive"], "age")
		}
	})

	b.Run("ParseInt/withRange", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = ParseInt(urls["positive"], "age", 18, 65)
		}
	})

	b.Run("GetInt/valid", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, _ = GetInt(urls["positive"], "age")
		}
	})

	b.Run("PullInt/valid", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = PullInt(urls["positive"], "age")
		}
	})

	// Benchmark integer slice parsing
	sliceUrls := map[string]*url.URL{
		"empty":     mustParseURL("http://example.com"),
		"single":    mustParseURL("http://example.com?ages=25"),
		"multiple":  mustParseURL("http://example.com?ages=25,30,35"),
		"separate":  mustParseURL("http://example.com?ages=25&ages=30&ages=35"),
		"invalid":   mustParseURL("http://example.com?ages=25,invalid,35"),
		"withEmpty": mustParseURL("http://example.com?ages="),
	}

	b.Run("ParseIntSlice/empty", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = ParseIntSlice(sliceUrls["empty"], "ages")
		}
	})

	b.Run("ParseIntSlice/multiple", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = ParseIntSlice(sliceUrls["multiple"], "ages")
		}
	})

	b.Run("ParseIntSlice/separate", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = ParseIntSlice(sliceUrls["separate"], "ages")
		}
	})
}

// BenchmarkStringParsing benchmarks string parsing functions
func BenchmarkStringParsing(b *testing.B) {
	urls := map[string]*url.URL{
		"empty":     mustParseURL("http://example.com"),
		"simple":    mustParseURL("http://example.com?name=john"),
		"complex":   mustParseURL("http://example.com?name=john+doe"),
		"multiple":  mustParseURL("http://example.com?name=john&name=doe"),
		"withEmpty": mustParseURL("http://example.com?name="),
	}

	b.Run("ParseString/empty", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = ParseString(urls["empty"], "name")
		}
	})

	b.Run("ParseString/valid", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = ParseString(urls["simple"], "name")
		}
	})

	b.Run("ParseString/withValidValues", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = ParseString(urls["simple"], "name", "john", "jane", "bob")
		}
	})

	b.Run("GetString/valid", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, _ = GetString(urls["simple"], "name")
		}
	})

	b.Run("PullString/valid", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = PullString(urls["simple"], "name")
		}
	})

	// Benchmark string slice parsing
	sliceUrls := map[string]*url.URL{
		"empty":     mustParseURL("http://example.com"),
		"single":    mustParseURL("http://example.com?names=john"),
		"multiple":  mustParseURL("http://example.com?names=john,jane,bob"),
		"separate":  mustParseURL("http://example.com?names=john&names=jane&names=bob"),
		"withEmpty": mustParseURL("http://example.com?names="),
	}

	b.Run("ParseStringSlice/empty", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = ParseStringSlice(sliceUrls["empty"], "names")
		}
	})

	b.Run("ParseStringSlice/multiple", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = ParseStringSlice(sliceUrls["multiple"], "names")
		}
	})

	b.Run("ParseStringSlice/separate", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = ParseStringSlice(sliceUrls["separate"], "names")
		}
	})
}

// BenchmarkUtilityFunctions benchmarks utility functions
func BenchmarkUtilityFunctions(b *testing.B) {
	urls := map[string]*url.URL{
		"empty":     mustParseURL("http://example.com"),
		"withParam": mustParseURL("http://example.com?param=value"),
		"withEmpty": mustParseURL("http://example.com?param="),
	}

	b.Run("Contains/absent", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = Contains(urls["empty"], "param")
		}
	})

	b.Run("Contains/present", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = Contains(urls["withParam"], "param")
		}
	})

	b.Run("Empty/absent", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = Empty(urls["empty"], "param")
		}
	})

	b.Run("Empty/present", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = Empty(urls["withEmpty"], "param")
		}
	})
}

// Helper function to parse URLs
func mustParseURL(s string) *url.URL {
	u, err := url.Parse(s)
	if err != nil {
		panic(err)
	}
	return u
}
