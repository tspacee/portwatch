package rules

import "fmt"

// Action defines what to do when a rule matches.
type Action string

const (
	ActionAlert  Action = "alert"
	ActionIgnore Action = "ignore"
)

// Rule defines a single port monitoring rule.
type Rule struct {
	Name      string `yaml:"name"`
	Port      int    `yaml:"port"`
	Protocol  string `yaml:"protocol"` // tcp or udp
	Expected  bool   `yaml:"expected"`  // true if port should be open
	Action    Action `yaml:"action"`
	Severity  string `yaml:"severity"` // info, warning, critical
}

// Validate checks that the rule fields are valid.
func (r *Rule) Validate() error {
	if r.Name == "" {
		return fmt.Errorf("rule name must not be empty")
	}
	if r.Port < 1 || r.Port > 65535 {
		return fmt.Errorf("rule %q: port %d is out of valid range (1-65535)", r.Name, r.Port)
	}
	if r.Protocol != "tcp" && r.Protocol != "udp" {
		return fmt.Errorf("rule %q: protocol must be 'tcp' or 'udp', got %q", r.Name, r.Protocol)
	}
	if r.Action != ActionAlert && r.Action != ActionIgnore {
		return fmt.Errorf("rule %q: action must be 'alert' or 'ignore', got %q", r.Name, r.Action)
	}
	validSeverities := map[string]bool{"info": true, "warning": true, "critical": true}
	if !validSeverities[r.Severity] {
		return fmt.Errorf("rule %q: severity must be info, warning, or critical, got %q", r.Name, r.Severity)
	}
	return nil
}
