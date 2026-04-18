package watch

import (
	"testing"
	"time"
)

func TestNewJitter_InvalidBase(t *testing.T) {
	_, err := NewJitter(0, 0.1)
	if err == nil {
		t.Fatal("expected error for zero base duration")
	}
}

func TestNewJitter_NegativeFactor(t *testing.T) {
	_, err := NewJitter(time.Second, -0.1)
	if err == nil {
		t.Fatal("expected error for negative factor")
	}
}

func TestNewJitter_FactorExceedsOne(t *testing.T) {
	_, err := NewJitter(time.Second, 1.5)
	if err == nil {
		t.Fatal("expected error for factor > 1")
	}
}

func TestNewJitter_Valid(t *testing.T) {
	j, err := NewJitter(10*time.Second, 0.2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if j.Base() != 10*time.Second {
		t.Errorf("expected base 10s, got %v", j.Base())
	}
	if j.Factor() != 0.2 {
		t.Errorf("expected factor 0.2, got %v", j.Factor())
	}
}

func TestJitter_Next_ZeroFactor_ReturnsBase(t *testing.T) {
	j, _ := NewJitter(5*time.Second, 0)
	for i := 0; i < 10; i++ {
		if got := j.Next(); got != 5*time.Second {
			t.Errorf("expected 5s, got %v", got)
		}
	}
}

func TestJitter_Next_WithFactor_InRange(t *testing.T) {
	base := 10 * time.Second
	factor := 0.3
	j, _ := NewJitter(base, factor)
	max := base + time.Duration(float64(base)*factor)
	for i := 0; i < 50; i++ {
		got := j.Next()
		if got < base || got > max {
			t.Errorf("Next() = %v out of range [%v, %v]", got, base, max)
		}
	}
}

func TestJitter_Next_ProducesVariance(t *testing.T) {
	j, _ := NewJitter(10*time.Second, 0.5)
	seen := make(map[time.Duration]bool)
	for i := 0; i < 20; i++ {
		seen[j.Next()] = true
	}
	if len(seen) < 2 {
		t.Error("expected variance in jitter output, got none")
	}
}
