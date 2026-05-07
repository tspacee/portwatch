package watch

import (
	"testing"
)

func TestNewIntent_Empty(t *testing.T) {
	i := NewIntent()
	if i.Len() != 0 {
		t.Fatalf("expected 0, got %d", i.Len())
	}
}

func TestIntent_Declare_Valid(t *testing.T) {
	i := NewIntent()
	if err := i.Declare(80, "http server"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if i.Len() != 1 {
		t.Fatalf("expected 1, got %d", i.Len())
	}
}

func TestIntent_Declare_InvalidPort(t *testing.T) {
	i := NewIntent()
	if err := i.Declare(0, "bad"); err == nil {
		t.Fatal("expected error for port 0")
	}
	if err := i.Declare(65536, "bad"); err == nil {
		t.Fatal("expected error for port 65536")
	}
}

func TestIntent_Declare_EmptyIntent(t *testing.T) {
	i := NewIntent()
	if err := i.Declare(443, ""); err == nil {
		t.Fatal("expected error for empty intent")
	}
}

func TestIntent_Lookup_Present(t *testing.T) {
	i := NewIntent()
	_ = i.Declare(22, "ssh access")
	v, ok := i.Lookup(22)
	if !ok {
		t.Fatal("expected intent to be found")
	}
	if v != "ssh access" {
		t.Fatalf("expected 'ssh access', got %q", v)
	}
}

func TestIntent_Lookup_Missing(t *testing.T) {
	i := NewIntent()
	_, ok := i.Lookup(9999)
	if ok {
		t.Fatal("expected no intent for unseen port")
	}
}

func TestIntent_Revoke_RemovesEntry(t *testing.T) {
	i := NewIntent()
	_ = i.Declare(8080, "dev server")
	i.Revoke(8080)
	if i.Len() != 0 {
		t.Fatalf("expected 0 after revoke, got %d", i.Len())
	}
	_, ok := i.Lookup(8080)
	if ok {
		t.Fatal("expected intent to be absent after revoke")
	}
}
