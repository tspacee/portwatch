package metrics_test

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/user/portwatch/internal/metrics"
)

func TestReporter_Print_ContainsExpectedFields(t *testing.T) {
	c := metrics.NewCollector()
	c.Record(metrics.ScanResult{
		Timestamp:    time.Now(),
		OpenPorts:    4,
		Violations:   1,
		ScanDuration: 15 * time.Millisecond,
	})

	var buf bytes.Buffer
	r := metrics.NewReporter(c, &buf)
	r.Print()

	out := buf.String()
	for _, want := range []string{
		"Total scans",
		"Errored scans",
		"Total violations",
		"Avg scan duration",
		"Last scan",
		"Last open ports",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("expected output to contain %q, got:\n%s", want, out)
		}
	}
}

func TestReporter_Print_NoScans_ShowsDash(t *testing.T) {
	c := metrics.NewCollector()

	var buf bytes.Buffer
	r := metrics.NewReporter(c, &buf)
	r.Print()

	if !strings.Contains(buf.String(), "—") {
		t.Error("expected dash placeholder when no scans recorded")
	}
}

func TestNewReporter_NilWriter_UsesStdout(t *testing.T) {
	// Should not panic when out is nil.
	c := metrics.NewCollector()
	r := metrics.NewReporter(c, nil)
	if r == nil {
		t.Fatal("expected non-nil reporter")
	}
}
