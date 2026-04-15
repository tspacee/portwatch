package history

import (
	"fmt"
	"strings"
	"time"
)

// Formatter renders history entries as human-readable text.
type Formatter struct {
	timeFormat string
}

// NewFormatter returns a Formatter with a default time layout.
func NewFormatter(timeFormat string) *Formatter {
	if timeFormat == "" {
		timeFormat = time.RFC3339
	}
	return &Formatter{timeFormat: timeFormat}
}

// FormatEntry returns a single-line string representation of an Entry.
func (f *Formatter) FormatEntry(e Entry) string {
	ports := make([]string, len(e.Ports))
	for i, p := range e.Ports {
		ports[i] = fmt.Sprintf("%d/%s", p.Port, p.Protocol)
	}
	portList := strings.Join(ports, ", ")
	if portList == "" {
		portList = "none"
	}
	return fmt.Sprintf("[%s] scan=%d open=%d ports=[%s]",
		e.Timestamp.Format(f.timeFormat),
		e.ScanDurationMs,
		len(e.Ports),
		portList,
	)
}

// FormatAll returns a slice of formatted strings for every entry.
func (f *Formatter) FormatAll(entries []Entry) []string {
	out := make([]string, len(entries))
	for i, e := range entries {
		out[i] = f.FormatEntry(e)
	}
	return out
}
