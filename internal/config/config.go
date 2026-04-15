package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Config holds the full portwatch configuration.
type Config struct {
	ScanIntervalSeconds int      `yaml:"scan_interval_seconds"`
	PortRange           PortRange `yaml:"port_range"`
	Protocols           []string  `yaml:"protocols"`
	SnapshotPath        string    `yaml:"snapshot_path"`
	WebhookURL          string    `yaml:"webhook_url"`
	Rules               []RuleConfig `yaml:"rules"`
}

// PortRange defines the inclusive range of ports to scan.
type PortRange struct {
	From int `yaml:"from"`
	To   int `yaml:"to"`
}

// RuleConfig mirrors rules.Rule for YAML unmarshalling.
type RuleConfig struct {
	Name     string `yaml:"name"`
	Port     int    `yaml:"port"`
	Protocol string `yaml:"protocol"`
	Action   string `yaml:"action"`
}

// Load reads and validates a YAML config file at path.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("config: read file %q: %w", path, err)
	}

	cfg := Default()
	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("config: parse yaml: %w", err)
	}

	if err := validate(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

func validate(cfg *Config) error {
	if cfg.PortRange.From < 1 || cfg.PortRange.To > 65535 {
		return fmt.Errorf("config: port range [%d-%d] out of valid bounds 1-65535",
			cfg.PortRange.From, cfg.PortRange.To)
	}
	if cfg.PortRange.From > cfg.PortRange.To {
		return fmt.Errorf("config: port_range.from (%d) must be <= port_range.to (%d)",
			cfg.PortRange.From, cfg.PortRange.To)
	}
	if cfg.ScanIntervalSeconds < 1 {
		return fmt.Errorf("config: scan_interval_seconds must be >= 1")
	}
	return nil
}
