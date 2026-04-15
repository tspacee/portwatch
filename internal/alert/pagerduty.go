package alert

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const pagerdutyEventsURL = "https://events.pagerduty.com/v2/enqueue"

// PagerDutyNotifier sends alerts to PagerDuty via the Events API v2.
type PagerDutyNotifier struct {
	integrationKey string
	client         *http.Client
	eventsURL      string
}

type pagerDutyPayload struct {
	RoutingKey  string            `json:"routing_key"`
	EventAction string            `json:"event_action"`
	Payload     pagerDutyDetail   `json:"payload"`
	Links       []pagerDutyLink   `json:"links,omitempty"`
}

type pagerDutyDetail struct {
	Summary  string `json:"summary"`
	Severity string `json:"severity"`
	Source   string `json:"source"`
}

type pagerDutyLink struct {
	Href string `json:"href"`
	Text string `json:"text"`
}

// NewPagerDutyNotifier creates a new PagerDutyNotifier.
// integrationKey is the PagerDuty Events API v2 integration key.
func NewPagerDutyNotifier(integrationKey string, timeout time.Duration) (*PagerDutyNotifier, error) {
	if integrationKey == "" {
		return nil, fmt.Errorf("pagerduty: integration key must not be empty")
	}
	if timeout <= 0 {
		timeout = 10 * time.Second
	}
	return &PagerDutyNotifier{
		integrationKey: integrationKey,
		client:         &http.Client{Timeout: timeout},
		eventsURL:      pagerdutyEventsURL,
	}, nil
}

// Notify sends the alert to PagerDuty as a trigger event.
func (p *PagerDutyNotifier) Notify(a Alert) error {
	body := pagerDutyPayload{
		RoutingKey:  p.integrationKey,
		EventAction: "trigger",
		Payload: pagerDutyDetail{
			Summary:  fmt.Sprintf("[%s] %s", a.Severity, a.Message),
			Severity: severityToPagerDuty(a.Severity),
			Source:   "portwatch",
		},
	}
	data, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("pagerduty: marshal payload: %w", err)
	}
	resp, err := p.client.Post(p.eventsURL, "application/json", bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("pagerduty: send event: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("pagerduty: unexpected status %d", resp.StatusCode)
	}
	return nil
}

func severityToPagerDuty(s string) string {
	switch s {
	case "critical":
		return "critical"
	case "warning":
		return "warning"
	default:
		return "info"
	}
}
