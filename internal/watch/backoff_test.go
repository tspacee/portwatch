package watch

import (
	"testing"
	"time"
)

func TestNewBackoff_InvalidBase(t *testing.T) {
	_, err := NewBackoff(0, time.Second)
	if err == nil {
		t.Fatal("expected error for zero base delay")
	}
}

func TestNewBackoff_MaxLessThanBase(t *testing.T) {
	_, err := NewBackoff(time.Second, time.Millisecond)
	if err == nil {
		t.Fatal("expected error when max < base")
	}
}

func TestNewBackoff_Valid(t *testing.T) {
	b, err := NewBackoff(100*time.Millisecond, 10*time.Second)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if b == nil {
		t.Fatal("expected non-nil Backoff")
	}
}

func TestBackoff_Next_Exponential(t *testing.T) {
	base := 100 * time.Millisecond
	b, _ := NewBackoff(base, time.Hour)

	d0 := b.Next() // attempt 0: base * 2^0 = 100ms
	d1 := b.Next() // attempt 1: base * 2^1 = 200ms
	d2 := b.Next() // attempt 2: base * 2^2 = 400ms

	if d0 != base {
		t.Errorf("expected %v, got %v", base, d0)
	}
	if d1 != 200*time.Millisecond {
		t.Errorf("expected 200ms, got %v", d1)
	}
	if d2 != 400*time.Millisecond {
		t.Errorf("expected 400ms, got %v", d2)
	}
}

func TestBackoff_Next_CapsAtMax(t *testing.T) {
	b, _ := NewBackoff(100*time.Millisecond, 250*time.Millisecond)

	b.Next() // 100ms
	b.Next() // 200ms
	d := b.Next() // would be 400ms but capped at 250ms

	if d != 250*time.Millisecond {
		t.Errorf("expected delay capped at 250ms, got %v", d)
	}
}

func TestBackoff_Reset(t *testing.T) {
	b, _ := NewBackoff(100*time.Millisecond, time.Second)

	b.Next()
	b.Next()
	if b.Attempt() != 2 {
		t.Fatalf("expected attempt 2, got %d", b.Attempt())
	}

	b.Reset()
	if b.Attempt() != 0 {
		t.Errorf("expected attempt 0 after reset, got %d", b.Attempt())
	}

	d := b.Next()
	if d != 100*time.Millisecond {
		t.Errorf("expected base delay after reset, got %v", d)
	}
}

func TestBackoff_Attempt_Increments(t *testing.T) {
	b, _ := NewBackoff(50*time.Millisecond, time.Second)
	for i := 0; i < 5; i++ {
		if b.Attempt() != i {
			t.Errorf("expected attempt %d, got %d", i, b.Attempt())
		}
		b.Next()
	}
}
