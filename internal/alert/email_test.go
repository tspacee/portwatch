package alert

import (
	"errors"
	"net/smtp"
	"strings"
	"testing"
	"time"
)

func TestNewEmailNotifier_MissingSMTPHost(t *testing.T) {
	_, err := NewEmailNotifier(EmailConfig{To: []string{"a@b.com"}})
	if err == nil {
		t.Fatal("expected error for missing smtp host")
	}
}

func TestNewEmailNotifier_MissingRecipients(t *testing.T) {
	_, err := NewEmailNotifier(EmailConfig{SMTPHost: "smtp.example.com"})
	if err == nil {
		t.Fatal("expected error for missing recipients")
	}
}

func TestNewEmailNotifier_DefaultPort(t *testing.T) {
	n, err := NewEmailNotifier(EmailConfig{
		SMTPHost: "smtp.example.com",
		To:       []string{"a@b.com"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n.cfg.SMTPPort != 587 {
		t.Errorf("expected default port 587, got %d", n.cfg.SMTPPort)
	}
}

func TestEmailNotifier_Notify_Success(t *testing.T) {
	var capturedAddr string
	var capturedMsg []byte

	n, _ := NewEmailNotifier(EmailConfig{
		SMTPHost: "smtp.example.com",
		SMTPPort: 587,
		From:     "from@example.com",
		To:       []string{"to@example.com"},
	})
	n.send = func(addr string, _ smtp.Auth, _ string, _ []string, msg []byte) error {
		capturedAddr = addr
		capturedMsg = msg
		return nil
	}

	a := Alert{
		Rule:      "block-ssh",
		Event:     "unexpected_open",
		Port:      22,
		Protocol:  "tcp",
		Timestamp: time.Now(),
	}
	if err := n.Notify(a); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if capturedAddr != "smtp.example.com:587" {
		t.Errorf("unexpected addr: %s", capturedAddr)
	}
	if !strings.Contains(string(capturedMsg), "block-ssh") {
		t.Error("expected rule name in message body")
	}
}

func TestEmailNotifier_Notify_SendError(t *testing.T) {
	n, _ := NewEmailNotifier(EmailConfig{
		SMTPHost: "smtp.example.com",
		To:       []string{"to@example.com"},
	})
	n.send = func(_ string, _ smtp.Auth, _ string, _ []string, _ []byte) error {
		return errors.New("connection refused")
	}

	err := n.Notify(Alert{Port: 80, Protocol: "tcp", Timestamp: time.Now()})
	if err == nil {
		t.Fatal("expected error from send failure")
	}
	if !strings.Contains(err.Error(), "send failed") {
		t.Errorf("unexpected error message: %v", err)
	}
}
