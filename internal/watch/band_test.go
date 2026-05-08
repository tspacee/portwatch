package watch

import (
	"testing"
)

func TestNewBand_Empty(t *testing.T) {
	b := NewBand()
	if b == nil {
		t.Fatal("expected non-nil Band")
	}
	if got := b.Hits("system"); got != 0 {
		t.Errorf("expected 0 hits, got %d", got)
	}
}

func TestBand_Define_Valid(t *testing.T) {
	b := NewBand()
	if err := b.Define("system", 1, 1023); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestBand_Define_EmptyName(t *testing.T) {
	b := NewBand()
	if err := b.Define("", 1, 1023); err == nil {
		t.Fatal("expected error for empty name")
	}
}

func TestBand_Define_InvalidRange(t *testing.T) {
	b := NewBand()
	if err := b.Define("bad", 1024, 80); err == nil {
		t.Fatal("expected error for inverted range")
	}
}

func TestBand_Define_PortOutOfBounds(t *testing.T) {
	b := NewBand()
	if err := b.Define("bad", 0, 100); err == nil {
		t.Fatal("expected error for port 0")
	}
}

func TestBand_Classify_IncrementsHits(t *testing.T) {
	b := NewBand()
	_ = b.Define("system", 1, 1023)
	b.Classify([]int{22, 80, 443})
	if got := b.Hits("system"); got != 3 {
		t.Errorf("expected 3 hits, got %d", got)
	}
}

func TestBand_Classify_IgnoresOutOfRange(t *testing.T) {
	b := NewBand()
	_ = b.Define("system", 1, 1023)
	b.Classify([]int{8080, 9090})
	if got := b.Hits("system"); got != 0 {
		t.Errorf("expected 0 hits, got %d", got)
	}
}

func TestBand_Reset_ClearsHits(t *testing.T) {
	b := NewBand()
	_ = b.Define("system", 1, 1023)
	b.Classify([]int{80})
	b.Reset()
	if got := b.Hits("system"); got != 0 {
		t.Errorf("expected 0 after reset, got %d", got)
	}
}

func TestBand_Snapshot_ReturnsCopy(t *testing.T) {
	b := NewBand()
	_ = b.Define("registered", 1024, 49151)
	b.Classify([]int{8080, 3000})
	snap := b.Snapshot()
	snap["registered"] = 999
	if got := b.Hits("registered"); got == 999 {
		t.Error("snapshot mutation affected original")
	}
}
