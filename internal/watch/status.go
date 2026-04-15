package watch

import (
	"fmt"
	"io"
	"os"
	"text/tabwriter"
	"time"

	"github.com/user/portwatch/internal/metrics"
)

// Status holds a point-in-time view of the watcher's health.
type Status struct {
	Running   bool
	StartedAt time.Time
	Summary   metrics.Summary
}

// StatusReporter prints a human-readable status block.
type StatusReporter struct {
	w io.Writer
}

// NewStatusReporter creates a StatusReporter that writes to w.
// If w is nil, os.Stdout is used.
func NewStatusReporter(w io.Writer) *StatusReporter {
	if w == nil {
		w = os.Stdout
	}
	return &StatusReporter{w: w}
}

// Print writes the status to the underlying writer.
func (r *StatusReporter) Print(s Status) error {
	tw := tabwriter.NewWriter(r.w, 0, 0, 2, ' ', 0)
	running := "no"
	if s.Running {
		running = "yes"
	}
	lines := []struct{ k, v string }{
		{"running", running},
		{"started_at", s.StartedAt.Format(time.RFC3339)},
		{"total_scans", fmt.Sprintf("%d", s.Summary.TotalScans)},
		{"total_errors", fmt.Sprintf("%d", s.Summary.TotalErrors)},
		{"last_scan", formatTime(s.Summary.LastScan)},
	}
	for _, l := range lines {
		if _, err := fmt.Fprintf(tw, "%s\t%s\n", l.k, l.v); err != nil {
			return err
		}
	}
	return tw.Flush()
}

func formatTime(t time.Time) string {
	if t.IsZero() {
		return "-"
	}
	return t.Format(time.RFC3339)
}
