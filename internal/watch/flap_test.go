package watch

import (
	"testing"
	"time"
)

func TestNewFlap_InvalidWindow(t *testing.T) {
	_, err := NewFlap(0, 3)
	if err == nil {
		t.Fatal("expected error for zero window")
	}
}

func TestNewFlap_ThresholdTooLow(t *testing.T) {
	_, err := NewFlap(time.Minute, 1)
	if err == nil {
		t.Fatal("expected error for threshold < 2")
	}
}

func TestNewFlap_Valid(t *testing.T) {
	f, err := NewFlap(time.Minute, 3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f == nil {
		t.Fatal("expected non-nil Flap")
	}
}

func TestFlap_Record_InvalidPort(t *testing.T) {
	f, _ := NewFlap(time.Minute, 2)
	err := f.Record(0, time.Now())
	if err == nil {
		t.Fatal("expected error for invalid port")
	}
}

func TestFlap_IsFlapping_BelowThreshold(t *testing.T) {
	f, _ := NewFlap(time.Minute, 3)
	now := time.Now()
	_ = f.Record(8080, now)
	_ = f.Record(8080, now.Add(time.Second))
	if f.IsFlapping(8080, now.Add(2*time.Second)) {
		t.Fatal("expected not flapping below threshold")
	}
}

func TestFlap_IsFlapping_AtThreshold(t *testing.T) {
	f, _ := NewFlap(time.Minute, 3)
	now := time.Now()
	_ = f.Record(8080, now)
	_ = f.Record(8080, now.Add(time.Second))
	_ = f.Record(8080, now.Add(2*time.Second))
	if !f.IsFlapping(8080, now.Add(3*time.Second)) {
		t.Fatal("expected flapping at threshold")
	}
}

func TestFlap_IsFlapping_EvictsExpired(t *testing.T) {
	f, _ := NewFlap(5*time.Second, 3)
	now := time.Now()
	_ = f.Record(443, now.Add(-10*time.Second))
	_ = f.Record(443, now.Add(-9*time.Second))
	_ = f.Record(443, now.Add(-8*time.Second))
	if f.IsFlapping(443, now) {
		t.Fatal("expected expired events to be evicted")
	}
}

func TestFlap_Reset_ClearsEvents(t *testing.T) {
	f, _ := NewFlap(time.Minute, 2)
	now := time.Now()
	_ = f.Record(22, now)
	_ = f.Record(22, now.Add(time.Second))
	f.Reset(22)
	if f.IsFlapping(22, now.Add(2*time.Second)) {
		t.Fatal("expected flap cleared after reset")
	}
}
