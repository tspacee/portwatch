package watch

import (
	"testing"
	"time"
)

func TestNewBurst_InvalidWindow(t *testing.T) {
	_, err := NewBurst(0, 3)
	if err == nil {
		t.Fatal("expected error for zero window")
	}
}

func TestNewBurst_InvalidThreshold(t *testing.T) {
	_, err := NewBurst(time.Second, 0)
	if err == nil {
		t.Fatal("expected error for zero threshold")
	}
}

func TestNewBurst_Valid(t *testing.T) {
	b, err := NewBurst(time.Second, 5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if b == nil {
		t.Fatal("expected non-nil Burst")
	}
}

func TestBurst_Count_Empty(t *testing.T) {
	b, _ := NewBurst(time.Second, 3)
	if b.Count() != 0 {
		t.Errorf("expected 0, got %d", b.Count())
	}
}

func TestBurst_Record_IncrementsCount(t *testing.T) {
	b, _ := NewBurst(time.Second, 5)
	b.Record()
	b.Record()
	if b.Count() != 2 {
		t.Errorf("expected 2, got %d", b.Count())
	}
}

func TestBurst_IsBursting_BelowThreshold(t *testing.T) {
	b, _ := NewBurst(time.Second, 5)
	b.Record()
	b.Record()
	if b.IsBursting() {
		t.Error("expected not bursting below threshold")
	}
}

func TestBurst_IsBursting_AtThreshold(t *testing.T) {
	b, _ := NewBurst(time.Second, 3)
	b.Record()
	b.Record()
	b.Record()
	if !b.IsBursting() {
		t.Error("expected bursting at threshold")
	}
}

func TestBurst_Reset_ClearsEvents(t *testing.T) {
	b, _ := NewBurst(time.Second, 3)
	b.Record()
	b.Record()
	b.Reset()
	if b.Count() != 0 {
		t.Errorf("expected 0 after reset, got %d", b.Count())
	}
}

func TestBurst_EvictsExpiredEvents(t *testing.T) {
	b, _ := NewBurst(50*time.Millisecond, 2)
	b.Record()
	b.Record()
	time.Sleep(60 * time.Millisecond)
	if b.Count() != 0 {
		t.Errorf("expected 0 after expiry, got %d", b.Count())
	}
}
