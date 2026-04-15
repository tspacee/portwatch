package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Config holds the full portwatch configuration.
type Config struct {
	ScanInterval int          `yaml:"scan_interval_seconds"`
	PortRange    PortRange    `yaml:"port_range"`
	Protocols    []string     `yaml:"protocols"`
	Rules        []RuleConfig `yaml:"rules"`
	Alert        AlertConfig  `yaml:"alert"`
}

// PortRange defines the inclusive start/end ports to scan.
type PortRange struct {
	Start int `yaml:"start"`
	End   int `yaml:"end"`
}

// RuleConfig mirrors the rule definition in YAML.
type RuleConfig struct {
	Name     string `yaml:"name"`
	Port     int    `yaml:"port"`
	Protocol string `yaml:"protocol"`
	Action   string `yaml:"action"`
}

// AlertConfig holds alerting backend settings.
type AlertConfig struct {
	Output string `yaml:"output"` // "stdout" or a file path
}

// Load reads and parses a YAML config file from the given path.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("config: read file %q: %w", path, err)
	}

	cfg := &Config{}
	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("config: parse yaml: %w", err)
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

// Validate checks that the configuration values are sensible.
func (c *Config) Validate() error {
	if c.ScanInterval <= 0 {
		c.ScanInterval = 30
	}
	if c.PortRange.Start <= 0 {
		c.PortRange.Start = 1
	}
	if c.PortRange.End <= 0 {
		c.PortRange.End = 65535
	}
	if c.PortRange.Start > c.PortRange.End {
		return fmt.Errorf("config: port_range start (%d) must be <= end (%d)",
			c.PortRange.Start, c.PortRange.End)
	}
	if len(c.Protocols) == 0 {
		c.Protocols = []string{"tcp"}
	}
	return nil
}
