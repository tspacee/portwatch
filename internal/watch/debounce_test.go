package watch

import (
	"sync/atomic"
	"testing"
	"time"
)

func TestNewDebounce_InvalidDelay(t *testing.T) {
	_, err := NewDebounce(0)
	if err == nil {
		t.Fatal("expected error for zero delay")
	}
	_, err = NewDebounce(-1 * time.Millisecond)
	if err == nil {
		t.Fatal("expected error for negative delay")
	}
}

func TestNewDebounce_Valid(t *testing.T) {
	d, err := NewDebounce(10 * time.Millisecond)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if d == nil {
		t.Fatal("expected non-nil Debounce")
	}
	d.Stop()
}

func TestDebounce_Trigger_CallsFn(t *testing.T) {
	d, _ := NewDebounce(20 * time.Millisecond)
	defer d.Stop()

	var called int32
	d.Trigger(func() { atomic.StoreInt32(&called, 1) })

	time.Sleep(50 * time.Millisecond)
	if atomic.LoadInt32(&called) != 1 {
		t.Error("expected fn to be called after delay")
	}
}

func TestDebounce_Trigger_Resets(t *testing.T) {
	d, _ := NewDebounce(40 * time.Millisecond)
	defer d.Stop()

	var count int32
	fn := func() { atomic.AddInt32(&count, 1) }

	// Trigger multiple times rapidly; only one call should result.
	for i := 0; i < 5; i++ {
		d.Trigger(fn)
		time.Sleep(10 * time.Millisecond)
	}

	time.Sleep(80 * time.Millisecond)
	if c := atomic.LoadInt32(&count); c != 1 {
		t.Errorf("expected 1 call, got %d", c)
	}
}

func TestDebounce_Stop_PreventsFn(t *testing.T) {
	d, _ := NewDebounce(50 * time.Millisecond)

	var called int32
	d.Trigger(func() { atomic.StoreInt32(&called, 1) })
	d.Stop()

	time.Sleep(80 * time.Millisecond)
	if atomic.LoadInt32(&called) != 0 {
		t.Error("expected fn to be suppressed after Stop")
	}
}

func TestDebounce_Stop_Idempotent(t *testing.T) {
	d, _ := NewDebounce(10 * time.Millisecond)
	d.Stop()
	d.Stop() // should not panic
}
