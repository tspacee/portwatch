package watch

import (
	"testing"
	"time"
)

func TestNewDrift_InvalidDecay_Zero(t *testing.T) {
	_, err := NewDrift(0)
	if err == nil {
		t.Fatal("expected error for decay=0")
	}
}

func TestNewDrift_InvalidDecay_ExceedsOne(t *testing.T) {
	_, err := NewDrift(1.1)
	if err == nil {
		t.Fatal("expected error for decay>1")
	}
}

func TestNewDrift_Valid(t *testing.T) {
	d, err := NewDrift(0.5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if d == nil {
		t.Fatal("expected non-nil Drift")
	}
}

func TestDrift_Observe_InvalidPort(t *testing.T) {
	d, _ := NewDrift(0.5)
	err := d.Observe(0, time.Now())
	if err == nil {
		t.Fatal("expected error for port=0")
	}
}

func TestDrift_Score_NeverObserved_IsZero(t *testing.T) {
	d, _ := NewDrift(0.5)
	if score := d.Score(8080); score != 0 {
		t.Fatalf("expected 0, got %v", score)
	}
}

func TestDrift_Observe_FirstCall_NoScore(t *testing.T) {
	d, _ := NewDrift(0.5)
	now := time.Now()
	_ = d.Observe(8080, now)
	if score := d.Score(8080); score != 0 {
		t.Fatalf("expected 0 after first observation, got %v", score)
	}
}

func TestDrift_Observe_SecondCall_SetsBaseline(t *testing.T) {
	d, _ := NewDrift(0.5)
	now := time.Now()
	_ = d.Observe(8080, now)
	_ = d.Observe(8080, now.Add(time.Second))
	// After two observations the baseline is set; score should still be 0
	// because there is no deviation yet.
	if score := d.Score(8080); score != 0 {
		t.Fatalf("expected 0 on stable baseline, got %v", score)
	}
}

func TestDrift_Observe_ThirdCall_DetectsDrift(t *testing.T) {
	d, _ := NewDrift(1.0) // decay=1 means only latest delta matters
	now := time.Now()
	_ = d.Observe(9000, now)
	_ = d.Observe(9000, now.Add(time.Second))    // baseline = 1s
	_ = d.Observe(9000, now.Add(4*time.Second))  // interval = 3s, delta = 2s
	if score := d.Score(9000); score <= 0 {
		t.Fatalf("expected positive drift score, got %v", score)
	}
}

func TestDrift_Reset_ClearsState(t *testing.T) {
	d, _ := NewDrift(1.0)
	now := time.Now()
	_ = d.Observe(443, now)
	_ = d.Observe(443, now.Add(time.Second))
	_ = d.Observe(443, now.Add(5*time.Second))
	d.Reset(443)
	if score := d.Score(443); score != 0 {
		t.Fatalf("expected 0 after reset, got %v", score)
	}
}
