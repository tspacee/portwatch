package metrics

import (
	"fmt"
	"io"
	"os"
	"text/tabwriter"
	"time"
)

// Reporter prints a human-readable summary of collected metrics.
type Reporter struct {
	collector *Collector
	out       io.Writer
}

// NewReporter creates a Reporter that writes to out.
// If out is nil, os.Stdout is used.
func NewReporter(c *Collector, out io.Writer) *Reporter {
	if out == nil {
		out = os.Stdout
	}
	return &Reporter{collector: c, out: out}
}

// Print writes the current metrics summary to the configured writer.
func (r *Reporter) Print() {
	s := r.collector.Summary()

	w := tabwriter.NewWriter(r.out, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "=== portwatch metrics ===")
	fmt.Fprintf(w, "Total scans:\t%d\n", s.TotalScans)
	fmt.Fprintf(w, "Errored scans:\t%d\n", s.ErroredScans)
	fmt.Fprintf(w, "Total violations:\t%d\n", s.TotalViolations)
	fAvg scan duration:\t%v\n", s.AvgScanDuration.Round(time.Millisecond))

	if !s.LastScan.IsZero() {
		fmt.Fprintf(w, "Last scan:\t%s\n", s.LastScan.Format(time.RFC3339))
		fmt.Fprintf(w, "Last open ports:\t%d\n", s.LastOpenPorts)
	} else {
		fmt.Fprintln(w, "Last scan:\t—")
	}

	w.Flush()
}

// PrintTo writes the current metrics summary to the provided writer,
// without changing the Reporter's configured output destination.
func (r *Reporter) PrintTo(out io.Writer) {
	tmp := r.out
	r.out = out
	r.Print()
	r.out = tmp
}
