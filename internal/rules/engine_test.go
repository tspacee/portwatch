package rules

import (
	"testing"
)

func TestNewEngine_InvalidRule(t *testing.T) {
	badRule := Rule{Name: "", Port: 80, Protocol: "tcp", Action: ActionAlert, Severity: "info"}
	_, err := NewEngine([]Rule{badRule})
	if err == nil {
		t.Error("expected error for invalid rule, got nil")
	}
}

func TestEngine_Evaluate_NoViolations(t *testing.T) {
	rules := []Rule{
		{Name: "http-open", Port: 80, Protocol: "tcp", Expected: true, Action: ActionAlert, Severity: "critical"},
	}
	engine, err := NewEngine(rules)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	openPorts := map[int]bool{80: true}
	matches := engine.Evaluate(openPorts)
	if len(matches) != 0 {
		t.Errorf("expected 0 matches, got %d", len(matches))
	}
}

func TestEngine_Evaluate_UnexpectedOpenPort(t *testing.T) {
	rules := []Rule{
		{Name: "telnet-closed", Port: 23, Protocol: "tcp", Expected: false, Action: ActionAlert, Severity: "critical"},
	}
	engine, _ := NewEngine(rules)
	openPorts := map[int]bool{23: true}
	matches := engine.Evaluate(openPorts)
	if len(matches) != 1 {
		t.Fatalf("expected 1 match, got %d", len(matches))
	}
	if !matches[0].Violation {
		t.Error("expected violation to be true")
	}
}

func TestEngine_Evaluate_IgnoreAction(t *testing.T) {
	rules := []Rule{
		{Name: "ignored-port", Port: 9999, Protocol: "tcp", Expected: false, Action: ActionIgnore, Severity: "info"},
	}
	engine, _ := NewEngine(rules)
	openPorts := map[int]bool{9999: true}
	matches := engine.Evaluate(openPorts)
	if len(matches) != 0 {
		t.Errorf("expected 0 matches for ignored rule, got %d", len(matches))
	}
}

func TestEngine_Rules_ReturnsCopy(t *testing.T) {
	rules := []Rule{
		{Name: "ssh", Port: 22, Protocol: "tcp", Expected: true, Action: ActionAlert, Severity: "warning"},
	}
	engine, _ := NewEngine(rules)
	copy := engine.Rules()
	copy[0].Name = "modified"
	if engine.Rules()[0].Name != "ssh" {
		t.Error("expected Rules() to return a copy, not a reference")
	}
}
