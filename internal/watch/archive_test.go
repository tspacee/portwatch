package watch

import (
	"testing"
)

func TestNewArchive_InvalidSize(t *testing.T) {
	_, err := NewArchive(0)
	if err == nil {
		t.Fatal("expected error for maxSize=0")
	}
}

func TestNewArchive_Valid(t *testing.T) {
	arch, err := NewArchive(10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if arch.Len() != 0 {
		t.Fatalf("expected 0 entries, got %d", arch.Len())
	}
}

func TestArchive_Store_And_Len(t *testing.T) {
	arch, _ := NewArchive(5)
	arch.Store([]int{80, 443})
	if arch.Len() != 1 {
		t.Fatalf("expected 1 entry, got %d", arch.Len())
	}
}

func TestArchive_Store_EvictsOldest(t *testing.T) {
	arch, _ := NewArchive(3)
	arch.Store([]int{80})
	arch.Store([]int{443})
	arch.Store([]int{8080})
	arch.Store([]int{9090})

	if arch.Len() != 3 {
		t.Fatalf("expected 3 entries, got %d", arch.Len())
	}
	entries := arch.Entries()
	if entries[0].Ports[0] != 443 {
		t.Fatalf("expected oldest evicted; first entry ports = %v", entries[0].Ports)
	}
}

func TestArchive_Entries_ReturnsCopy(t *testing.T) {
	arch, _ := NewArchive(5)
	arch.Store([]int{22, 80})

	entries := arch.Entries()
	entries[0].Ports[0] = 9999

	original := arch.Entries()
	if original[0].Ports[0] == 9999 {
		t.Fatal("Entries should return a deep copy")
	}
}

func TestArchive_Clear_RemovesAll(t *testing.T) {
	arch, _ := NewArchive(5)
	arch.Store([]int{80})
	arch.Store([]int{443})
	arch.Clear()

	if arch.Len() != 0 {
		t.Fatalf("expected 0 entries after Clear, got %d", arch.Len())
	}
}

func TestArchive_Entries_HasTimestamp(t *testing.T) {
	arch, _ := NewArchive(5)
	arch.Store([]int{80})

	entries := arch.Entries()
	if entries[0].Timestamp.IsZero() {
		t.Fatal("expected non-zero timestamp on archive entry")
	}
}
