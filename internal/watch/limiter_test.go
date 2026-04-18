package watch

import (
	"testing"
	"time"
)

func TestNewScanLimiter_InvalidMinInterval(t *testing.T) {
	_, err := NewScanLimiter(0, time.Minute, 10)
	if err == nil {
		t.Fatal("expected error for zero minInterval")
	}
}

func TestNewScanLimiter_InvalidWindow(t *testing.T) {
	_, err := NewScanLimiter(time.Second, 0, 10)
	if err == nil {
		t.Fatal("expected error for zero window")
	}
}

func TestNewScanLimiter_InvalidMaxPerWindow(t *testing.T) {
	_, err := NewScanLimiter(time.Second, time.Minute, 0)
	if err == nil {
		t.Fatal("expected error for zero maxPerWindow")
	}
}

func TestNewScanLimiter_Valid(t *testing.T) {
	l, err := NewScanLimiter(time.Second, time.Minute, 5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if l == nil {
		t.Fatal("expected non-nil limiter")
	}
}

func TestScanLimiter_Allow_FirstCall(t *testing.T) {
	l, _ := NewScanLimiter(time.Second, time.Minute, 5)
	if !l.Allow(time.Now()) {
		t.Fatal("expected first call to be allowed")
	}
}

func TestScanLimiter_Allow_BlockedByMinInterval(t *testing.T) {
	l, _ := NewScanLimiter(5*time.Second, time.Minute, 10)
	now := time.Now()
	l.Allow(now)
	if l.Allow(now.Add(time.Second)) {
		t.Fatal("expected call within minInterval to be blocked")
	}
}

func TestScanLimiter_Allow_PermitsAfterMinInterval(t *testing.T) {
	l, _ := NewScanLimiter(time.Second, time.Minute, 10)
	now := time.Now()
	l.Allow(now)
	if !l.Allow(now.Add(2 * time.Second)) {
		t.Fatal("expected call after minInterval to be allowed")
	}
}

func TestScanLimiter_Allow_BlockedByWindowLimit(t *testing.T) {
	l, _ := NewScanLimiter(time.Millisecond, time.Minute, 3)
	now := time.Now()
	for i := 0; i < 3; i++ {
		l.Allow(now.Add(time.Duration(i) * 10 * time.Millisecond))
	}
	if l.Allow(now.Add(50 * time.Millisecond)) {
		t.Fatal("expected call exceeding window limit to be blocked")
	}
}

func TestScanLimiter_Reset_ClearsState(t *testing.T) {
	l, _ := NewScanLimiter(time.Minute, time.Hour, 2)
	now := time.Now()
	l.Allow(now)
	l.Allow(now.Add(2 * time.Minute))
	l.Reset()
	if !l.Allow(now.Add(3 * time.Minute)) {
		t.Fatal("expected allow after reset")
	}
}
