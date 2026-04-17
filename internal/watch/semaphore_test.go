package watch

import (
	"sync"
	"testing"
)

func TestNewSemaphore_InvalidLimit(t *testing.T) {
	_, err := NewSemaphore(0)
	if err != ErrInvalidConcurrency {
		t.Fatalf("expected ErrInvalidConcurrency, got %v", err)
	}
}

func TestNewSemaphore_Valid(t *testing.T) {
	s, err := NewSemaphore(3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.Available() != 3 {
		t.Fatalf("expected 3 available, got %d", s.Available())
	}
}

func TestSemaphore_AcquireRelease(t *testing.T) {
	s, _ := NewSemaphore(2)
	s.Acquire()
	if s.Available() != 1 {
		t.Fatalf("expected 1 available after acquire, got %d", s.Available())
	}
	s.Release()
	if s.Available() != 2 {
		t.Fatalf("expected 2 available after release, got %d", s.Available())
	}
}

func TestSemaphore_TryAcquire_Success(t *testing.T) {
	s, _ := NewSemaphore(1)
	if !s.TryAcquire() {
		t.Fatal("expected TryAcquire to succeed")
	}
	if s.Available() != 0 {
		t.Fatalf("expected 0 available, got %d", s.Available())
	}
}

func TestSemaphore_TryAcquire_Fails_WhenFull(t *testing.T) {
	s, _ := NewSemaphore(1)
	s.Acquire()
	if s.TryAcquire() {
		t.Fatal("expected TryAcquire to fail when semaphore is full")
	}
}

func TestSemaphore_ConcurrentAcquire(t *testing.T) {
	const limit = 3
	const workers = 9
	s, _ := NewSemaphore(limit)
	var wg sync.WaitGroup
	counter := make(chan struct{}, workers)
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			s.Acquire()
			counter <- struct{}{}
			s.Release()
		}()
	}
	wg.Wait()
	if len(counter) != workers {
		t.Fatalf("expected %d completions, got %d", workers, len(counter))
	}
}
