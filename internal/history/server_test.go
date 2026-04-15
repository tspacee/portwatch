package history

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNewServer_NilHistory(t *testing.T) {
	_, err := NewServer(nil)
	if err == nil {
		t.Fatal("expected error for nil history")
	}
}

func TestNewServer_Valid(t *testing.T) {
	h := New(10)
	s, err := NewServer(h)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.Handler() == nil {
		t.Fatal("expected non-nil handler")
	}
}

func TestServer_HandleHistory_Empty(t *testing.T) {
	h := New(10)
	s, _ := NewServer(h)

	req := httptest.NewRequest(http.MethodGet, "/history", nil)
	rec := httptest.NewRecorder()
	s.Handler().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	var entries []Entry
	if err := json.NewDecoder(rec.Body).Decode(&entries); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if len(entries) != 0 {
		t.Fatalf("expected 0 entries, got %d", len(entries))
	}
}

func TestServer_HandleHistory_WithEntries(t *testing.T) {
	h := New(10)
	for i := 0; i < 5; i++ {
		h.Add(Entry{Timestamp: time.Now(), Added: []int{8080 + i}})
	}
	s, _ := NewServer(h)

	req := httptest.NewRequest(http.MethodGet, "/history", nil)
	rec := httptest.NewRecorder()
	s.Handler().ServeHTTP(rec, req)

	var entries []Entry
	json.NewDecoder(rec.Body).Decode(&entries)
	if len(entries) != 5 {
		t.Fatalf("expected 5 entries, got %d", len(entries))
	}
}

func TestServer_HandleHistory_LimitParam(t *testing.T) {
	h := New(20)
	for i := 0; i < 10; i++ {
		h.Add(Entry{Timestamp: time.Now(), Added: []int{9000 + i}})
	}
	s, _ := NewServer(h)

	req := httptest.NewRequest(http.MethodGet, "/history?limit=3", nil)
	rec := httptest.NewRecorder()
	s.Handler().ServeHTTP(rec, req)

	var entries []Entry
	json.NewDecoder(rec.Body).Decode(&entries)
	if len(entries) != 3 {
		t.Fatalf("expected 3 entries with limit, got %d", len(entries))
	}
}
