package watch

import (
	"testing"
)

func TestNewCounter_Empty(t *testing.T) {
	c := NewCounter()
	if c == nil {
		t.Fatal("expected non-nil counter")
	}
	if got := c.Get(80); got != 0 {
		t.Fatalf("expected 0, got %d", got)
	}
}

func TestCounter_Increment_Valid(t *testing.T) {
	c := NewCounter()
	if err := c.Increment(443); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := c.Increment(443); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := c.Get(443); got != 2 {
		t.Fatalf("expected 2, got %d", got)
	}
}

func TestCounter_Increment_InvalidPort(t *testing.T) {
	c := NewCounter()
	if err := c.Increment(0); err == nil {
		t.Fatal("expected error for port 0")
	}
	if err := c.Increment(65536); err == nil {
		t.Fatal("expected error for port 65536")
	}
}

func TestCounter_Reset_ClearsCount(t *testing.T) {
	c := NewCounter()
	_ = c.Increment(8080)
	c.Reset(8080)
	if got := c.Get(8080); got != 0 {
		t.Fatalf("expected 0 after reset, got %d", got)
	}
}

func TestCounter_Snapshot_ReturnsCopy(t *testing.T) {
	c := NewCounter()
	_ = c.Increment(22)
	_ = c.Increment(22)
	_ = c.Increment(80)
	snap := c.Snapshot()
	if snap[22] != 2 {
		t.Fatalf("expected 2 for port 22, got %d", snap[22])
	}
	if snap[80] != 1 {
		t.Fatalf("expected 1 for port 80, got %d", snap[80])
	}
	// mutating snapshot should not affect counter
	snap[22] = 99
	if got := c.Get(22); got != 2 {
		t.Fatalf("counter should not be affected by snapshot mutation, got %d", got)
	}
}

func TestCounter_MultiplePorts_Independent(t *testing.T) {
	c := NewCounter()
	_ = c.Increment(53)
	_ = c.Increment(443)
	_ = c.Increment(443)
	if c.Get(53) != 1 {
		t.Fatalf("expected 1 for port 53")
	}
	if c.Get(443) != 2 {
		t.Fatalf("expected 2 for port 443")
	}
}
