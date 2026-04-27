package watch

import (
	"fmt"
	"io"
	"os"
)

// PulseReporter formats and writes pulse summary statistics to an io.Writer.
type PulseReporter struct {
	pulse  *Pulse
	writer io.Writer
}

// NewPulseReporter returns a PulseReporter backed by the given Pulse.
// If writer is nil, os.Stdout is used.
func NewPulseReporter(p *Pulse, w io.Writer) *PulseReporter {
	if w == nil {
		w = os.Stdout
	}
	return &PulseReporter{pulse: p, writer: w}
}

// Print writes a human-readable summary of the current pulse statistics.
// If no intervals have been recorded, a placeholder line is printed.
func (r *PulseReporter) Print() {
	s, err := r.pulse.Summary()
	if err != nil {
		fmt.Fprintln(r.writer, "pulse: no data")
		return
	}
	fmt.Fprintf(r.writer,
		"pulse: count=%d avg=%v min=%v max=%v\n",
		s.Count, s.Avg.Round(1), s.Min.Round(1), s.Max.Round(1),
	)
}
