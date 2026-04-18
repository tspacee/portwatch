package watch

import (
	"testing"
)

func TestNewTagRegistry_Empty(t *testing.T) {
	r := NewTagRegistry()
	if len(r.All()) != 0 {
		t.Fatal("expected empty registry")
	}
}

func TestTagRegistry_Set_And_Get(t *testing.T) {
	r := NewTagRegistry()
	if err := r.Set(80, "http"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	label, ok := r.Get(80)
	if !ok || label != "http" {
		t.Fatalf("expected label http, got %q ok=%v", label, ok)
	}
}

func TestTagRegistry_Set_InvalidPort(t *testing.T) {
	r := NewTagRegistry()
	if err := r.Set(0, "zero"); err == nil {
		t.Fatal("expected error for port 0")
	}
	if err := r.Set(70000, "big"); err == nil {
		t.Fatal("expected error for port 70000")
	}
}

func TestTagRegistry_Set_EmptyLabel(t *testing.T) {
	r := NewTagRegistry()
	if err := r.Set(443, ""); err == nil {
		t.Fatal("expected error for empty label")
	}
}

func TestTagRegistry_Get_Missing(t *testing.T) {
	r := NewTagRegistry()
	_, ok := r.Get(9999)
	if ok {
		t.Fatal("expected missing port to return false")
	}
}

func TestTagRegistry_Delete(t *testing.T) {
	r := NewTagRegistry()
	_ = r.Set(22, "ssh")
	r.Delete(22)
	_, ok := r.Get(22)
	if ok {
		t.Fatal("expected port to be deleted")
	}
}

func TestTagRegistry_All_ReturnsCopy(t *testing.T) {
	r := NewTagRegistry()
	_ = r.Set(80, "http")
	_ = r.Set(443, "https")
	all := r.All()
	if len(all) != 2 {
		t.Fatalf("expected 2 tags, got %d", len(all))
	}
	// mutate copy
	all[0].Label = "changed"
	label, _ := r.Get(all[0].Port)
	if label == "changed" {
		t.Fatal("All should return a copy")
	}
}
