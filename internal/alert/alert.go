package alert

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/user/portwatch/internal/rules"
)

// Level represents the severity of an alert.
type Level string

const (
	LevelInfo  Level = "INFO"
	LevelWarn  Level = "WARN"
	LevelAlert Level = "ALERT"
)

// Alert represents a notification triggered by a rule violation.
type Alert struct {
	Timestamp time.Time
	Level     Level
	Violation rules.Violation
	Message   string
}

// Notifier sends alerts to a destination.
type Notifier interface {
	Notify(a Alert) error
}

// LogNotifier writes alerts as formatted lines to a writer.
type LogNotifier struct {
	Out io.Writer
}

// NewLogNotifier creates a LogNotifier that writes to stdout by default.
func NewLogNotifier(out io.Writer) *LogNotifier {
	if out == nil {
		out = os.Stdout
	}
	return &LogNotifier{Out: out}
}

// Notify formats and writes the alert to the configured writer.
func (n *LogNotifier) Notify(a Alert) error {
	_, err := fmt.Fprintf(
		n.Out,
		"[%s] %s | port=%d proto=%s rule=%s msg=%s\n",
		a.Timestamp.Format(time.RFC3339),
		a.Level,
		a.Violation.Port,
		a.Violation.Protocol,
		a.Violation.RuleName,
		a.Message,
	)
	return err
}

// ViolationToAlert converts a rules.Violation into an Alert.
func ViolationToAlert(v rules.Violation) Alert {
	return Alert{
		Timestamp: time.Now(),
		Level:     LevelAlert,
		Violation: v,
		Message:   fmt.Sprintf("unexpected port state: port %d/%s triggered rule '%s'", v.Port, v.Protocol, v.RuleName),
	}
}
