package watch

import (
	"testing"
	"time"
)

func TestNewBeacon_InvalidInterval(t *testing.T) {
	_, err := NewBeacon(0)
	if err == nil {
		t.Fatal("expected error for zero interval")
	}
	_, err = NewBeacon(-time.Second)
	if err == nil {
		t.Fatal("expected error for negative interval")
	}
}

func TestNewBeacon_Valid(t *testing.T) {
	b, err := NewBeacon(10 * time.Millisecond)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if b == nil {
		t.Fatal("expected non-nil beacon")
	}
}

func TestBeacon_EmitsSignal(t *testing.T) {
	b, _ := NewBeacon(20 * time.Millisecond)
	b.Start()
	defer b.Stop()

	select {
	case <-b.C():
		// received signal
	case <-time.After(200 * time.Millisecond):
		t.Fatal("expected beacon signal, timed out")
	}
}

func TestBeacon_Stop_PreventsSignals(t *testing.T) {
	b, _ := NewBeacon(20 * time.Millisecond)
	b.Start()
	b.Stop()

	// drain any buffered signal
	select {
	case <-b.C():
	default:
	}

	time.Sleep(60 * time.Millisecond)
	select {
	case <-b.C():
		t.Fatal("beacon should not emit after stop")
	default:
	}
}

func TestBeacon_Running_ReflectsState(t *testing.T) {
	b, _ := NewBeacon(50 * time.Millisecond)
	if b.Running() {
		t.Fatal("should not be running before Start")
	}
	b.Start()
	if !b.Running() {
		t.Fatal("should be running after Start")
	}
	b.Stop()
	if b.Running() {
		t.Fatal("should not be running after Stop")
	}
}

func TestBeacon_Start_Idempotent(t *testing.T) {
	b, _ := NewBeacon(50 * time.Millisecond)
	b.Start()
	b.Start() // should not panic or spawn extra goroutines
	defer b.Stop()
	if !b.Running() {
		t.Fatal("expected beacon to be running")
	}
}
