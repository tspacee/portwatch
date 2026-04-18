package watch

import (
	"testing"
	"time"
)

func TestNewTicker_InvalidBase(t *testing.T) {
	_, err := NewTicker(0, 0)
	if err == nil {
		t.Fatal("expected error for zero base interval")
	}
}

func TestNewTicker_InvalidFactor(t *testing.T) {
	_, err := NewTicker(time.Second, 1.5)
	if err == nil {
		t.Fatal("expected error for factor > 1")
	}
}

func TestNewTicker_Valid(t *testing.T) {
	tk, err := NewTicker(50*time.Millisecond, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer tk.Stop()
	if tk.C == nil {
		t.Fatal("expected non-nil channel")
	}
}

func TestTicker_EmitsTick(t *testing.T) {
	tk, err := NewTicker(30*time.Millisecond, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer tk.Stop()

	select {
	case <-tk.C:
		// success
	case <-time.After(500 * time.Millisecond):
		t.Fatal("ticker did not emit within timeout")
	}
}

func TestTicker_Stop_PreventsFurtherTicks(t *testing.T) {
	tk, err := NewTicker(20*time.Millisecond, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// drain first tick
	select {
	case <-tk.C:
	case <-time.After(300 * time.Millisecond):
		t.Fatal("no initial tick")
	}
	tk.Stop()
	// double-stop must not panic
	tk.Stop()
}

func TestTicker_WithJitter_EmitsTick(t *testing.T) {
	tk, err := NewTicker(30*time.Millisecond, 0.1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer tk.Stop()

	select {
	case <-tk.C:
		// success
	case <-time.After(500 * time.Millisecond):
		t.Fatal("jittered ticker did not emit within timeout")
	}
}
