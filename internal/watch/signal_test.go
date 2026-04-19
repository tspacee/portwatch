package watch

import (
	"testing"
	"time"
)

func TestNewSignal_Empty(t *testing.T) {
	s := NewSignal()
	if s == nil {
		t.Fatal("expected non-nil Signal")
	}
	if s.Len("boot") != 0 {
		t.Fatal("expected zero subscribers")
	}
}

func TestSignal_Subscribe_EmptyEvent_ReturnsError(t *testing.T) {
	s := NewSignal()
	_, err := s.Subscribe("")
	if err == nil {
		t.Fatal("expected error for empty event")
	}
}

func TestSignal_Fire_EmptyEvent_ReturnsError(t *testing.T) {
	s := NewSignal()
	if err := s.Fire(""); err == nil {
		t.Fatal("expected error for empty event")
	}
}

func TestSignal_Subscribe_IncrementsLen(t *testing.T) {
	s := NewSignal()
	s.Subscribe("ready") //nolint
	s.Subscribe("ready") //nolint
	if got := s.Len("ready"); got != 2 {
		t.Fatalf("expected 2 subscribers, got %d", got)
	}
}

func TestSignal_Fire_NotifiesSubscriber(t *testing.T) {
	s := NewSignal()
	ch, err := s.Subscribe("ready")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := s.Fire("ready"); err != nil {
		t.Fatalf("unexpected fire error: %v", err)
	}
	select {
	case <-ch:
		// ok
	case <-time.After(100 * time.Millisecond):
		t.Fatal("subscriber channel was not closed after Fire")
	}
}

func TestSignal_Fire_ClearsSubscribers(t *testing.T) {
	s := NewSignal()
	s.Subscribe("done") //nolint
	s.Fire("done")       //nolint
	if got := s.Len("done"); got != 0 {
		t.Fatalf("expected 0 subscribers after fire, got %d", got)
	}
}

func TestSignal_Fire_UnknownEvent_NoError(t *testing.T) {
	s := NewSignal()
	if err := s.Fire("unknown"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
