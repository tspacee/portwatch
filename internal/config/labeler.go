package config

import "errors"

// PortLabel maps a port number to a descriptive label.
type PortLabel struct {
	Port  int    `yaml:"port"`
	Label string `yaml:"label"`
}

// LabelerConfig holds port label configuration.
type LabelerConfig struct {
	Labels []PortLabel `yaml:"labels"`
}

// defaultLabelerConfig returns a LabelerConfig with common well-known ports.
func defaultLabelerConfig() LabelerConfig {
	return LabelerConfig{
		Labels: []PortLabel{
			{Port: 22, Label: "ssh"},
			{Port: 80, Label: "http"},
			{Port: 443, Label: "https"},
			{Port: 3306, Label: "mysql"},
			{Port: 5432, Label: "postgres"},
		},
	}
}

// Validate checks that all labels are non-empty and ports are in range.
func (c LabelerConfig) Validate() error {
	for _, pl := range c.Labels {
		if pl.Port < 1 || pl.Port > 65535 {
			return errors.New("labeler: port out of valid range")
		}
		if pl.Label == "" {
			return errors.New("labeler: label must not be empty")
		}
	}
	return nil
}
