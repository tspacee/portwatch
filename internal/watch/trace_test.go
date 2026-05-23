package watch

import (
	"testing"
	"time"
)

func TestNewTrace_InvalidMax(t *testing.T) {
	_, err := NewTrace(0)
	if err == nil {
		t.Fatal("expected error for zero maxEntries")
	}
}

func TestNewTrace_Valid(t *testing.T) {
	tr, err := NewTrace(10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tr.Len() != 0 {
		t.Fatalf("expected empty trace, got %d", tr.Len())
	}
}

func TestTrace_Record_InvalidPort(t *testing.T) {
	tr, _ := NewTrace(5)
	if err := tr.Record(0, true); err == nil {
		t.Fatal("expected error for port 0")
	}
	if err := tr.Record(65536, false); err == nil {
		t.Fatal("expected error for port 65536")
	}
}

func TestTrace_Record_And_Len(t *testing.T) {
	tr, _ := NewTrace(5)
	if err := tr.Record(80, true); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tr.Len() != 1 {
		t.Fatalf("expected 1 entry, got %d", tr.Len())
	}
}

func TestTrace_Record_EvictsOldest(t *testing.T) {
	tr, _ := NewTrace(3)
	for _, p := range []int{80, 443, 8080, 9090} {
		_ = tr.Record(p, true)
	}
	if tr.Len() != 3 {
		t.Fatalf("expected 3 entries after eviction, got %d", tr.Len())
	}
	entries := tr.Entries()
	if entries[0].Port != 443 {
		t.Fatalf("expected oldest surviving entry port 443, got %d", entries[0].Port)
	}
}

func TestTrace_Entries_ReturnsCopy(t *testing.T) {
	tr, _ := NewTrace(5)
	_ = tr.Record(22, false)
	e1 := tr.Entries()
	e1[0].Port = 9999
	e2 := tr.Entries()
	if e2[0].Port == 9999 {
		t.Fatal("Entries should return an independent copy")
	}
}

func TestTrace_Record_SetsTimestamp(t *testing.T) {
	before := time.Now()
	tr, _ := NewTrace(5)
	_ = tr.Record(8080, true)
	after := time.Now()
	entry := tr.Entries()[0]
	if entry.ObservedAt.Before(before) || entry.ObservedAt.After(after) {
		t.Fatal("ObservedAt timestamp out of expected range")
	}
}

func TestTrace_Clear_ResetsEntries(t *testing.T) {
	tr, _ := NewTrace(5)
	_ = tr.Record(80, true)
	_ = tr.Record(443, true)
	tr.Clear()
	if tr.Len() != 0 {
		t.Fatalf("expected 0 entries after Clear, got %d", tr.Len())
	}
}
