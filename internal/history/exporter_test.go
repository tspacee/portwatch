package history

import (
	"bytes"
	"encoding/csv"
	"strings"
	"testing"
	"time"
)

func TestNewExporter_NilHistory(t *testing.T) {
	_, err := NewExporter(nil)
	if err == nil {
		t.Fatal("expected error for nil history, got nil")
	}
}

func TestNewExporter_Valid(t *testing.T) {
	h := New(10)
	ex, err := NewExporter(h)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ex == nil {
		t.Fatal("expected non-nil exporter")
	}
}

func TestWriteCSV_HeaderOnly(t *testing.T) {
	h := New(10)
	ex, _ := NewExporter(h)

	var buf bytes.Buffer
	if err := ex.WriteCSV(&buf); err != nil {
		t.Fatalf("WriteCSV error: %v", err)
	}

	r := csv.NewReader(&buf)
	records, err := r.ReadAll()
	if err != nil {
		t.Fatalf("csv parse error: %v", err)
	}
	if len(records) != 1 {
		t.Fatalf("expected 1 row (header), got %d", len(records))
	}
	if records[0][0] != "timestamp" {
		t.Errorf("expected first header 'timestamp', got %q", records[0][0])
	}
}

func TestWriteCSV_ContainsEntries(t *testing.T) {
	h := New(10)
	h.Add(Entry{
		Timestamp: time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
		Port:      8080,
		Protocol:  "tcp",
		Event:     "opened",
		Detail:    "new port detected",
	})

	ex, _ := NewExporter(h)
	var buf bytes.Buffer
	if err := ex.WriteCSV(&buf); err != nil {
		t.Fatalf("WriteCSV error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "8080") {
		t.Errorf("expected port 8080 in output, got: %s", output)
	}
	if !strings.Contains(output, "opened") {
		t.Errorf("expected event 'opened' in output, got: %s", output)
	}
	if !strings.Contains(output, "2024-01-15T10:00:00Z") {
		t.Errorf("expected RFC3339 timestamp in output, got: %s", output)
	}

	r := csv.NewReader(&buf)
	records, _ := r.ReadAll()
	// header + 1 entry
	if len(records) != 2 {
		t.Fatalf("expected 2 rows, got %d", len(records))
	}
}
