package alert

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// WebhookNotifier sends alert payloads to an HTTP endpoint.
type WebhookNotifier struct {
	url    string
	client *http.Client
}

// WebhookPayload is the JSON body sent to the webhook endpoint.
type WebhookPayload struct {
	Timestamp string `json:"timestamp"`
	Port      int    `json:"port"`
	Protocol  string `json:"protocol"`
	Action    string `json:"action"`
	Rule      string `json:"rule"`
	Message   string `json:"message"`
}

// NewWebhookNotifier creates a WebhookNotifier that posts to the given URL.
// A zero timeout defaults to 5 seconds.
func NewWebhookNotifier(url string, timeout time.Duration) *WebhookNotifier {
	if timeout == 0 {
		timeout = 5 * time.Second
	}
	return &WebhookNotifier{
		url:    url,
		client: &http.Client{Timeout: timeout},
	}
}

// Notify serialises the Alert as JSON and POSTs it to the configured URL.
func (w *WebhookNotifier) Notify(a Alert) error {
	payload := WebhookPayload{
		Timestamp: a.Timestamp.UTC().Format(time.RFC3339),
		Port:      a.Port,
		Protocol:  a.Protocol,
		Action:    a.Action,
		Rule:      a.Rule,
		Message:   a.Message,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("webhook: marshal payload: %w", err)
	}

	resp, err := w.client.Post(w.url, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("webhook: post to %s: %w", w.url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("webhook: unexpected status %d from %s", resp.StatusCode, w.url)
	}
	return nil
}
