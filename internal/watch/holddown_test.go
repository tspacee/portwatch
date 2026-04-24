package watch

import (
	"testing"
	"time"
)

func TestNewHoldDown_InvalidWindow(t *testing.T) {
	_, err := NewHoldDown(0)
	if err == nil {
		t.Fatal("expected error for zero window")
	}
	_, err = NewHoldDown(-1 * time.Second)
	if err == nil {
		t.Fatal("expected error for negative window")
	}
}

func TestNewHoldDown_Valid(t *testing.T) {
	h, err := NewHoldDown(time.Second)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if h == nil {
		t.Fatal("expected non-nil HoldDown")
	}
}

func TestHoldDown_Suppressed_FirstCall_ReturnsFalse(t *testing.T) {
	h, _ := NewHoldDown(time.Second)
	if h.Suppressed(8080) {
		t.Error("first call should not be suppressed")
	}
}

func TestHoldDown_Suppressed_SecondCall_ReturnsTrue(t *testing.T) {
	h, _ := NewHoldDown(time.Second)
	h.Suppressed(8080)
	if !h.Suppressed(8080) {
		t.Error("second call within window should be suppressed")
	}
}

func TestHoldDown_Suppressed_AfterExpiry_ReturnsFalse(t *testing.T) {
	h, _ := NewHoldDown(20 * time.Millisecond)
	h.Suppressed(9090)
	time.Sleep(40 * time.Millisecond)
	if h.Suppressed(9090) {
		t.Error("call after window expiry should not be suppressed")
	}
}

func TestHoldDown_Release_AllowsNextTrigger(t *testing.T) {
	h, _ := NewHoldDown(time.Hour)
	h.Suppressed(443)
	h.Release(443)
	if h.Suppressed(443) {
		t.Error("after release, port should not be suppressed")
	}
}

func TestHoldDown_Len_CountsActive(t *testing.T) {
	h, _ := NewHoldDown(time.Hour)
	if h.Len() != 0 {
		t.Fatal("expected 0 initially")
	}
	h.Suppressed(80)
	h.Suppressed(443)
	if h.Len() != 2 {
		t.Errorf("expected 2, got %d", h.Len())
	}
	h.Release(80)
	if h.Len() != 1 {
		t.Errorf("expected 1 after release, got %d", h.Len())
	}
}

func TestHoldDown_IndependentPorts(t *testing.T) {
	h, _ := NewHoldDown(time.Hour)
	h.Suppressed(22)
	if h.Suppressed(3306) {
		t.Error("unrelated port should not be suppressed")
	}
}
