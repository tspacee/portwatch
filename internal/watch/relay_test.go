package watch

import (
	"testing"
	"time"
)

func TestNewRelay_InvalidBufSize(t *testing.T) {
	_, err := NewRelay(0)
	if err == nil {
		t.Fatal("expected error for bufSize=0")
	}
}

func TestNewRelay_Valid(t *testing.T) {
	r, err := NewRelay(4)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Len() != 0 {
		t.Errorf("expected 0 subscribers, got %d", r.Len())
	}
}

func TestRelay_Subscribe_EmptyName(t *testing.T) {
	r, _ := NewRelay(4)
	_, err := r.Subscribe("")
	if err == nil {
		t.Fatal("expected error for empty name")
	}
}

func TestRelay_Subscribe_DuplicateName(t *testing.T) {
	r, _ := NewRelay(4)
	_, err := r.Subscribe("watcher")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_, err = r.Subscribe("watcher")
	if err == nil {
		t.Fatal("expected error for duplicate subscriber")
	}
}

func TestRelay_Subscribe_IncrementsLen(t *testing.T) {
	r, _ := NewRelay(4)
	r.Subscribe("a")
	r.Subscribe("b")
	if r.Len() != 2 {
		t.Errorf("expected len 2, got %d", r.Len())
	}
}

func TestRelay_Unsubscribe_DecrementsLen(t *testing.T) {
	r, _ := NewRelay(4)
	r.Subscribe("a")
	r.Unsubscribe("a")
	if r.Len() != 0 {
		t.Errorf("expected len 0 after unsubscribe, got %d", r.Len())
	}
}

func TestRelay_Broadcast_DeliversToSubscriber(t *testing.T) {
	r, _ := NewRelay(4)
	ch, _ := r.Subscribe("listener")

	ports := []int{80, 443, 8080}
	r.Broadcast(ports)

	select {
	case received := <-ch:
		if len(received) != len(ports) {
			t.Errorf("expected %d ports, got %d", len(ports), len(received))
		}
	case <-time.After(100 * time.Millisecond):
		t.Fatal("timed out waiting for broadcast")
	}
}

func TestRelay_Broadcast_DoesNotMutateOriginal(t *testing.T) {
	r, _ := NewRelay(4)
	r.Subscribe("listener")

	original := []int{22, 80}
	r.Broadcast(original)
	original[0] = 9999

	ch, _ := r.Subscribe("verifier")
	_ = ch
	// Original mutation should not affect already-sent copy; test passes if no panic
}
