package watch_test

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/user/portwatch/internal/metrics"
	"github.com/user/portwatch/internal/watch"
)

func TestNewStatusReporter_NilWriter(t *testing.T) {
	r := watch.NewStatusReporter(nil)
	if r == nil {
		t.Fatal("expected non-nil reporter")
	}
}

func TestPrint_RunningStatus(t *testing.T) {
	var buf bytes.Buffer
	r := watch.NewStatusReporter(&buf)
	s := watch.Status{
		Running:   true,
		StartedAt: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		Summary: metrics.Summary{
			TotalScans:  5,
			TotalErrors: 1,
		},
	}
	if err := r.Print(s); err != nil {
		t.Fatalf("Print: %v", err)
	}
	out := buf.String()
	for _, want := range []string{"running", "yes", "total_scans", "5", "total_errors", "1"} {
		if !strings.Contains(out, want) {
			t.Errorf("output missing %q\ngot: %s", want, out)
		}
	}
}

func TestPrint_NotRunning(t *testing.T) {
	var buf bytes.Buffer
	r := watch.NewStatusReporter(&buf)
	if err := r.Print(watch.Status{Running: false}); err != nil {
		t.Fatalf("Print: %v", err)
	}
	if !strings.Contains(buf.String(), "no") {
		t.Error("expected 'no' in output for non-running status")
	}
}

func TestPrint_LastScanZero_ShowsDash(t *testing.T) {
	var buf bytes.Buffer
	r := watch.NewStatusReporter(&buf)
	r.Print(watch.Status{}) //nolint:errcheck
	if !strings.Contains(buf.String(), "-") {
		t.Error("expected '-' when LastScan is zero")
	}
}

func TestPrint_LastScanSet(t *testing.T) {
	var buf bytes.Buffer
	r := watch.NewStatusReporter(&buf)
	now := time.Now().UTC().Truncate(time.Second)
	s := watch.Status{
		Summary: metrics.Summary{LastScan: now},
	}
	r.Print(s) //nolint:errcheck
	if !strings.Contains(buf.String(), now.Format("2006")) {
		t.Error("expected year in last_scan output")
	}
}
