package watch

import (
	"testing"
)

func TestNewMatcher_Empty(t *testing.T) {
	m := NewMatcher()
	if m == nil {
		t.Fatal("expected non-nil Matcher")
	}
	if len(m.Snapshot()) != 0 {
		t.Error("expected empty snapshot")
	}
}

func TestMatcher_Add_Valid(t *testing.T) {
	m := NewMatcher()
	if err := m.Add(80, "http"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	labels := m.Match(80)
	if len(labels) != 1 || labels[0] != "http" {
		t.Errorf("expected [http], got %v", labels)
	}
}

func TestMatcher_Add_InvalidPort(t *testing.T) {
	m := NewMatcher()
	if err := m.Add(0, "zero"); err == nil {
		t.Error("expected error for port 0")
	}
	if err := m.Add(65536, "high"); err == nil {
		t.Error("expected error for port 65536")
	}
}

func TestMatcher_Add_EmptyLabel(t *testing.T) {
	m := NewMatcher()
	if err := m.Add(443, ""); err == nil {
		t.Error("expected error for empty label")
	}
}

func TestMatcher_Add_MultipleLabels(t *testing.T) {
	m := NewMatcher()
	_ = m.Add(443, "https")
	_ = m.Add(443, "tls")
	labels := m.Match(443)
	if len(labels) != 2 {
		t.Errorf("expected 2 labels, got %d", len(labels))
	}
}

func TestMatcher_Match_Missing(t *testing.T) {
	m := NewMatcher()
	if labels := m.Match(9999); labels != nil {
		t.Errorf("expected nil, got %v", labels)
	}
}

func TestMatcher_Has(t *testing.T) {
	m := NewMatcher()
	_ = m.Add(22, "ssh")
	if !m.Has(22) {
		t.Error("expected Has(22) to be true")
	}
	if m.Has(23) {
		t.Error("expected Has(23) to be false")
	}
}

func TestMatcher_Remove(t *testing.T) {
	m := NewMatcher()
	_ = m.Add(8080, "proxy")
	m.Remove(8080)
	if m.Has(8080) {
		t.Error("expected port to be removed")
	}
}

func TestMatcher_Snapshot_ReturnsCopy(t *testing.T) {
	m := NewMatcher()
	_ = m.Add(53, "dns")
	snap := m.Snapshot()
	snap[53] = append(snap[53], "injected")
	if len(m.Match(53)) != 1 {
		t.Error("snapshot mutation affected internal state")
	}
}
