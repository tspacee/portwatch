package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Config holds the full portwatch configuration.
type Config struct {
	ScanInterval int         `yaml:"scan_interval"`
	PortRange    PortRange   `yaml:"port_range"`
	Protocols    []string    `yaml:"protocols"`
	SnapshotPath string      `yaml:"snapshot_path"`
	Rules        []RuleEntry `yaml:"rules"`
	Alerts       AlertConfig `yaml:"alerts"`
}

// PortRange defines the start and end ports to scan.
type PortRange struct {
	Start int `yaml:"start"`
	End   int `yaml:"end"`
}

// RuleEntry maps to a rule definition in YAML.
type RuleEntry struct {
	Name     string `yaml:"name"`
	Port     int    `yaml:"port"`
	Protocol string `yaml:"protocol"`
	Action   string `yaml:"action"`
}

// AlertConfig holds notifier configurations.
type AlertConfig struct {
	Log     bool        `yaml:"log"`
	Webhook string      `yaml:"webhook"`
	Slack   string      `yaml:"slack"`
	Email   EmailAlert  `yaml:"email"`
}

// EmailAlert holds email-specific alert configuration.
type EmailAlert struct {
	Enabled  bool     `yaml:"enabled"`
	SMTPHost string   `yaml:"smtp_host"`
	SMTPPort int      `yaml:"smtp_port"`
	Username string   `yaml:"username"`
	Password string   `yaml:"password"`
	From     string   `yaml:"from"`
	To       []string `yaml:"to"`
}

// Load reads and validates a config file at the given path.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("config: read file: %w", err)
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
	if cfg.PortRange.Start < 1 || cfg.PortRange.End > 65535 || cfg.PortRange.Start > cfg.PortRange.End {
		return fmt.Errorf("config: invalid port range %d-%d", cfg.PortRange.Start, cfg.PortRange.End)
	}
	if cfg.ScanInterval < 1 {
		return fmt.Errorf("config: scan_interval must be >= 1")
	}
	return nil
}
