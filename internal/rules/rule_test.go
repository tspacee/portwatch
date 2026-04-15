package rules

import (
	"testing"
)

func validRule() Rule {
	return Rule{
		Name:     "test-rule",
		Port:     8080,
		Protocol: "tcp",
		Expected: true,
		Action:   ActionAlert,
		Severity: "warning",
	}
}

func TestRule_Validate_Valid(t *testing.T) {
	r := validRule()
	if err := r.Validate(); err != nil {
		t.Errorf("expected no error, got: %v", err)
	}
}

func TestRule_Validate_EmptyName(t *testing.T) {
	r := validRule()
	r.Name = ""
	if err := r.Validate(); err == nil {
		t.Error("expected error for empty name")
	}
}

func TestRule_Validate_InvalidPort(t *testing.T) {
	for _, port := range []int{0, -1, 65536, 99999} {
		r := validRule()
		r.Port = port
		if err := r.Validate(); err == nil {
			t.Errorf("expected error for port %d", port)
		}
	}
}

func TestRule_Validate_InvalidProtocol(t *testing.T) {
	r := validRule()
	r.Protocol = "http"
	if err := r.Validate(); err == nil {
		t.Error("expected error for invalid protocol")
	}
}

func TestRule_Validate_InvalidAction(t *testing.T) {
	r := validRule()
	r.Action = "notify"
	if err := r.Validate(); err == nil {
		t.Error("expected error for invalid action")
	}
}

func TestRule_Validate_InvalidSeverity(t *testing.T) {
	r := validRule()
	r.Severity = "high"
	if err := r.Validate(); err == nil {
		t.Error("expected error for invalid severity")
	}
}
