package watch_test

import (
	"testing"

	"github.com/user/portwatch/internal/watch"
)

func TestNewLabeler_Empty(t *testing.T) {
	l := watch.NewLabeler()
	if l == nil {
		t.Fatal("expected non-nil Labeler")
	}
	if got := len(l.All()); got != 0 {
		t.Fatalf("expected 0 labels, got %d", got)
	}
}

func TestLabeler_Register_Valid(t *testing.T) {
	l := watch.NewLabeler()
	if err := l.Register(80, "http"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := l.Label(80); got != "http" {
		t.Fatalf("expected 'http', got %q", got)
	}
}

func TestLabeler_Register_InvalidPort(t *testing.T) {
	l := watch.NewLabeler()
	if err := l.Register(0, "zero"); err == nil {
		t.Fatal("expected error for port 0")
	}
	if err := l.Register(65536, "toobig"); err == nil {
		t.Fatal("expected error for port 65536")
	}
}

func TestLabeler_Register_EmptyLabel(t *testing.T) {
	l := watch.NewLabeler()
	if err := l.Register(443, ""); err == nil {
		t.Fatal("expected error for empty label")
	}
}

func TestLabeler_Label_DefaultFallback(t *testing.T) {
	l := watch.NewLabeler()
	got := l.Label(9999)
	if got != "port/9999" {
		t.Fatalf("expected 'port/9999', got %q", got)
	}
}

func TestLabeler_All_ReturnsCopy(t *testing.T) {
	l := watch.NewLabeler()
	_ = l.Register(22, "ssh")
	all := l.All()
	all[22] = "mutated"
	if l.Label(22) != "ssh" {
		t.Fatal("All() should return a copy, not a reference")
	}
}
