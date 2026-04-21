package watch_test

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/watch"
)

func makeLedgerEntry(port int, seen time.Time) watch.LedgerEntry {
	return watch.LedgerEntry{
		Port:     port,
		FirstSeen: seen,
		LastSeen:  seen,
		SeenCount: 1,
	}
}

func TestNewLedger_Empty(t *testing.T) {
	l := watch.NewLedger()
	if l == nil {
		t.Fatal("expected non-nil ledger")
	}
	if l.Len() != 0 {
		t.Fatalf("expected empty ledger, got len=%d", l.Len())
	}
}

func TestLedger_Record_Valid(t *testing.T) {
	l := watch.NewLedger()
	now := time.Now()

	err := l.Record(8080, now)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if l.Len() != 1 {
		t.Fatalf("expected len=1, got %d", l.Len())
	}
}

func TestLedger_Record_InvalidPort(t *testing.T) {
	l := watch.NewLedger()
	err := l.Record(0, time.Now())
	if err == nil {
		t.Fatal("expected error for invalid port")
	}
}

func TestLedger_Record_IncrementsSeenCount(t *testing.T) {
	l := watch.NewLedger()
	now := time.Now()

	_ = l.Record(443, now)
	_ = l.Record(443, now.Add(time.Second))
	_ = l.Record(443, now.Add(2*time.Second))

	entry, ok := l.Get(443)
	if !ok {
		t.Fatal("expected entry for port 443")
	}
	if entry.SeenCount != 3 {
		t.Fatalf("expected SeenCount=3, got %d", entry.SeenCount)
	}
}

func TestLedger_Record_UpdatesLastSeen(t *testing.T) {
	l := watch.NewLedger()
	first := time.Now()
	second := first.Add(5 * time.Second)

	_ = l.Record(22, first)
	_ = l.Record(22, second)

	entry, ok := l.Get(22)
	if !ok {
		t.Fatal("expected entry for port 22")
	}
	if !entry.LastSeen.Equal(second) {
		t.Fatalf("expected LastSeen=%v, got %v", second, entry.LastSeen)
	}
	if !entry.FirstSeen.Equal(first) {
		t.Fatalf("expected FirstSeen=%v, got %v", first, entry.FirstSeen)
	}
}

func TestLedger_Get_Missing(t *testing.T) {
	l := watch.NewLedger()
	_, ok := l.Get(9999)
	if ok {
		t.Fatal("expected missing entry")
	}
}

func TestLedger_Delete_RemovesEntry(t *testing.T) {
	l := watch.NewLedger()
	_ = l.Record(3000, time.Now())

	l.Delete(3000)

	_, ok := l.Get(3000)
	if ok {
		t.Fatal("expected entry to be deleted")
	}
	if l.Len() != 0 {
		t.Fatalf("expected len=0 after delete, got %d", l.Len())
	}
}

func TestLedger_Snapshot_ReturnsCopy(t *testing.T) {
	l := watch.NewLedger()
	_ = l.Record(80, time.Now())
	_ = l.Record(443, time.Now())

	snap := l.Snapshot()
	if len(snap) != 2 {
		t.Fatalf("expected snapshot len=2, got %d", len(snap))
	}

	// Mutating snapshot should not affect ledger
	delete(snap, 80)
	if l.Len() != 2 {
		t.Fatal("mutating snapshot should not affect ledger")
	}
}
