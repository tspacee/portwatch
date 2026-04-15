package alert_test

import (
	"testing"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/config"
)

func baseConfig() *config.Config {
	cfg := config.Default()
	return cfg
}

func TestBuildNotifier_DefaultsToLogNotifier(t *testing.T) {
	cfg := baseConfig()
	cfg.Alerts.Log.Enabled = false

	n, err := alert.BuildNotifier(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n == nil {
		t.Fatal("expected non-nil notifier")
	}
}

func TestBuildNotifier_LogEnabled(t *testing.T) {
	cfg := baseConfig()
	cfg.Alerts.Log.Enabled = true

	n, err := alert.BuildNotifier(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n == nil {
		t.Fatal("expected non-nil notifier")
	}
}

func TestBuildNotifier_WebhookConfigured(t *testing.T) {
	cfg := baseConfig()
	cfg.Alerts.Webhook.URL = "http://example.com/hook"

	n, err := alert.BuildNotifier(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n == nil {
		t.Fatal("expected non-nil notifier")
	}
}

func TestBuildNotifier_SlackConfigured(t *testing.T) {
	cfg := baseConfig()
	cfg.Alerts.Slack.WebhookURL = "http://hooks.slack.com/test"

	n, err := alert.BuildNotifier(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n == nil {
		t.Fatal("expected non-nil notifier")
	}
}

func TestBuildNotifier_InvalidSlackURL(t *testing.T) {
	cfg := baseConfig()
	cfg.Alerts.Slack.WebhookURL = ""
	cfg.Alerts.Log.Enabled = false

	// No slack URL means slack notifier is skipped; should still succeed.
	n, err := alert.BuildNotifier(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n == nil {
		t.Fatal("expected fallback notifier")
	}
}

func TestBuildNotifier_MultipleNotifiers(t *testing.T) {
	cfg := baseConfig()
	cfg.Alerts.Log.Enabled = true
	cfg.Alerts.Webhook.URL = "http://example.com/hook"
	cfg.Alerts.Slack.WebhookURL = "http://hooks.slack.com/test"

	n, err := alert.BuildNotifier(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n == nil {
		t.Fatal("expected non-nil multi-notifier")
	}
}
