package rules

import "fmt"

// Match represents the result of evaluating a rule against a port state.
type Match struct {
	Rule      Rule
	Port      int
	IsOpen    bool
	Violation bool // true if the actual state differs from expected
	Message   string
}

// Engine holds a set of rules and evaluates port states against them.
type Engine struct {
	rules []Rule
}

// NewEngine creates an Engine with the given rules, validating each one.
func NewEngine(rules []Rule) (*Engine, error) {
	for _, r := range rules {
		if err := r.Validate(); err != nil {
			return nil, fmt.Errorf("invalid rule: %w", err)
		}
	}
	return &Engine{rules: rules}, nil
}

// Evaluate checks the given set of open ports against all rules and returns matches.
func (e *Engine) Evaluate(openPorts map[int]bool) []Match {
	var matches []Match
	for _, rule := range e.rules {
		if rule.Action == ActionIgnore {
			continue
		}
		isOpen := openPorts[rule.Port]
		violation := isOpen != rule.Expected
		if !violation {
			continue
		}
		var msg string
		if isOpen && !rule.Expected {
			msg = fmt.Sprintf("[%s] port %d/%s is OPEN but expected CLOSED", rule.Severity, rule.Port, rule.Protocol)
		} else {
			msg = fmt.Sprintf("[%s] port %d/%s is CLOSED but expected OPEN", rule.Severity, rule.Port, rule.Protocol)
		}
		matches = append(matches, Match{
			Rule:      rule,
			Port:      rule.Port,
			IsOpen:    isOpen,
			Violation: true,
			Message:   msg,
		})
	}
	return matches
}

// Rules returns a copy of the engine's rules.
func (e *Engine) Rules() []Rule {
	copy := make([]Rule, len(e.rules))
	for i, r := range e.rules {
		copy[i] = r
	}
	return copy
}
