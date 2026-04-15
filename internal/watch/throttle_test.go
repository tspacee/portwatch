package watch

import (
	"testing"
	"time"
)

func TestNewThrottle_InvalidWindow(t *testing.T) {
	_, err := NewThrottle(0, 5)
	if err != ErrInvalidWindow {
		t.Fatalf("expected ErrInvalidWindow, got %v", err)
	}
}

func TestNewThrottle_InvalidMaxCount(t *testing.T) {
	_, err := NewThrottle(time.Second, 0)
	if err != ErrInvalidMaxCount {
		t.Fatalf("expected ErrInvalidMaxCount, got %v", err)
	}
}

func TestNewThrottle_Valid(t *testing.T) {
	th, err := NewThrottle(time.Second, 3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if th == nil {
		t.Fatal("expected non-nil Throttle")
	}
}

func TestThrottle_Allow_WithinLimit(t *testing.T) {
	th, _ := NewThrottle(time.Second, 3)
	for i := 0; i < 3; i++ {
		if !th.Allow() {
			t.Fatalf("expected Allow()=true on call %d", i+1)
		}
	}
}

func TestThrottle_Allow_ExceedsLimit(t *testing.T) {
	th, _ := NewThrottle(time.Second, 2)
	th.Allow()
	th.Allow()
	if th.Allow() {
		t.Fatal("expected Allow()=false after limit reached")
	}
}

func TestThrottle_Reset_ClearsEvents(t *testing.T) {
	th, _ := NewThrottle(time.Second, 2)
	th.Allow()
	th.Allow()
	th.Reset()
	if !th.Allow() {
		t.Fatal("expected Allow()=true after Reset")
	}
}

func TestThrottle_Count_ReflectsWindow(t *testing.T) {
	th, _ := NewThrottle(200*time.Millisecond, 10)
	th.Allow()
	th.Allow()
	if got := th.Count(); got != 2 {
		t.Fatalf("expected count 2, got %d", got)
	}
	time.Sleep(250 * time.Millisecond)
	if got := th.Count(); got != 0 {
		t.Fatalf("expected count 0 after window expiry, got %d", got)
	}
}

func TestThrottle_Allow_AfterWindowExpiry(t *testing.T) {
	th, _ := NewThrottle(100*time.Millisecond, 1)
	if !th.Allow() {
		t.Fatal("first Allow should succeed")
	}
	if th.Allow() {
		t.Fatal("second Allow should be blocked")
	}
	time.Sleep(150 * time.Millisecond)
	if !th.Allow() {
		t.Fatal("Allow should succeed after window expires")
	}
}
