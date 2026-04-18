package watch

import (
	"fmt"
	"io"
	"os"
	"text/tabwriter"
	"time"
)

// AuditWriter formats and writes audit log entries to an io.Writer.
type AuditWriter struct {
	w io.Writer
}

// NewAuditWriter creates an AuditWriter. If w is nil, os.Stdout is used.
func NewAuditWriter(w io.Writer) *AuditWriter {
	if w == nil {
		w = os.Stdout
	}
	return &AuditWriter{w: w}
}

// Write formats and outputs all entries from the given AuditLog.
func (aw *AuditWriter) Write(log *AuditLog) error {
	if log == nil {
		return fmt.Errorf("audit: log must not be nil")
	}
	entries := log.Entries()
	tw := tabwriter.NewWriter(aw.w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "TIMESTAMP\tPORT\tPROTOCOL\tEVENT\tSOURCE")
	for _, e := range entries {
		fmt.Fprintf(tw, "%s\t%d\t%s\t%s\t%s\n",
			e.Timestamp.Format(time.RFC3339),
			e.Port,
			e.Protocol,
			e.Event,
			e.Source,
		)
	}
	return tw.Flush()
}
