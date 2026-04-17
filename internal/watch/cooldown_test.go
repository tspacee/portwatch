package watch

import (
	"testing"
	"time"
)

func TestNewCooldown_InvalidWindow(t *testing.T) {
	_, err := NewCooldown(0)
	if err == nil {
		t.Fatal("expected error for zero window")
	}
	if err != ErrInvalidCooldownWindow {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestNewCooldown_Valid(t *testing.T) {
	c, err := NewCooldown(100 * time.Millisecond)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c == nil {
		t.Fatal("expected non-nil cooldown")
	}
}

func TestCooldown_Ready_FirstCall(t *testing.T) {
	c, _ := NewCooldown(100 * time.Millisecond)
	if !c.Ready() {
		t.Fatal("expected Ready to return true on first call")
	}
}

func TestCooldown_Ready_BlockedWithinWindow(t *testing.T) {
	c, _ := NewCooldown(200 * time.Millisecond)
	c.Ready() // fire once
	if c.Ready() {
		t.Fatal("expected Ready to return false within cooldown window")
	}
}

func TestCooldown_Ready_PermitsAfterWindow(t *testing.T) {
	c, _ := NewCooldown(30 * time.Millisecond)
	c.Ready()
	time.Sleep(50 * time.Millisecond)
	if !c.Ready() {
		t.Fatal("expected Ready to return true after window elapsed")
	}
}

func TestCooldown_Reset_MakesReady(t *testing.T) {
	c, _ := NewCooldown(500 * time.Millisecond)
	c.Ready()
	c.Reset()
	if !c.Ready() {
		t.Fatal("expected Ready to return true after Reset")
	}
}

func TestCooldown_Remaining_ZeroBeforeFire(t *testing.T) {
	c, _ := NewCooldown(100 * time.Millisecond)
	if r := c.Remaining(); r != 0 {
		t.Fatalf("expected 0 remaining before first fire, got %v", r)
	}
}

func TestCooldown_Remaining_PositiveAfterFire(t *testing.T) {
	c, _ := NewCooldown(500 * time.Millisecond)
	c.Ready()
	if r := c.Remaining(); r <= 0 {
		t.Fatalf("expected positive remaining after fire, got %v", r)
	}
}
