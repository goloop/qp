package qp

import (
	"net/url"
	"slices"
	"testing"
)

func TestEachClamp(t *testing.T) {
	q := From(url.Values{"n": {"-2,5,20"}})
	r := q.ParseIntSlice("n", Each(Clamp(0, 10)))
	if r.Error != nil || !slices.Equal(r.Value, []int{0, 5, 10}) {
		t.Fatalf("result = %+v", r)
	}
}
