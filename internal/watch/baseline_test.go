package watch

import (
	"testing"
)

func TestNewBaseline_Empty(t *testing.T) {
	b := NewBaseline()
	if b.Len() != 0 {
		t.Fatalf("expected 0 ports, got %d", b.Len())
	}
	if b.IsFrozen() {
		t.Fatal("expected baseline to be unfrozen")
	}
}

func TestBaseline_Learn_Valid(t *testing.T) {
	b := NewBaseline()
	if err := b.Learn(8080); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if b.Len() != 1 {
		t.Fatalf("expected 1 port, got %d", b.Len())
	}
}

func TestBaseline_Learn_InvalidPort(t *testing.T) {
	b := NewBaseline()
	if err := b.Learn(0); err == nil {
		t.Fatal("expected error for port 0")
	}
	if err := b.Learn(65536); err == nil {
		t.Fatal("expected error for port 65536")
	}
}

func TestBaseline_Learn_AfterFreeze_ReturnsError(t *testing.T) {
	b := NewBaseline()
	b.Freeze()
	if err := b.Learn(9090); err == nil {
		t.Fatal("expected error when learning after freeze")
	}
}

func TestBaseline_Freeze_SetsFlag(t *testing.T) {
	b := NewBaseline()
	b.Freeze()
	if !b.IsFrozen() {
		t.Fatal("expected baseline to be frozen")
	}
}

func TestBaseline_Contains_True(t *testing.T) {
	b := NewBaseline()
	_ = b.Learn(443)
	if !b.Contains(443) {
		t.Fatal("expected baseline to contain port 443")
	}
}

func TestBaseline_Contains_False(t *testing.T) {
	b := NewBaseline()
	if b.Contains(443) {
		t.Fatal("expected baseline not to contain port 443")
	}
}

func TestBaseline_Snapshot_ReturnsCopy(t *testing.T) {
	b := NewBaseline()
	_ = b.Learn(80)
	_ = b.Learn(443)
	snap := b.Snapshot()
	if len(snap) != 2 {
		t.Fatalf("expected 2 ports in snapshot, got %d", len(snap))
	}
	// mutating the snapshot must not affect the baseline
	snap[0] = 9999
	if b.Len() != 2 {
		t.Fatal("mutating snapshot affected baseline")
	}
}

func TestBaseline_Learn_Idempotent(t *testing.T) {
	b := NewBaseline()
	_ = b.Learn(8080)
	_ = b.Learn(8080)
	if b.Len() != 1 {
		t.Fatalf("expected 1 unique port, got %d", b.Len())
	}
}
