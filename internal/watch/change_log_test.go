package watch

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func makeChangeEntry(added, removed []int) ChangeEntry {
	return ChangeEntry{Timestamp: time.Now(), Added: added, Removed: removed}
}

func TestNewChangeLog_InvalidSize(t *testing.T) {
	_, err := NewChangeLog(0)
	if err == nil {
		t.Fatal("expected error for size 0")
	}
}

func TestNewChangeLog_Valid(t *testing.T) {
	cl, err := NewChangeLog(10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cl.Len() != 0 {
		t.Errorf("expected empty log")
	}
}

func TestChangeLog_Add_And_Len(t *testing.T) {
	cl, _ := NewChangeLog(5)
	cl.Add(makeChangeEntry([]int{80}, nil))
	cl.Add(makeChangeEntry([]int{443}, nil))
	if cl.Len() != 2 {
		t.Errorf("expected 2, got %d", cl.Len())
	}
}

func TestChangeLog_Add_Evicts_Oldest(t *testing.T) {
	cl, _ := NewChangeLog(2)
	cl.Add(makeChangeEntry([]int{1}, nil))
	cl.Add(makeChangeEntry([]int{2}, nil))
	cl.Add(makeChangeEntry([]int{3}, nil))
	if cl.Len() != 2 {
		t.Errorf("expected 2, got %d", cl.Len())
	}
	if cl.Entries()[0].Added[0] != 2 {
		t.Errorf("expected oldest evicted")
	}
}

func TestChangeLog_Entries_ReturnsCopy(t *testing.T) {
	cl, _ := NewChangeLog(5)
	cl.Add(makeChangeEntry([]int{80}, nil))
	e := cl.Entries()
	e[0].Added = []int{9999}
	if cl.Entries()[0].Added[0] == 9999 {
		t.Error("entries should be a copy")
	}
}

func TestChangeLog_Print_Output(t *testing.T) {
	cl, _ := NewChangeLog(5)
	cl.Add(makeChangeEntry([]int{80, 443}, []int{22}))
	var buf bytes.Buffer
	cl.Print(&buf)
	out := buf.String()
	if !strings.Contains(out, "80,443") {
		t.Errorf("expected added ports in output, got: %s", out)
	}
	if !strings.Contains(out, "22") {
		t.Errorf("expected removed port in output, got: %s", out)
	}
}

func TestChangeLog_Print_NilWriter(t *testing.T) {
	cl, _ := NewChangeLog(3)
	cl.Add(makeChangeEntry(nil, nil))
	// Should not panic with nil writer.
	cl.Print(nil)
}
