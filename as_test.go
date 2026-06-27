package qp

import (
	"net/netip"
	"testing"
	"time"
)

// TestParseAs covers the generic custom-parser hook with a stdlib type that the
// package does not handle natively (netip.Addr), including option validation.
func TestParseAs(t *testing.T) {
	q := New(mustURL(t, "ip=192.168.0.1&bad=notanip"))

	addr, ok := GetAs(q, "ip", netip.ParseAddr)
	if !ok || addr.String() != "192.168.0.1" {
		t.Errorf("GetAs ip = (%v, %v), want 192.168.0.1/true", addr, ok)
	}

	if _, ok := GetAs(q, "bad", netip.ParseAddr); ok {
		t.Error("GetAs on an unparseable address should report ok=false")
	}

	if p := PullAs(q, "missing", netip.ParseAddr); p != nil {
		t.Error("PullAs on an absent key should be nil")
	}

	// OneOf still works through the custom parser.
	loopback := netip.MustParseAddr("127.0.0.1")
	r := ParseAs(New(mustURL(t, "ip=192.168.0.1")), "ip",
		netip.ParseAddr, OneOf(loopback))
	if r.Error == nil {
		t.Error("value outside the allowed set should be rejected")
	}
}

// TestParseSliceAs covers the slice custom-parser hook in all three families.
func TestParseSliceAs(t *testing.T) {
	q := New(mustURL(t, "ips=10.0.0.1, 10.0.0.2 ,10.0.0.3"))

	addrs, ok := GetSliceAs(q, "ips", netip.ParseAddr)
	if !ok || len(addrs) != 3 || addrs[1].String() != "10.0.0.2" {
		t.Errorf("GetSliceAs = (%v, %v), want 3 trimmed addrs", addrs, ok)
	}

	if r := ParseSliceAs(q, "ips", netip.ParseAddr); r.Error != nil || len(r.Value) != 3 {
		t.Errorf("ParseSliceAs = %v err=%v, want 3 addrs", r.Value, r.Error)
	}

	if got := PullSliceAs(q, "ips", netip.ParseAddr); len(got) != 3 {
		t.Errorf("PullSliceAs = %v, want 3 addrs", got)
	}
	if got := PullSliceAs(q, "missing", netip.ParseAddr); got != nil {
		t.Errorf("PullSliceAs absent = %v, want nil", got)
	}
}

// TestPullAsValue covers the non-nil path of PullAs.
func TestPullAsValue(t *testing.T) {
	q := New(mustURL(t, "ip=192.168.0.1"))
	if p := PullAs(q, "ip", netip.ParseAddr); p == nil || p.String() != "192.168.0.1" {
		t.Errorf("PullAs = %v, want pointer to 192.168.0.1", p)
	}
}

// TestDuration covers the Duration helper end to end, including range validation
// (Duration is ordered, so Min/Max/Between apply).
func TestDuration(t *testing.T) {
	if got := Duration(mustURL(t, "ttl=1.5h"), "ttl"); got != 90*time.Minute {
		t.Errorf("Duration = %v, want 1h30m", got)
	}

	if got := Duration(mustURL(t, "other=1"), "ttl",
		Default(30*time.Second)); got != 30*time.Second {
		t.Errorf("Duration default = %v, want 30s", got)
	}

	// Over the Max -> invalid -> default returned, ok=false.
	v, ok := GetDuration(mustURL(t, "ttl=2h"), "ttl",
		Default(time.Minute), Max(time.Hour))
	if ok || v != time.Minute {
		t.Errorf("GetDuration over max = (%v, %v), want (1m, false)", v, ok)
	}

	if _, ok := GetDuration(mustURL(t, "ttl=nope"), "ttl"); ok {
		t.Error("unparseable duration should report ok=false")
	}

	if p := PullDuration(mustURL(t, "other=1"), "ttl"); p != nil {
		t.Error("PullDuration on absent key should be nil")
	}
	if p := PullDuration(mustURL(t, "ttl=5s"), "ttl"); p == nil || *p != 5*time.Second {
		t.Errorf("PullDuration present = %v, want pointer to 5s", p)
	}
	if r := ParseDuration(mustURL(t, "ttl=5s"), "ttl"); r.Error != nil || r.Value != 5*time.Second {
		t.Errorf("ParseDuration = %v err=%v, want 5s", r.Value, r.Error)
	}
}

// TestTime covers the Time helper with an explicit layout.
func TestTime(t *testing.T) {
	got := Time(mustURL(t, "from=2026-06-27T10:00:00Z"), "from", time.RFC3339)
	want := time.Date(2026, 6, 27, 10, 0, 0, 0, time.UTC)
	if !got.Equal(want) {
		t.Errorf("Time = %v, want %v", got, want)
	}

	if _, ok := GetTime(mustURL(t, "from=27/06/2026"), "from", time.RFC3339); ok {
		t.Error("a value not matching the layout should report ok=false")
	}

	if p := PullTime(mustURL(t, "other=1"), "from", time.RFC3339); p != nil {
		t.Error("PullTime on absent key should be nil")
	}
	if p := PullTime(mustURL(t, "from=2026-06-27T10:00:00Z"), "from",
		time.RFC3339); p == nil || !p.Equal(want) {
		t.Errorf("PullTime present = %v, want pointer to %v", p, want)
	}
	if r := ParseTime(mustURL(t, "from=2026-06-27T10:00:00Z"), "from",
		time.RFC3339); r.Error != nil || !r.Value.Equal(want) {
		t.Errorf("ParseTime = %v err=%v, want %v", r.Value, r.Error, want)
	}
}
