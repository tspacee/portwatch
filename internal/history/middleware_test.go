package history

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/user/portwatch/internal/metrics"
)

func TestNewMiddleware_NilCollector(t *testing.T) {
	_, err := NewMiddleware(nil, http.NotFoundHandler())
	if err == nil {
		t.Fatal("expected error for nil collector, got nil")
	}
}

func TestNewMiddleware_NilHandler(t *testing.T) {
	c := metrics.NewCollector()
	_, err := NewMiddleware(c, nil)
	if err == nil {
		t.Fatal("expected error for nil handler, got nil")
	}
}

func TestNewMiddleware_Valid(t *testing.T) {
	c := metrics.NewCollector()
	m, err := NewMiddleware(c, http.NotFoundHandler())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m == nil {
		t.Fatal("expected non-nil middleware")
	}
}

func TestMiddleware_RecordsMetrics(t *testing.T) {
	c := metrics.NewCollector()
	inner := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	m, err := NewMiddleware(c, inner)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/history", nil)
	rec := httptest.NewRecorder()
	m.ServeHTTP(rec, req)

	summary := c.Summary()
	if summary.TotalScans != 1 {
		t.Errorf("expected 1 scan recorded, got %d", summary.TotalScans)
	}
	if summary.ErrorCount != 0 {
		t.Errorf("expected 0 errors, got %d", summary.ErrorCount)
	}
}

func TestMiddleware_RecordsError_On5xx(t *testing.T) {
	c := metrics.NewCollector()
	inner := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})
	m, _ := NewMiddleware(c, inner)

	req := httptest.NewRequest(http.MethodGet, "/history", nil)
	rec := httptest.NewRecorder()
	m.ServeHTTP(rec, req)

	summary := c.Summary()
	if summary.ErrorCount != 1 {
		t.Errorf("expected 1 error recorded, got %d", summary.ErrorCount)
	}
}

func TestParseLimit_Default(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/history", nil)
	if got := parseLimit(req, 50); got != 50 {
		t.Errorf("expected default 50, got %d", got)
	}
}

func TestParseLimit_Valid(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/history?limit=10", nil)
	if got := parseLimit(req, 50); got != 10 {
		t.Errorf("expected 10, got %d", got)
	}
}

func TestParseLimit_Invalid(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/history?limit=abc", nil)
	if got := parseLimit(req, 25); got != 25 {
		t.Errorf("expected default 25, got %d", got)
	}
}
