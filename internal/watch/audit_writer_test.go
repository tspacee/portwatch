package watch

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func TestNewAuditWriter_NilWriter(t *testing.T) {
	aw := NewAuditWriter(nil)
	if aw == nil {
		t.Fatal("expected non-nil AuditWriter")
	}
}

func TestAuditWriter_Write_NilLog(t *testing.T) {
	var buf bytes.Buffer
	aw := NewAuditWriter(&buf)
	err := aw.Write(nil)
	if err == nil {
		t.Fatal("expected error for nil log")
	}
}

func TestAuditWriter_Write_EmptyLog(t *testing.T) {
	var buf bytes.Buffer
	aw := NewAuditWriter(&buf)
	al, _ := NewAuditLog(10)
	if err := aw.Write(al); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	output := buf.String()
	if !strings.Contains(output, "TIMESTAMP") {
		t.Error("expected header in output")
	}
}

func TestAuditWriter_Write_ContainsEntries(t *testing.T) {
	var buf bytes.Buffer
	aw := NewAuditWriter(&buf)
	al, _ := NewAuditLog(10)
	al.Add(AuditEntry{
		Timestamp: time.Now(),
		Port:      8080,
		Protocol:  "tcp",
		Event:     "opened",
		Source:    "scanner",
	})
	if err := aw.Write(al); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	output := buf.String()
	if !strings.Contains(output, "8080") {
		t.Errorf("expected port 8080 in output, got: %s", output)
	}
	if !strings.Contains(output, "opened") {
		t.Errorf("expected event 'opened' in output, got: %s", output)
	}
	if !strings.Contains(output, "scanner") {
		t.Errorf("expected source 'scanner' in output, got: %s", output)
	}
}
