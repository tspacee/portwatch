package watch

import (
	"testing"
	"time"
)

func TestNewHorizon_InvalidCutoff(t *testing.T) {
	_, err := NewHorizon(0)
	if err == nil {
		t.Fatal("expected error for zero cutoff")
	}
	_, err = NewHorizon(-time.Second)
	if err == nil {
		t.Fatal("expected error for negative cutoff")
	}
}

func TestNewHorizon_Valid(t *testing.T) {
	h, err := NewHorizon(time.Minute)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if h.Len() != 0 {
		t.Errorf("expected empty horizon, got %d", h.Len())
	}
}

func TestHorizon_Observe_InvalidPort(t *testing.T) {
	h, _ := NewHorizon(time.Minute)
	if err := h.Observe(0); err == nil {
		t.Error("expected error for port 0")
	}
	if err := h.Observe(70000); err == nil {
		t.Error("expected error for port 70000")
	}
}

func TestHorizon_Observe_TracksPort(t *testing.T) {
	h, _ := NewHorizon(time.Minute)
	if err := h.Observe(8080); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if h.Len() != 1 {
		t.Errorf("expected 1 tracked port, got %d", h.Len())
	}
}

func TestHorizon_Observe_Idempotent(t *testing.T) {
	h, _ := NewHorizon(time.Minute)
	h.Observe(443)
	h.Observe(443)
	if h.Len() != 1 {
		t.Errorf("expected 1 tracked port after double observe, got %d", h.Len())
	}
}

func TestHorizon_Age_NeverObserved(t *testing.T) {
	h, _ := NewHorizon(time.Minute)
	if age := h.Age(9090); age != 0 {
		t.Errorf("expected zero age for unseen port, got %v", age)
	}
}

func TestHorizon_Beyond_NotYet(t *testing.T) {
	h, _ := NewHorizon(time.Hour)
	h.Observe(80)
	if h.Beyond(80) {
		t.Error("expected port not yet beyond horizon")
	}
}

func TestHorizon_Forget_RemovesPort(t *testing.T) {
	h, _ := NewHorizon(time.Minute)
	h.Observe(22)
	h.Forget(22)
	if h.Len() != 0 {
		t.Errorf("expected 0 tracked ports after forget, got %d", h.Len())
	}
	if h.Age(22) != 0 {
		t.Error("expected zero age after forget")
	}
}
