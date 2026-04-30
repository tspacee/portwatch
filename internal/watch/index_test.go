package watch

import (
	"testing"
)

func TestNewIndex_Empty(t *testing.T) {
	idx := NewIndex()
	if idx.Len() != 0 {
		t.Fatalf("expected 0 entries, got %d", idx.Len())
	}
}

func TestIndex_Set_Valid(t *testing.T) {
	idx := NewIndex()
	if err := idx.Set(80, "web"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if idx.Len() != 1 {
		t.Fatalf("expected 1 entry, got %d", idx.Len())
	}
}

func TestIndex_Set_InvalidPort(t *testing.T) {
	idx := NewIndex()
	if err := idx.Set(0, "web"); err == nil {
		t.Fatal("expected error for port 0")
	}
	if err := idx.Set(65536, "web"); err == nil {
		t.Fatal("expected error for port 65536")
	}
}

func TestIndex_Set_EmptyCategory(t *testing.T) {
	idx := NewIndex()
	if err := idx.Set(443, ""); err == nil {
		t.Fatal("expected error for empty category")
	}
}

func TestIndex_Get_Present(t *testing.T) {
	idx := NewIndex()
	_ = idx.Set(22, "ssh")
	cat, ok := idx.Get(22)
	if !ok {
		t.Fatal("expected entry to be present")
	}
	if cat != "ssh" {
		t.Fatalf("expected \"ssh\", got %q", cat)
	}
}

func TestIndex_Get_Missing(t *testing.T) {
	idx := NewIndex()
	_, ok := idx.Get(9999)
	if ok {
		t.Fatal("expected entry to be absent")
	}
}

func TestIndex_Delete_RemovesEntry(t *testing.T) {
	idx := NewIndex()
	_ = idx.Set(8080, "proxy")
	idx.Delete(8080)
	if idx.Len() != 0 {
		t.Fatalf("expected 0 entries after delete, got %d", idx.Len())
	}
}

func TestIndex_Snapshot_ReturnsCopy(t *testing.T) {
	idx := NewIndex()
	_ = idx.Set(80, "web")
	_ = idx.Set(443, "tls")
	snap := idx.Snapshot()
	if len(snap) != 2 {
		t.Fatalf("expected 2 entries in snapshot, got %d", len(snap))
	}
	// Mutating the snapshot must not affect the index.
	delete(snap, 80)
	if idx.Len() != 2 {
		t.Fatal("snapshot mutation affected the original index")
	}
}
