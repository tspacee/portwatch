package watch

import (
	"context"
	"errors"
	"testing"
	"time"
)

func makeBackoff(t *testing.T) *Backoff {
	t.Helper()
	b, err := NewBackoff(10*time.Millisecond, 50*time.Millisecond, 2.0)
	if err != nil {
		t.Fatalf("NewBackoff: %v", err)
	}
	return b
}

func TestNewSupervisor_NilWorker(t *testing.T) {
	_, err := NewSupervisor(nil, makeBackoff(t), 3)
	if err == nil {
		t.Fatal("expected error for nil worker")
	}
}

func TestNewSupervisor_NilBackoff(t *testing.T) {
	_, err := NewSupervisor(func(ctx context.Context) error { return nil }, nil, 3)
	if err == nil {
		t.Fatal("expected error for nil backoff")
	}
}

func TestNewSupervisor_Valid(t *testing.T) {
	sup, err := NewSupervisor(func(ctx context.Context) error { return nil }, makeBackoff(t), 3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if sup == nil {
		t.Fatal("expected non-nil supervisor")
	}
}

func TestSupervisor_Run_SucceedsImmediately(t *testing.T) {
	sup, _ := NewSupervisor(func(ctx context.Context) error { return nil }, makeBackoff(t), 3)
	if err := sup.Run(context.Background()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if sup.Restarts() != 0 {
		t.Fatalf("expected 0 restarts, got %d", sup.Restarts())
	}
}

func TestSupervisor_Run_RestartsOnError(t *testing.T) {
	calls := 0
	worker := func(ctx context.Context) error {
		calls++
		if calls < 3 {
			return errors.New("transient")
		}
		return nil
	}
	sup, _ := NewSupervisor(worker, makeBackoff(t), 5)
	if err := sup.Run(context.Background()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if sup.Restarts() != 2 {
		t.Fatalf("expected 2 restarts, got %d", sup.Restarts())
	}
}

func TestSupervisor_Run_StopsAtMaxRestarts(t *testing.T) {
	worker := func(ctx context.Context) error { return errors.New("fail") }
	sup, _ := NewSupervisor(worker, makeBackoff(t), 2)
	err := sup.Run(context.Background())
	if !errors.Is(err, ErrSupervisorStopped) {
		t.Fatalf("expected ErrSupervisorStopped, got %v", err)
	}
}

func TestSupervisor_Run_StopsOnContextCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	worker := func(ctx context.Context) error {
		cancel()
		return errors.New("fail")
	}
	sup, _ := NewSupervisor(worker, makeBackoff(t), 0)
	err := sup.Run(ctx)
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected context.Canceled, got %v", err)
	}
}
