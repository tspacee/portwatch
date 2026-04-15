package alert

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// SlackNotifier sends alert notifications to a Slack webhook URL.
type SlackNotifier struct {
	webhookURL string
	client     *http.Client
}

type slackPayload struct {
	Text string `json:"text"`
}

// NewSlackNotifier creates a SlackNotifier with the given webhook URL.
// An optional timeout may be provided; if zero, a default of 10 seconds is used.
func NewSlackNotifier(webhookURL string, timeout time.Duration) (*SlackNotifier, error) {
	if webhookURL == "" {
		return nil, fmt.Errorf("slack webhook URL must not be empty")
	}
	if timeout <= 0 {
		timeout = 10 * time.Second
	}
	return &SlackNotifier{
		webhookURL: webhookURL,
		client:     &http.Client{Timeout: timeout},
	}, nil
}

// Notify sends the alert as a formatted Slack message.
func (s *SlackNotifier) Notify(a Alert) error {
	msg := fmt.Sprintf("*[portwatch] %s* — port %d/%s on %s\n%s",
		a.Severity, a.Port, a.Protocol, a.Host, a.Message)

	payload := slackPayload{Text: msg}
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("slack: marshal payload: %w", err)
	}

	resp, err := s.client.Post(s.webhookURL, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("slack: send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("slack: unexpected status code %d", resp.StatusCode)
	}
	return nil
}
