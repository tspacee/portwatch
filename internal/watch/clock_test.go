package watch

import (
	"testing"
	"time"
)

func TestNewClock_InitialState(t *testing.T) {
	c := NewClock()
	if c.Count() != 0 {
		t.Fatalf("expected count 0, got %d", c.Count())
	}
	if !c.LastScan().IsZero() {
		t.Fatal("expected zero LastScan")
	}
	if c.Elapsed() != 0 {
		t.Fatal("expected zero Elapsed")
	}
}

func TestClock_StartStop_RecordsElapsed(t *testing.T) {
	c := NewClock()
	if err := c.Start(); err != nil {
		t.Fatalf("unexpected error on Start: %v", err)
	}
	time.Sleep(5 * time.Millisecond)
	if err := c.Stop(); err != nil {
		t.Fatalf("unexpected error on Stop: %v", err)
	}
	if c.Elapsed() < 5*time.Millisecond {
		t.Fatalf("expected elapsed >= 5ms, got %v", c.Elapsed())
	}
	if c.Count() != 1 {
		t.Fatalf("expected count 1, got %d", c.Count())
	}
	if c.LastScan().IsZero() {
		t.Fatal("expected non-zero LastScan")
	}
}

func TestClock_Stop_WithoutStart_ReturnsError(t *testing.T) {
	c := NewClock()
	if err := c.Stop(); err == nil {
		t.Fatal("expected error stopping without start")
	}
}

func TestClock_Start_WhileInProgress_ReturnsError(t *testing.T) {
	c := NewClock()
	_ = c.Start()
	if err := c.Start(); err == nil {
		t.Fatal("expected error starting while scan in progress")
	}
}

func TestClock_MultipleRounds_IncrementsCount(t *testing.T) {
	c := NewClock()
	for i := 0; i < 3; i++ {
		_ = c.Start()
		_ = c.Stop()
	}
	if c.Count() != 3 {
		t.Fatalf("expected count 3, got %d", c.Count())
	}
}

func TestClock_Elapsed_ReflectsMostRecent(t *testing.T) {
	c := NewClock()
	_ = c.Start()
	time.Sleep(5 * time.Millisecond)
	_ = c.Stop()
	first := c.Elapsed()

	_ = c.Start()
	time.Sleep(10 * time.Millisecond)
	_ = c.Stop()
	second := c.Elapsed()

	if second <= first {
		t.Fatalf("expected second elapsed > first, got %v <= %v", second, first)
	}
}
