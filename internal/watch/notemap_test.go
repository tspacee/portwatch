package watch

import (
	"testing"
)

func TestNewNoteMap_Empty(t *testing.T) {
	nm := NewNoteMap()
	if nm == nil {
		t.Fatal("expected non-nil NoteMap")
	}
	if got := nm.Get(80); got != "" {
		t.Fatalf("expected empty note, got %q", got)
	}
}

func TestNoteMap_Set_And_Get(t *testing.T) {
	nm := NewNoteMap()
	if err := nm.Set(443, "HTTPS"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := nm.Get(443); got != "HTTPS" {
		t.Fatalf("expected 'HTTPS', got %q", got)
	}
}

func TestNoteMap_Set_InvalidPort(t *testing.T) {
	nm := NewNoteMap()
	if err := nm.Set(0, "bad"); err != ErrInvalidNotePort {
		t.Fatalf("expected ErrInvalidNotePort, got %v", err)
	}
	if err := nm.Set(65536, "bad"); err != ErrInvalidNotePort {
		t.Fatalf("expected ErrInvalidNotePort, got %v", err)
	}
}

func TestNoteMap_Set_EmptyNote(t *testing.T) {
	nm := NewNoteMap()
	if err := nm.Set(80, ""); err != ErrEmptyNote {
		t.Fatalf("expected ErrEmptyNote, got %v", err)
	}
}

func TestNoteMap_Delete_RemovesNote(t *testing.T) {
	nm := NewNoteMap()
	_ = nm.Set(22, "SSH")
	nm.Delete(22)
	if got := nm.Get(22); got != "" {
		t.Fatalf("expected empty after delete, got %q", got)
	}
}

func TestNoteMap_Snapshot_ReturnsCopy(t *testing.T) {
	nm := NewNoteMap()
	_ = nm.Set(80, "HTTP")
	_ = nm.Set(443, "HTTPS")
	snap := nm.Snapshot()
	if len(snap) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(snap))
	}
	// Mutating snapshot must not affect NoteMap
	snap[80] = "modified"
	if got := nm.Get(80); got != "HTTP" {
		t.Fatalf("snapshot mutation affected NoteMap: got %q", got)
	}
}
