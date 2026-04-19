package watch

import (
	"testing"
)

func TestNewTally_Empty(t *testing.T) {
	tal := NewTally()
	if tal == nil {
		t.Fatal("expected non-nil Tally")
	}
	if got := tal.Get(80); got != 0 {
		t.Fatalf("expected 0, got %d", got)
	}
}

func TestTally_Inc_Valid(t *testing.T) {
	tal := NewTally()
	if err := tal.Inc(443); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := tal.Inc(443); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := tal.Get(443); got != 2 {
		t.Fatalf("expected 2, got %d", got)
	}
}

func TestTally_Inc_InvalidPort(t *testing.T) {
	tal := NewTally()
	if err := tal.Inc(0); err == nil {
		t.Fatal("expected error for port 0")
	}
	if err := tal.Inc(70000); err == nil {
		t.Fatal("expected error for port 70000")
	}
}

func TestTally_Reset_ClearsCount(t *testing.T) {
	tal := NewTally()
	_ = tal.Inc(22)
	tal.Reset(22)
	if got := tal.Get(22); got != 0 {
		t.Fatalf("expected 0 after reset, got %d", got)
	}
}

func TestTally_Snapshot_ReturnsCopy(t *testing.T) {
	tal := NewTally()
	_ = tal.Inc(80)
	_ = tal.Inc(80)
	_ = tal.Inc(8080)
	snap := tal.Snapshot()
	if snap[80] != 2 {
		t.Fatalf("expected 2 for port 80, got %d", snap[80])
	}
	if snap[8080] != 1 {
		t.Fatalf("expected 1 for port 8080, got %d", snap[8080])
	}
	// mutating copy should not affect tally
	snap[80] = 999
	if tal.Get(80) != 2 {
		t.Fatal("snapshot mutation affected original tally")
	}
}

func TestTally_Top_ReturnsHighestPort(t *testing.T) {
	tal := NewTally()
	_ = tal.Inc(22)
	_ = tal.Inc(443)
	_ = tal.Inc(443)
	_ = tal.Inc(443)
	_ = tal.Inc(80)
	_ = tal.Inc(80)
	if top := tal.Top(); top != 443 {
		t.Fatalf("expected top port 443, got %d", top)
	}
}

func TestTally_Top_Empty(t *testing.T) {
	tal := NewTally()
	if top := tal.Top(); top != 0 {
		t.Fatalf("expected 0 for empty tally, got %d", top)
	}
}
