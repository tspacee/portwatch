package watch

import (
	"testing"
	"time"
)

func TestNewCircuitBreaker_InvalidThreshold(t *testing.T) {
	_, err := NewCircuitBreaker(0, time.Second)
	if err == nil {
		t.Fatal("expected error for zero threshold")
	}
}

func TestNewCircuitBreaker_InvalidRecovery(t *testing.T) {
	_, err := NewCircuitBreaker(1, 0)
	if err == nil {
		t.Fatal("expected error for zero recovery")
	}
}

func TestNewCircuitBreaker_Valid(t *testing.T) {
	cb, err := NewCircuitBreaker(3, time.Second)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cb.State() != StateClosed {
		t.Errorf("expected StateClosed, got %v", cb.State())
	}
}

func TestCircuitBreaker_Allow_ClosedState(t *testing.T) {
	cb, _ := NewCircuitBreaker(3, time.Second)
	if err := cb.Allow(); err != nil {
		t.Errorf("expected nil, got %v", err)
	}
}

func TestCircuitBreaker_TripsAfterThreshold(t *testing.T) {
	cb, _ := NewCircuitBreaker(2, time.Second)
	cb.RecordFailure()
	if cb.State() != StateClosed {
		t.Error("expected still closed after one failure")
	}
	cb.RecordFailure()
	if cb.State() != StateOpen {
		t.Error("expected StateOpen after threshold reached")
	}
	if err := cb.Allow(); err != ErrCircuitOpen {
		t.Errorf("expected ErrCircuitOpen, got %v", err)
	}
}

func TestCircuitBreaker_RecoveryTransition(t *testing.T) {
	cb, _ := NewCircuitBreaker(1, 10*time.Millisecond)
	cb.RecordFailure()
	if err := cb.Allow(); err != ErrCircuitOpen {
		t.Fatal("expected open circuit")
	}
	time.Sleep(20 * time.Millisecond)
	if err := cb.Allow(); err != nil {
		t.Errorf("expected half-open to allow, got %v", err)
	}
	if cb.State() != StateHalfOpen {
		t.Errorf("expected StateHalfOpen, got %v", cb.State())
	}
}

func TestCircuitBreaker_RecordSuccess_Resets(t *testing.T) {
	cb, _ := NewCircuitBreaker(1, 10*time.Millisecond)
	cb.RecordFailure()
	time.Sleep(20 * time.Millisecond)
	_ = cb.Allow()
	cb.RecordSuccess()
	if cb.State() != StateClosed {
		t.Errorf("expected StateClosed after success, got %v", cb.State())
	}
}
