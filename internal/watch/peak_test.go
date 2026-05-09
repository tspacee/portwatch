package watch

import (
	"testing"
)

func TestNewPeak_InitialState(t *testing.T) {
	p := NewPeak()
	v, ok := p.Value()
	if ok {
		t.Fatalf("expected no value set, got %d", v)
	}
}

func TestPeak_Observe_FirstValue(t *testing.T) {
	p := NewPeak()
	if err := p.Observe(5); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	v, ok := p.Value()
	if !ok {
		t.Fatal("expected value to be set")
	}
	if v != 5 {
		t.Fatalf("expected 5, got %d", v)
	}
}

func TestPeak_Observe_UpdatesOnHigher(t *testing.T) {
	p := NewPeak()
	_ = p.Observe(3)
	_ = p.Observe(7)
	v, _ := p.Value()
	if v != 7 {
		t.Fatalf("expected peak 7, got %d", v)
	}
}

func TestPeak_Observe_IgnoresLower(t *testing.T) {
	p := NewPeak()
	_ = p.Observe(10)
	_ = p.Observe(4)
	v, _ := p.Value()
	if v != 10 {
		t.Fatalf("expected peak 10, got %d", v)
	}
}

func TestPeak_Observe_NegativeCount_ReturnsError(t *testing.T) {
	p := NewPeak()
	err := p.Observe(-1)
	if err == nil {
		t.Fatal("expected error for negative count")
	}
}

func TestPeak_Reset_ClearsPeak(t *testing.T) {
	p := NewPeak()
	_ = p.Observe(8)
	p.Reset()
	_, ok := p.Value()
	if ok {
		t.Fatal("expected peak to be cleared after reset")
	}
}

func TestPeak_Observe_AfterReset(t *testing.T) {
	p := NewPeak()
	_ = p.Observe(15)
	p.Reset()
	_ = p.Observe(6)
	v, ok := p.Value()
	if !ok {
		t.Fatal("expected value after reset and re-observe")
	}
	if v != 6 {
		t.Fatalf("expected 6, got %d", v)
	}
}
