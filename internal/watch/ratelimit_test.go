package watch

import (
	"testing"
	"time"
)

func TestNewRateLimiter_InvalidInterval(t *testing.T) {
	_, err := NewRateLimiter(0)
	if err == nil {
		t.Fatal("expected error for zero interval, got nil")
	}

	_, err = NewRateLimiter(-1 * time.Second)
	if err == nil {
		t.Fatal("expected error for negative interval, got nil")
	}
}

func TestNewRateLimiter_Valid(t *testing.T) {
	rl, err := NewRateLimiter(100 * time.Millisecond)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rl == nil {
		t.Fatal("expected non-nil RateLimiter")
	}
}

func TestRateLimiter_Allow_FirstCall(t *testing.T) {
	rl, _ := NewRateLimiter(100 * time.Millisecond)
	if err := rl.Allow(); err != nil {
		t.Fatalf("first Allow() should succeed, got: %v", err)
	}
}

func TestRateLimiter_Allow_BlocksWithinInterval(t *testing.T) {
	rl, _ := NewRateLimiter(500 * time.Millisecond)

	if err := rl.Allow(); err != nil {
		t.Fatalf("first Allow() should succeed: %v", err)
	}

	err := rl.Allow()
	if err == nil {
		t.Fatal("second Allow() within interval should be rate-limited")
	}

	var limited *ErrRateLimited
	if rl, ok := err.(*ErrRateLimited); !ok || rl == nil {
		t.Fatalf("expected *ErrRateLimited, got %T", err)
	}
	_ = limited
}

func TestRateLimiter_Allow_PermitsAfterInterval(t *testing.T) {
	rl, _ := NewRateLimiter(50 * time.Millisecond)

	if err := rl.Allow(); err != nil {
		t.Fatalf("first Allow() should succeed: %v", err)
	}

	time.Sleep(60 * time.Millisecond)

	if err := rl.Allow(); err != nil {
		t.Fatalf("Allow() after interval should succeed: %v", err)
	}
}

func TestRateLimiter_Skipped_Increments(t *testing.T) {
	rl, _ := NewRateLimiter(500 * time.Millisecond)
	rl.Allow()
	rl.Allow()
	rl.Allow()

	if got := rl.Skipped(); got != 2 {
		t.Fatalf("expected 2 skipped, got %d", got)
	}
}

func TestRateLimiter_Reset_AllowsImmediateScan(t *testing.T) {
	rl, _ := NewRateLimiter(500 * time.Millisecond)
	rl.Allow()

	rl.Reset()

	if err := rl.Allow(); err != nil {
		t.Fatalf("Allow() after Reset() should succeed: %v", err)
	}
}
