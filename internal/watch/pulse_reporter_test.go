package watch

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func TestPulseReporter_Print_NoData(t *testing.T) {
	p := NewPulse()
	var buf bytes.Buffer
	r := NewPulseReporter(p, &buf)
	r.Print()

	if !strings.Contains(buf.String(), "no data") {
		t.Errorf("expected 'no data', got %q", buf.String())
	}
}

func TestPulseReporter_Print_WithData(t *testing.T) {
	p := NewPulse()
	p.Beat()
	time.Sleep(15 * time.Millisecond)
	p.Beat()

	var buf bytes.Buffer
	r := NewPulseReporter(p, &buf)
	r.Print()

	out := buf.String()
	for _, field := range []string{"count=", "avg=", "min=", "max="} {
		if !strings.Contains(out, field) {
			t.Errorf("expected field %q in output %q", field, out)
		}
	}
}

func TestNewPulseReporter_NilWriter_UsesStdout(t *testing.T) {
	p := NewPulse()
	r := NewPulseReporter(p, nil)
	if r.writer == nil {
		t.Fatal("expected non-nil writer when nil passed")
	}
}
