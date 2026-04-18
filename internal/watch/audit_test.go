package watch

import (
	"testing"
	"time"
)

func makeAuditEntry(port int, event string) AuditEntry {
	return AuditEntry{
		Timestamp: time.Now(),
		Port:      port,
		Protocol:  "tcp",
		Event:     event,
		Source:    "test",
	}
}

func TestNewAuditLog_InvalidSize(t *testing.T) {
	_, err := NewAuditLog(0)
	if err == nil {
		t.Fatal("expected error for size 0")
	}
}

func TestNewAuditLog_Valid(t *testing.T) {
	al, err := NewAuditLog(10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if al.Len() != 0 {
		t.Errorf("expected 0 entries, got %d", al.Len())
	}
}

func TestAuditLog_Add_And_Len(t *testing.T) {
	al, _ := NewAuditLog(10)
	al.Add(makeAuditEntry(80, "opened"))
	al.Add(makeAuditEntry(443, "opened"))
	if al.Len() != 2 {
		t.Errorf("expected 2, got %d", al.Len())
	}
}

func TestAuditLog_Add_Evicts_Oldest(t *testing.T) {
	al, _ := NewAuditLog(2)
	al.Add(makeAuditEntry(80, "opened"))
	al.Add(makeAuditEntry(443, "opened"))
	al.Add(makeAuditEntry(8080, "opened"))
	if al.Len() != 2 {
		t.Errorf("expected 2 after eviction, got %d", al.Len())
	}
	entries := al.Entries()
	if entries[0].Port != 443 {
		t.Errorf("expected oldest evicted, got port %d", entries[0].Port)
	}
}

func TestAuditLog_Entries_ReturnsCopy(t *testing.T) {
	al, _ := NewAuditLog(10)
	al.Add(makeAuditEntry(80, "opened"))
	e := al.Entries()
	e[0].Port = 9999
	if al.Entries()[0].Port == 9999 {
		t.Error("expected copy, but original was mutated")
	}
}

func TestAuditLog_Clear(t *testing.T) {
	al, _ := NewAuditLog(10)
	al.Add(makeAuditEntry(80, "opened"))
	al.Clear()
	if al.Len() != 0 {
		t.Errorf("expected 0 after clear, got %d", al.Len())
	}
}
