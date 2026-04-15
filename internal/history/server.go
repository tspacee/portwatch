package history

import (
	"encoding/json"
	"net/http"
	"strconv"
)

// Server exposes history entries over HTTP as JSON.
type Server struct {
	history *History
	mux     *http.ServeMux
}

// NewServer creates a new HTTP server backed by the given History.
func NewServer(h *History) (*Server, error) {
	if h == nil {
		return nil, ErrNilHistory
	}
	s := &Server{
		history: h,
		mux:     http.NewServeMux(),
	}
	s.mux.HandleFunc("/history", s.handleHistory)
	return s, nil
}

// Handler returns the underlying http.Handler.
func (s *Server) Handler() http.Handler {
	return s.mux
}

// handleHistory writes history entries as JSON. Accepts optional ?limit= query param.
func (s *Server) handleHistory(w http.ResponseWriter, r *http.Request) {
	entries := s.history.Entries()

	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if n, err := strconv.Atoi(limitStr); err == nil && n > 0 && n < len(entries) {
			entries = entries[len(entries)-n:]
		}
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(entries); err != nil {
		http.Error(w, "failed to encode history", http.StatusInternalServerError)
	}
}
