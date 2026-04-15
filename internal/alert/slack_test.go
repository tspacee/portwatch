package alert

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestNewSlackNotifier_EmptyURL(t *testing.T) {
	_, err := NewSlackNotifier("", 0)
	if err == nil {
		t.Fatal("expected error for empty webhook URL")
	}
}

func TestNewSlackNotifier_DefaultTimeout(t *testing.T) {
	n, err := NewSlackNotifier("https://hooks.slack.com/test", 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n.client.Timeout != 10*time.Second {
		t.Errorf("expected default timeout 10s, got %v", n.client.Timeout)
	}
}

func TestSlackNotifier_Notify_Success(t *testing.T) {
	var received map[string]string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Errorf("decode body: %v", err)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	n, _ := NewSlackNotifier(ts.URL, time.Second)
	err := n.Notify(sampleViolation())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(received["text"], "portwatch") {
		t.Errorf("expected 'portwatch' in slack text, got: %s", received["text"])
	}
}

func TestSlackNotifier_Notify_NonOKStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer ts.Close()

	n, _ := NewSlackNotifier(ts.URL, time.Second)
	err := n.Notify(sampleViolation())
	if err == nil {
		t.Fatal("expected error for non-OK status")
	}
	if !strings.Contains(err.Error(), "403") {
		t.Errorf("expected 403 in error, got: %v", err)
	}
}

func TestSlackNotifier_Notify_UnreachableURL(t *testing.T) {
	n, _ := NewSlackNotifier("http://127.0.0.1:1", 500*time.Millisecond)
	err := n.Notify(sampleViolation())
	if err == nil {
		t.Fatal("expected error for unreachable URL")
	}
}
