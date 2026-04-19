package watch

import (
	"sync"
	"testing"
)

func TestNewThreshold_InvalidLimit(t *testing.T) {
	_, err := NewThreshold(0, func(port, value int) {})
	if err == nil {
		t.Fatal("expected error for limit < 1")
	}
}

func TestNewThreshold_NilCallback(t *testing.T) {
	_, err := NewThreshold(3, nil)
	if err == nil {
		t.Fatal("expected error for nil callback")
	}
}

func TestNewThreshold_Valid(t *testing.T) {
	th, err := NewThreshold(3, func(port, value int) {})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if th == nil {
		t.Fatal("expected non-nil Threshold")
	}
}

func TestThreshold_Count_Empty(t *testing.T) {
	th, _ := NewThreshold(3, func(port, value int) {})
	if th.Count(80) != 0 {
		t.Fatal("expected 0 for unseen port")
	}
}

func TestThreshold_Record_BelowLimit_NoCallback(t *testing.T) {
	called := false
	th, _ := NewThreshold(3, func(port, value int) { called = true })
	th.Record(80)
	th.Record(80)
	th.Record(80)
	if called {
		t.Fatal("callback should not fire at or below limit")
	}
	if th.Count(80) != 3 {
		t.Fatalf("expected count 3, got %d", th.Count(80))
	}
}

func TestThreshold_Record_ExceedsLimit_FiresCallback(t *testing.T) {
	var mu sync.Mutex
	var firedPort, firedValue int
	th, _ := NewThreshold(2, func(port, value int) {
		mu.Lock()
		defer mu.Unlock()
		firedPort = port
		firedValue = value
	})
	th.Record(443)
	th.Record(443)
	th.Record(443)
	mu.Lock()
	defer mu.Unlock()
	if firedPort != 443 {
		t.Fatalf("expected port 443, got %d", firedPort)
	}
	if firedValue != 3 {
		t.Fatalf("expected value 3, got %d", firedValue)
	}
}

func TestThreshold_Reset_ClearsCount(t *testing.T) {
	th, _ := NewThreshold(5, func(port, value int) {})
	th.Record(22)
	th.Record(22)
	th.Reset(22)
	if th.Count(22) != 0 {
		t.Fatal("expected count 0 after reset")
	}
}
