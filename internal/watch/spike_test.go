package watch

import (
	"testing"
	"time"
)

func TestNewSpike_InvalidWindow(t *testing.T) {
	_, err := NewSpike(0, 3)
	if err == nil {
		t.Fatal("expected error for zero window")
	}
}

func TestNewSpike_InvalidThreshold(t *testing.T) {
	_, err := NewSpike(time.Second, 0)
	if err == nil {
		t.Fatal("expected error for zero threshold")
	}
}

func TestNewSpike_Valid(t *testing.T) {
	s, err := NewSpike(time.Second, 5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s == nil {
		t.Fatal("expected non-nil Spike")
	}
}

func TestSpike_Count_Empty(t *testing.T) {
	s, _ := NewSpike(time.Second, 3)
	if got := s.Count(); got != 0 {
		t.Fatalf("expected 0, got %d", got)
	}
}

func TestSpike_Record_IncrementsCount(t *testing.T) {
	s, _ := NewSpike(time.Second, 3)
	s.Record()
	s.Record()
	if got := s.Count(); got != 2 {
		t.Fatalf("expected 2, got %d", got)
	}
}

func TestSpike_Triggered_BelowThreshold(t *testing.T) {
	s, _ := NewSpike(time.Second, 5)
	s.Record()
	s.Record()
	if s.Triggered() {
		t.Fatal("expected not triggered below threshold")
	}
}

func TestSpike_Triggered_AtThreshold(t *testing.T) {
	s, _ := NewSpike(time.Second, 3)
	s.Record()
	s.Record()
	s.Record()
	if !s.Triggered() {
		t.Fatal("expected triggered at threshold")
	}
}

func TestSpike_Reset_ClearsCount(t *testing.T) {
	s, _ := NewSpike(time.Second, 3)
	s.Record()
	s.Record()
	s.Reset()
	if got := s.Count(); got != 0 {
		t.Fatalf("expected 0 after reset, got %d", got)
	}
}

func TestSpike_Prune_RemovesExpiredEvents(t *testing.T) {
	s, _ := NewSpike(50*time.Millisecond, 2)
	s.Record()
	s.Record()
	time.Sleep(80 * time.Millisecond)
	if got := s.Count(); got != 0 {
		t.Fatalf("expected 0 after window expiry, got %d", got)
	}
}
