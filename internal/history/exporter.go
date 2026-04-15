package history

import (
	"encoding/csv"
	"fmt"
	"io"
	"time"
)

// Exporter writes history entries to an external format.
type Exporter struct {
	h *History
}

// NewExporter creates an Exporter backed by the given History.
func NewExporter(h *History) (*Exporter, error) {
	if h == nil {
		return nil, fmt.Errorf("history must not be nil")
	}
	return &Exporter{h: h}, nil
}

// WriteCSV writes all history entries as CSV rows to w.
// Columns: timestamp, port, protocol, event, detail
func (e *Exporter) WriteCSV(w io.Writer) error {
	cw := csv.NewWriter(w)

	if err := cw.Write([]string{"timestamp", "port", "protocol", "event", "detail"}); err != nil {
		return fmt.Errorf("write csv header: %w", err)
	}

	for _, entry := range e.h.Entries() {
		row := []string{
			entry.Timestamp.UTC().Format(time.RFC3339),
			fmt.Sprintf("%d", entry.Port),
			entry.Protocol,
			entry.Event,
			entry.Detail,
		}
		if err := cw.Write(row); err != nil {
			return fmt.Errorf("write csv row: %w", err)
		}
	}

	cw.Flush()
	return cw.Error()
}
