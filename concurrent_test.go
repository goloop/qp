package qp

import (
	"sync"
	"testing"
)

// TestConcurrentReads runs many goroutines against one Query to prove the
// holder is safe for concurrent reads (run with -race). A Query never mutates
// its decoded values after construction.
func TestConcurrentReads(t *testing.T) {
	q := New(mustURL(t, "age=30&ids=1,2,3&name=alice&active=true"))

	var wg sync.WaitGroup
	for i := 0; i < 64; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 1000; j++ {
				if q.Int("age") != 30 {
					t.Error("age mismatch under concurrency")
					return
				}
				if len(q.IntSlice("ids")) != 3 {
					t.Error("ids mismatch under concurrency")
					return
				}
				_ = q.String("name")
				_ = q.Bool("active")
			}
		}()
	}
	wg.Wait()
}
