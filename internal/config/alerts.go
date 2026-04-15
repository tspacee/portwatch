package config

import "time"

// AlertsConfig holds all alert channel configurations.
type AlertsConfig struct {
	Log       LogAlertConfig       `yaml:"log"`
	Webhook   WebhookAlertConfig   `yaml:"webhook"`
	Slack     SlackAlertConfig     `yaml:"slack"`
	PagerDuty PagerDutyAlertConfig `yaml:"pagerduty"`
	Email     EmailAlertConfig     `yaml:"email"`
}

// LogAlertConfig configures the stdout/log notifier.
type LogAlertConfig struct {
	Enabled bool `yaml:"enabled"`
}

// WebhookAlertConfig configures a generic HTTP webhook notifier.
type WebhookAlertConfig struct {
	URL     string        `yaml:"url"`
	Timeout time.Duration `yaml:"timeout"`
}

// SlackAlertConfig configures the Slack incoming webhook notifier.
type SlackAlertConfig struct {
	WebhookURL string        `yaml:"webhook_url"`
	Timeout    time.Duration `yaml:"timeout"`
}

// PagerDutyAlertConfig configures the PagerDuty Events API notifier.
type PagerDutyAlertConfig struct {
	IntegrationKey string        `yaml:"integration_key"`
	Timeout        time.Duration `yaml:"timeout"`
}

// EmailAlertConfig configures the SMTP email notifier.
type EmailAlertConfig struct {
	SMTPHost   string   `yaml:"smtp_host"`
	SMTPPort   int      `yaml:"smtp_port"`
	From       string   `yaml:"from"`
	Recipients []string `yaml:"recipients"`
}

// defaultAlertsConfig returns sensible defaults for all alert channels.
func defaultAlertsConfig() AlertsConfig {
	return AlertsConfig{
		Log: LogAlertConfig{
			Enabled: true,
		},
		Webhook: WebhookAlertConfig{
			Timeout: 10 * time.Second,
		},
		Slack: SlackAlertConfig{
			Timeout: 10 * time.Second,
		},
		PagerDuty: PagerDutyAlertConfig{
			Timeout: 10 * time.Second,
		},
		Email: EmailAlertConfig{
			SMTPPort: 587,
		},
	}
}
