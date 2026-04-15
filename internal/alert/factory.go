package alert

import (
	"fmt"

	"github.com/user/portwatch/internal/config"
)

// BuildNotifier constructs a Notifier (or MultiNotifier) from the given config.
func BuildNotifier(cfg *config.Config) (Notifier, error) {
	var notifiers []Notifier

	if cfg.Alerts.Log.Enabled {
		notifiers = append(notifiers, NewLogNotifier(nil))
	}

	if cfg.Alerts.Webhook.URL != "" {
		n, err := NewWebhookNotifier(cfg.Alerts.Webhook.URL, cfg.Alerts.Webhook.Timeout)
		if err != nil {
			return nil, fmt.Errorf("webhook notifier: %w", err)
		}
		notifiers = append(notifiers, n)
	}

	if cfg.Alerts.Slack.WebhookURL != "" {
		n, err := NewSlackNotifier(cfg.Alerts.Slack.WebhookURL, cfg.Alerts.Slack.Timeout)
		if err != nil {
			return nil, fmt.Errorf("slack notifier: %w", err)
		}
		notifiers = append(notifiers, n)
	}

	if cfg.Alerts.PagerDuty.IntegrationKey != "" {
		n, err := NewPagerDutyNotifier(cfg.Alerts.PagerDuty.IntegrationKey, cfg.Alerts.PagerDuty.Timeout)
		if err != nil {
			return nil, fmt.Errorf("pagerduty notifier: %w", err)
		}
		notifiers = append(notifiers, n)
	}

	if cfg.Alerts.Email.SMTPHost != "" {
		n, err := NewEmailNotifier(
			cfg.Alerts.Email.SMTPHost,
			cfg.Alerts.Email.SMTPPort,
			cfg.Alerts.Email.From,
			cfg.Alerts.Email.Recipients,
		)
		if err != nil {
			return nil, fmt.Errorf("email notifier: %w", err)
		}
		notifiers = append(notifiers, n)
	}

	if len(notifiers) == 0 {
		return NewLogNotifier(nil), nil
	}

	return NewMultiNotifier(notifiers)
}
