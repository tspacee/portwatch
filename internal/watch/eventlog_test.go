package watch

import (
	"testing"
	"time"
)

func makeEvent(t EventType, port int) PortEvent {
	return PortEvent{Type: t, Port: port, Protocol: "tcp", DetectedAt: time.Now()}
}

func TestNewEventLog_InvalidSize(t *testing.T) {
	_, err := NewEventLog(0)
	if err != ErrInvalidEventLogSize {
		t.Fatalf("expected ErrInvalidEventLogSize, got %v", err)
	}
}

func TestNewEventLog_Valid(t *testing.T) {
	l, err := NewEventLog(10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if l.Len() != 0 {
		t.Fatalf("expected empty log")
	}
}

func TestEventLog_Add_And_Len(t *testing.T) {
	l, _ := NewEventLog(5)
	l.Add(makeEvent(EventPortOpened, 8080))
	l.Add(makeEvent(EventPortClosed, 9090))
	if l.Len() != 2 {
		t.Fatalf("expected 2, got %d", l.Len())
	}
}

func TestEventLog_Entries_ReturnsCopy(t *testing.T) {
	l, _ := NewEventLog(5)
	l.Add(makeEvent(EventPortOpened, 3000))
	entries := l.Entries()
	entries[0].Port = 9999
	if l.Entries()[0].Port == 9999 {
		t.Fatal("entries should be a copy")
	}
}

func TestEventLog_Add_EvictsOldest(t *testing.T) {
	l, _ := NewEventLog(3)
	l.Add(makeEvent(EventPortOpened, 1))
	l.Add(makeEvent(EventPortOpened, 2))
	l.Add(makeEvent(EventPortOpened, 3))
	l.Add(makeEvent(EventPortOpened, 4))
	if l.Len() != 3 {
		t.Fatalf("expected 3, got %d", l.Len())
	}
	if l.Entries()[0].Port != 2 {
		t.Fatalf("expected oldest evicted, first port should be 2")
	}
}

func TestEventLog_EventTypes(t *testing.T) {
	l, _ := NewEventLog(10)
	l.Add(makeEvent(EventPortOpened, 80))
	l.Add(makeEvent(EventPortClosed, 80))
	entries := l.Entries()
	if entries[0].Type != EventPortOpened {
		t.Fatalf("expected opened")
	}
	if entries[1].Type != EventPortClosed {
		t.Fatalf("expected closed")
	}
}
