package watch

import (
	"testing"
)

func TestNewMarker_Empty(t *testing.T) {
	m := NewMarker()
	if len(m.Snapshot()) != 0 {
		t.Fatal("expected empty marker")
	}
}

func TestMarker_Mark_Valid(t *testing.T) {
	m := NewMarker()
	if err := m.Mark(443, "tls audit"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !m.IsMarked(443) {
		t.Fatal("expected port to be marked")
	}
}

func TestMarker_Mark_InvalidPort(t *testing.T) {
	m := NewMarker()
	if err := m.Mark(0, "bad"); err == nil {
		t.Fatal("expected error for port 0")
	}
	if err := m.Mark(65536, "bad"); err == nil {
		t.Fatal("expected error for port 65536")
	}
}

func TestMarker_Mark_EmptyReason(t *testing.T) {
	m := NewMarker()
	if err := m.Mark(80, ""); err == nil {
		t.Fatal("expected error for empty reason")
	}
}

func TestMarker_Reason_ReturnsValue(t *testing.T) {
	m := NewMarker()
	_ = m.Mark(22, "ssh review")
	if got := m.Reason(22); got != "ssh review" {
		t.Fatalf("expected 'ssh review', got %q", got)
	}
}

func TestMarker_Reason_Missing(t *testing.T) {
	m := NewMarker()
	if got := m.Reason(9999); got != "" {
		t.Fatalf("expected empty string, got %q", got)
	}
}

func TestMarker_Unmark_RemovesMark(t *testing.T) {
	m := NewMarker()
	_ = m.Mark(8080, "check")
	m.Unmark(8080)
	if m.IsMarked(8080) {
		t.Fatal("expected port to be unmarked")
	}
}

func TestMarker_Snapshot_ReturnsCopy(t *testing.T) {
	m := NewMarker()
	_ = m.Mark(80, "http")
	_ = m.Mark(443, "https")
	snap := m.Snapshot()
	if len(snap) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(snap))
	}
	delete(snap, 80)
	if !m.IsMarked(80) {
		t.Fatal("snapshot mutation affected original")
	}
}
