package alert

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNewWebhookNotifier_DefaultTimeout(t *testing.T) {
	w := NewWebhookNotifier("http://example.com", 0)
	if w.client.Timeout != 5*time.Second {
		t.Errorf("expected default timeout 5s, got %v", w.client.Timeout)
	}
}

func TestWebhookNotifier_Notify_Success(t *testing.T) {
	var received WebhookPayload

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if ct := r.Header.Get("Content-Type"); ct != "application/json" {
			t.Errorf("expected application/json, got %s", ct)
		}
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Errorf("decode body: %v", err)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	notifier := NewWebhookNotifier(server.URL, time.Second)
	alert := sampleViolation()
	a := ViolationToAlert(alert)

	if err := notifier.Notify(a); err != nil {
		t.Fatalf("Notify returned error: %v", err)
	}

	if received.Port != a.Port {
		t.Errorf("port: want %d, got %d", a.Port, received.Port)
	}
	if received.Rule != a.Rule {
		t.Errorf("rule: want %s, got %s", a.Rule, received.Rule)
	}
	if received.Action != a.Action {
		t.Errorf("action: want %s, got %s", a.Action, received.Action)
	}
}

func TestWebhookNotifier_Notify_NonOKStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	notifier := NewWebhookNotifier(server.URL, time.Second)
	a := ViolationToAlert(sampleViolation())

	if err := notifier.Notify(a); err == nil {
		t.Error("expected error for non-2xx status, got nil")
	}
}

func TestWebhookNotifier_Notify_UnreachableURL(t *testing.T) {
	notifier := NewWebhookNotifier("http://127.0.0.1:1", 200*time.Millisecond)
	a := ViolationToAlert(sampleViolation())

	if err := notifier.Notify(a); err == nil {
		t.Error("expected error for unreachable URL, got nil")
	}
}
