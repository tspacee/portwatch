package alert_test

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/rules"
)

func sampleViolation() rules.Violation {
	return rules.Violation{
		RuleName: "block-8080",
		Port:     8080,
		Protocol: "tcp",
	}
}

func TestViolationToAlert_Fields(t *testing.T) {
	v := sampleViolation()
	before := time.Now()
	a := alert.ViolationToAlert(v)
	after := time.Now()

	if a.Level != alert.LevelAlert {
		t.Errorf("expected level ALERT, got %s", a.Level)
	}
	if a.Violation != v {
		t.Errorf("violation mismatch")
	}
	if a.Timestamp.Before(before) || a.Timestamp.After(after) {
		t.Errorf("timestamp out of expected range")
	}
	if !strings.Contains(a.Message, "8080") {
		t.Errorf("message should contain port number, got: %s", a.Message)
	}
}

func TestLogNotifier_Notify_Output(t *testing.T) {
	var buf bytes.Buffer
	n := alert.NewLogNotifier(&buf)
	v := sampleViolation()
	a := alert.ViolationToAlert(v)

	if err := n.Notify(a); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	for _, want := range []string{"ALERT", "8080", "tcp", "block-8080"} {
		if !strings.Contains(out, want) {
			t.Errorf("output missing %q, got: %s", want, out)
		}
	}
}

func TestNewLogNotifier_NilWriterUsesStdout(t *testing.T) {
	n := alert.NewLogNotifier(nil)
	if n == nil {
		t.Fatal("expected non-nil notifier")
	}
	if n.Out == nil {
		t.Error("expected Out to be set to stdout")
	}
}
