package watch

import (
	"errors"
	"testing"
)

func TestNewHealthChecker_InvalidMaxFails(t *testing.T) {
	_, err := NewHealthChecker(0)
	if err == nil {
		t.Fatal("expected error for maxFails=0")
	}
}

func TestNewHealthChecker_Valid(t *testing.T) {
	hc, err := NewHealthChecker(3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !hc.Status().Healthy {
		t.Error("expected healthy on creation")
	}
}

func TestHealthChecker_RecordSuccess_ResetsState(t *testing.T) {
	hc, _ := NewHealthChecker(2)
	hc.RecordFailure(errors.New("oops"))
	hc.RecordSuccess()
	s := hc.Status()
	if !s.Healthy {
		t.Error("expected healthy after success")
	}
	if s.ConsecFails != 0 {
		t.Errorf("expected 0 consec fails, got %d", s.ConsecFails)
	}
	if s.LastError != nil {
		t.Error("expected nil LastError after success")
	}
	if s.LastSuccess.IsZero() {
		t.Error("expected LastSuccess to be set")
	}
}

func TestHealthChecker_RecordFailure_BelowThreshold(t *testing.T) {
	hc, _ := NewHealthChecker(3)
	hc.RecordFailure(errors.New("err"))
	if !hc.Status().Healthy {
		t.Error("should still be healthy below threshold")
	}
}

func TestHealthChecker_RecordFailure_TripsAtThreshold(t *testing.T) {
	hc, _ := NewHealthChecker(2)
	hc.RecordFailure(errors.New("e1"))
	hc.RecordFailure(errors.New("e2"))
	s := hc.Status()
	if s.Healthy {
		t.Error("expected unhealthy after reaching maxFails")
	}
	if s.ConsecFails != 2 {
		t.Errorf("expected ConsecFails=2, got %d", s.ConsecFails)
	}
}

func TestHealthChecker_Status_ReturnsCopy(t *testing.T) {
	hc, _ := NewHealthChecker(1)
	s1 := hc.Status()
	hc.RecordFailure(errors.New("fail"))
	s2 := hc.Status()
	if s1.Healthy == s2.Healthy && s1.ConsecFails == s2.ConsecFails {
		t.Error("expected status to differ after failure")
	}
}
