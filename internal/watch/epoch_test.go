package watch

import (
	"testing"
	"time"
)

func TestNewEpoch_InitialState(t *testing.T) {
	e := NewEpoch()
	if e.Current() != 0 {
		t.Fatalf("expected current=0, got %d", e.Current())
	}
}

func TestEpoch_Advance_IncrementsCounter(t *testing.T) {
	e := NewEpoch()
	n := e.Advance()
	if n != 1 {
		t.Fatalf("expected 1, got %d", n)
	}
	if e.Current() != 1 {
		t.Fatalf("expected current=1, got %d", e.Current())
	}
}

func TestEpoch_Advance_Monotonic(t *testing.T) {
	e := NewEpoch()
	for i := uint64(1); i <= 5; i++ {
		n := e.Advance()
		if n != i {
			t.Fatalf("expected %d, got %d", i, n)
		}
	}
}

func TestEpoch_Since_ValidEpoch(t *testing.T) {
	e := NewEpoch()
	e.Advance()
	time.Sleep(2 * time.Millisecond)
	d, err := e.Since(1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if d < time.Millisecond {
		t.Fatalf("expected duration >= 1ms, got %v", d)
	}
}

func TestEpoch_Since_ZeroEpoch_ReturnsError(t *testing.T) {
	e := NewEpoch()
	e.Advance()
	_, err := e.Since(0)
	if err == nil {
		t.Fatal("expected error for epoch=0")
	}
}

func TestEpoch_Since_OutOfRange_ReturnsError(t *testing.T) {
	e := NewEpoch()
	e.Advance()
	_, err := e.Since(99)
	if err == nil {
		t.Fatal("expected error for out-of-range epoch")
	}
}

func TestEpoch_Reset_ClearsState(t *testing.T) {
	e := NewEpoch()
	e.Advance()
	e.Advance()
	e.Reset()
	if e.Current() != 0 {
		t.Fatalf("expected 0 after reset, got %d", e.Current())
	}
	if e.Age() != 0 {
		t.Fatalf("expected zero age after reset")
	}
}

func TestEpoch_Age_ZeroBeforeAdvance(t *testing.T) {
	e := NewEpoch()
	if e.Age() != 0 {
		t.Fatal("expected zero age before first advance")
	}
}

func TestEpoch_Age_NonZeroAfterAdvance(t *testing.T) {
	e := NewEpoch()
	e.Advance()
	time.Sleep(2 * time.Millisecond)
	if e.Age() < time.Millisecond {
		t.Fatal("expected non-zero age after advance")
	}
}
