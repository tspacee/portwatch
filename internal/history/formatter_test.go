package history

import (
	"strings"
	"testing"
	"time"

	"github.com/user/portwatch/internal/portscanner"
)

func makeFormatter() *Formatter {
	return NewFormatter("2006-01-02T15:04:05Z07:00")
}

func TestNewFormatter_DefaultTimeFormat(t *testing.T) {
	f := NewFormatter("")
	if f.timeFormat == "" {
		t.Fatal("expected non-empty default time format")
	}
}

func TestFormatEntry_ContainsTimestamp(t *testing.T) {
	f := makeFormatter()
	e := Entry{
		Timestamp:      time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
		ScanDurationMs: 42,
		Ports:          []portscanner.Port{{Port: 80, Protocol: "tcp"}},
	}
	result := f.FormatEntry(e)
	if !strings.Contains(result, "2024-01-15") {
		t.Errorf("expected timestamp in output, got: %s", result)
	}
}

func TestFormatEntry_ContainsPortInfo(t *testing.T) {
	f := makeFormatter()
	e := Entry{
		Timestamp:      time.Now(),
		ScanDurationMs: 10,
		Ports:          []portscanner.Port{{Port: 443, Protocol: "tcp"}, {Port: 22, Protocol: "tcp"}},
	}
	result := f.FormatEntry(e)
	if !strings.Contains(result, "443/tcp") {
		t.Errorf("expected port 443/tcp in output, got: %s", result)
	}
	if !strings.Contains(result, "22/tcp") {
		t.Errorf("expected port 22/tcp in output, got: %s", result)
	}
}

func TestFormatEntry_NoPorts_ShowsNone(t *testing.T) {
	f := makeFormatter()
	e := Entry{
		Timestamp:      time.Now(),
		ScanDurationMs: 5,
		Ports:          []portscanner.Port{},
	}
	result := f.FormatEntry(e)
	if !strings.Contains(result, "none") {
		t.Errorf("expected 'none' for empty ports, got: %s", result)
	}
}

func TestFormatAll_ReturnsOneLinePerEntry(t *testing.T) {
	f := makeFormatter()
	entries := []Entry{
		{Timestamp: time.Now(), ScanDurationMs: 1, Ports: []portscanner.Port{}},
		{Timestamp: time.Now(), ScanDurationMs: 2, Ports: []portscanner.Port{}},
		{Timestamp: time.Now(), ScanDurationMs: 3, Ports: []portscanner.Port{}},
	}
	results := f.FormatAll(entries)
	if len(results) != 3 {
		t.Errorf("expected 3 formatted lines, got %d", len(results))
	}
}

func TestFormatAll_EmptyEntries(t *testing.T) {
	f := makeFormatter()
	results := f.FormatAll([]Entry{})
	if len(results) != 0 {
		t.Errorf("expected empty result for empty input, got %d items", len(results))
	}
}
