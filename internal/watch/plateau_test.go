package watch

import (
	"testing"
	"time"
)

func TestNewPlateau_InvalidDuration(t *testing.T) {
	_, err := NewPlateau(0)
	if err == nil {
		t.Fatal("expected error for zero duration")
	}
	_, err = NewPlateau(-1 * time.Second)
	if err == nil {
		t.Fatal("expected error for negative duration")
	}
}

func TestNewPlateau_Valid(t *testing.T) {
	p, err := NewPlateau(100 * time.Millisecond)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.Len() != 0 {
		t.Errorf("expected empty plateau, got len=%d", p.Len())
	}
}

func TestPlateau_Observe_InvalidPort(t *testing.T) {
	p, _ := NewPlateau(50 * time.Millisecond)
	if err := p.Observe(0, true); err == nil {
		t.Error("expected error for port 0")
	}
	if err := p.Observe(65536, true); err == nil {
		t.Error("expected error for port 65536")
	}
}

func TestPlateau_Stable_NeverObserved(t *testing.T) {
	p, _ := NewPlateau(50 * time.Millisecond)
	if p.Stable(8080) {
		t.Error("expected false for unseen port")
	}
}

func TestPlateau_Stable_NotYet(t *testing.T) {
	p, _ := NewPlateau(200 * time.Millisecond)
	_ = p.Observe(8080, true)
	if p.Stable(8080) {
		t.Error("expected not stable immediately after observe")
	}
}

func TestPlateau_Stable_AfterDuration(t *testing.T) {
	p, _ := NewPlateau(30 * time.Millisecond)
	_ = p.Observe(8080, true)
	time.Sleep(50 * time.Millisecond)
	if !p.Stable(8080) {
		t.Error("expected stable after minStable elapsed")
	}
}

func TestPlateau_Observe_StateChange_ResetsTimer(t *testing.T) {
	p, _ := NewPlateau(30 * time.Millisecond)
	_ = p.Observe(8080, true)
	time.Sleep(50 * time.Millisecond)
	// state changes — timer should reset
	_ = p.Observe(8080, false)
	if p.Stable(8080) {
		t.Error("expected not stable after state change")
	}
}

func TestPlateau_Reset_ClearsState(t *testing.T) {
	p, _ := NewPlateau(30 * time.Millisecond)
	_ = p.Observe(8080, true)
	_ = p.Observe(9090, false)
	if p.Len() != 2 {
		t.Fatalf("expected len 2, got %d", p.Len())
	}
	p.Reset()
	if p.Len() != 0 {
		t.Errorf("expected empty after reset, got len=%d", p.Len())
	}
	if p.Stable(8080) {
		t.Error("expected not stable after reset")
	}
}

func TestPlateau_Observe_SameState_DoesNotResetTimer(t *testing.T) {
	p, _ := NewPlateau(30 * time.Millisecond)
	_ = p.Observe(443, true)
	time.Sleep(40 * time.Millisecond)
	// same state — should not reset the clock
	_ = p.Observe(443, true)
	if !p.Stable(443) {
		t.Error("expected stable: same state re-observe should not reset timer")
	}
}
