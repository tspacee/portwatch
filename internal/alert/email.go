package alert

import (
	"fmt"
	"net/smtp"
	"strings"
)

// EmailConfig holds configuration for the email notifier.
type EmailConfig struct {
	SMTPHost string
	SMTPPort int
	Username string
	Password string
	From     string
	To       []string
}

// EmailNotifier sends alert notifications via email.
type EmailNotifier struct {
	cfg  EmailConfig
	send func(addr string, a smtp.Auth, from string, to []string, msg []byte) error
}

// NewEmailNotifier creates a new EmailNotifier with the given config.
// Returns an error if required fields are missing.
func NewEmailNotifier(cfg EmailConfig) (*EmailNotifier, error) {
	if cfg.SMTPHost == "" {
		return nil, fmt.Errorf("email notifier: smtp host is required")
	}
	if len(cfg.To) == 0 {
		return nil, fmt.Errorf("email notifier: at least one recipient is required")
	}
	if cfg.SMTPPort == 0 {
		cfg.SMTPPort = 587
	}
	return &EmailNotifier{
		cfg:  cfg,
		send: smtp.SendMail,
	}, nil
}

// Notify sends the alert as an email.
func (e *EmailNotifier) Notify(a Alert) error {
	subject := fmt.Sprintf("[portwatch] %s on port %d/%s", a.Event, a.Port, a.Protocol)
	body := fmt.Sprintf(
		"To: %s\r\nSubject: %s\r\nMIME-Version: 1.0\r\nContent-Type: text/plain\r\n\r\nRule: %s\nEvent: %s\nPort: %d\nProtocol: %s\nTimestamp: %s\n",
		strings.Join(e.cfg.To, ", "),
		subject,
		a.Rule,
		a.Event,
		a.Port,
		a.Protocol,
		a.Timestamp.Format("2006-01-02T15:04:05Z07:00"),
	)
	var auth smtp.Auth
	if e.cfg.Username != "" {
		auth = smtp.PlainAuth("", e.cfg.Username, e.cfg.Password, e.cfg.SMTPHost)
	}
	addr := fmt.Sprintf("%s:%d", e.cfg.SMTPHost, e.cfg.SMTPPort)
	if err := e.send(addr, auth, e.cfg.From, e.cfg.To, []byte(body)); err != nil {
		return fmt.Errorf("email notifier: send failed: %w", err)
	}
	return nil
}
