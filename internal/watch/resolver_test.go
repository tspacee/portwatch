package watch

import (
	"testing"
)

func TestNewResolver_ContainsWellKnown(t *testing.T) {
	r := NewResolver()
	if r.Resolve(80) != "http" {
		t.Errorf("expected http, got %s", r.Resolve(80))
	}
	if r.Resolve(443) != "https" {
		t.Errorf("expected https, got %s", r.Resolve(443))
	}
}

func TestResolver_Resolve_Unknown(t *testing.T) {
	r := NewResolver()
	if got := r.Resolve(9999); got != "unknown" {
		t.Errorf("expected unknown, got %s", got)
	}
}

func TestResolver_Register_Valid(t *testing.T) {
	r := NewResolver()
	if err := r.Register(9090, "prometheus"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := r.Resolve(9090); got != "prometheus" {
		t.Errorf("expected prometheus, got %s", got)
	}
}

func TestResolver_Register_InvalidPort(t *testing.T) {
	r := NewResolver()
	if err := r.Register(0, "zero"); err == nil {
		t.Error("expected error for port 0")
	}
	if err := r.Register(65536, "overflow"); err == nil {
		t.Error("expected error for port 65536")
	}
}

func TestResolver_Register_EmptyService(t *testing.T) {
	r := NewResolver()
	if err := r.Register(8080, ""); err == nil {
		t.Error("expected error for empty service name")
	}
}

func TestResolver_Len_IncludesWellKnown(t *testing.T) {
	r := NewResolver()
	base := r.Len()
	_ = r.Register(19999, "custom")
	if r.Len() != base+1 {
		t.Errorf("expected %d, got %d", base+1, r.Len())
	}
}

func TestResolver_Register_Overwrite(t *testing.T) {
	r := NewResolver()
	_ = r.Register(22, "custom-ssh")
	if got := r.Resolve(22); got != "custom-ssh" {
		t.Errorf("expected custom-ssh, got %s", got)
	}
}

func TestResolver_Register_Overwrite_DoesNotIncreaseLen(t *testing.T) {
	r := NewResolver()
	base := r.Len()
	_ = r.Register(22, "custom-ssh")
	if r.Len() != base {
		t.Errorf("expected Len to remain %d after overwrite, got %d", base, r.Len())
	}
}
