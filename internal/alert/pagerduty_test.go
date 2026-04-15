package alert

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNewPagerDutyNotifier_EmptyKey(t *testing.T) {
	_, err := NewPagerDutyNotifier("", 0)
	if err == nil {
		t.Fatal("expected error for empty integration key")
	}
}

func TestNewPagerDutyNotifier_DefaultTimeout(t *testing.T) {
	n, err := NewPagerDutyNotifier("test-key", 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n.client.Timeout != 10*time.Second {
		t.Errorf("expected default timeout 10s, got %v", n.client.Timeout)
	}
}

func TestPagerDutyNotifier_Notify_Success(t *testing.T) {
	var received pagerDutyPayload
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Errorf("decode body: %v", err)
		}
		w.WriteHeader(http.StatusAccepted)
	}))
	defer ts.Close()

	n, _ := NewPagerDutyNotifier("key-abc", 5*time.Second)
	n.eventsURL = ts.URL

	a := Alert{Message: "port 22 opened", Severity: "critical"}
	if err := n.Notify(a); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if received.RoutingKey != "key-abc" {
		t.Errorf("expected routing key 'key-abc', got %q", received.RoutingKey)
	}
	if received.EventAction != "trigger" {
		t.Errorf("expected event_action 'trigger', got %q", received.EventAction)
	}
	if received.Payload.Severity != "critical" {
		t.Errorf("expected severity 'critical', got %q", received.Payload.Severity)
	}
	if received.Payload.Source != "portwatch" {
		t.Errorf("expected source 'portwatch', got %q", received.Payload.Source)
	}
}

func TestPagerDutyNotifier_Notify_NonOKStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
	}))
	defer ts.Close()

	n, _ := NewPagerDutyNotifier("key-abc", 5*time.Second)
	n.eventsURL = ts.URL

	err := n.Notify(Alert{Message: "test", Severity: "warning"})
	if err == nil {
		t.Fatal("expected error for non-2xx status")
	}
}

func TestPagerDutyNotifier_Notify_UnreachableURL(t *testing.T) {
	n, _ := NewPagerDutyNotifier("key-abc", time.Second)
	n.eventsURL = "http://127.0.0.1:1"

	err := n.Notify(Alert{Message: "test", Severity: "info"})
	if err == nil {
		t.Fatal("expected error for unreachable URL")
	}
}

func TestSeverityToPagerDuty(t *testing.T) {
	cases := []struct{ in, want string }{
		{"critical", "critical"},
		{"warning", "warning"},
		{"info", "info"},
		{"unknown", "info"},
	}
	for _, c := range cases {
		if got := severityToPagerDuty(c.in); got != c.want {
			t.Errorf("severityToPagerDuty(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}
