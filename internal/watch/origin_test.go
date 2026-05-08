package watch

import (
	"testing"
)

func TestNewOrigin_Empty(t *testing.T) {
	o := NewOrigin()
	if o.Len() != 0 {
		t.Fatalf("expected 0, got %d", o.Len())
	}
}

func TestOrigin_Record_Valid(t *testing.T) {
	o := NewOrigin()
	if err := o.Record(80, "scanner"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if o.Len() != 1 {
		t.Fatalf("expected 1, got %d", o.Len())
	}
}

func TestOrigin_Record_InvalidPort(t *testing.T) {
	o := NewOrigin()
	if err := o.Record(0, "scanner"); err == nil {
		t.Fatal("expected error for port 0")
	}
	if err := o.Record(65536, "scanner"); err == nil {
		t.Fatal("expected error for port 65536")
	}
}

func TestOrigin_Record_EmptySource(t *testing.T) {
	o := NewOrigin()
	if err := o.Record(443, ""); err == nil {
		t.Fatal("expected error for empty source")
	}
}

func TestOrigin_Record_Idempotent(t *testing.T) {
	o := NewOrigin()
	_ = o.Record(22, "first")
	_ = o.Record(22, "second")
	e, ok := o.Get(22)
	if !ok {
		t.Fatal("expected entry")
	}
	if e.Source != "first" {
		t.Fatalf("expected 'first', got %q", e.Source)
	}
}

func TestOrigin_Get_Present(t *testing.T) {
	o := NewOrigin()
	_ = o.Record(8080, "daemon")
	e, ok := o.Get(8080)
	if !ok {
		t.Fatal("expected entry to be present")
	}
	if e.Port != 8080 || e.Source != "daemon" {
		t.Fatalf("unexpected entry: %+v", e)
	}
}

func TestOrigin_Get_Missing(t *testing.T) {
	o := NewOrigin()
	_, ok := o.Get(9999)
	if ok {
		t.Fatal("expected missing entry")
	}
}

func TestOrigin_Reset_ClearsAll(t *testing.T) {
	o := NewOrigin()
	_ = o.Record(80, "a")
	_ = o.Record(443, "b")
	o.Reset()
	if o.Len() != 0 {
		t.Fatalf("expected 0 after reset, got %d", o.Len())
	}
}
