package history

import (
	"net/http"
	"strconv"
	"time"

	"github.com/user/portwatch/internal/metrics"
)

// Middleware wraps an http.Handler to record request metrics into a Collector.
type Middleware struct {
	collector *metrics.Collector
	next      http.Handler
}

// NewMiddleware returns a new Middleware that records scan/request metrics.
// It returns an error if collector or next are nil.
func NewMiddleware(collector *metrics.Collector, next http.Handler) (*Middleware, error) {
	if collector == nil {
		return nil, ErrNilHistory
	}
	if next == nil {
		return nil, ErrNilHistory
	}
	return &Middleware{collector: collector, next: next}, nil
}

// ServeHTTP records the duration of each request and delegates to the
// wrapped handler. The duration is recorded as a scan cycle in the collector
// so that the metrics reporter can surface history-endpoint latency.
func (m *Middleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	rw := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
	m.next.ServeHTTP(rw, r)
	duration := time.Since(start)

	// Record the request as a lightweight scan cycle entry.
	m.collector.Record(0, duration)
	if rw.status >= http.StatusInternalServerError {
		m.collector.RecordError()
	}
}

// statusRecorder wraps http.ResponseWriter to capture the response status code.
type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (sr *statusRecorder) WriteHeader(code int) {
	sr.status = code
	sr.ResponseWriter.WriteHeader(code)
}

// parseLimit extracts and validates the "limit" query parameter.
// Returns defaultVal if the parameter is absent or invalid.
func parseLimit(r *http.Request, defaultVal int) int {
	raw := r.URL.Query().Get("limit")
	if raw == "" {
		return defaultVal
	}
	v, err := strconv.Atoi(raw)
	if err != nil || v <= 0 {
		return defaultVal
	}
	return v
}
