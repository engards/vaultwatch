// Package audit provides structured audit logging for secret access and renewal events.
package audit

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"
)

// EventType categorizes the kind of audit event.
type EventType string

const (
	EventSecretChecked EventType = "secret_checked"
	EventSecretExpired EventType = "secret_expired"
	EventRenewalTriggered EventType = "renewal_triggered"
	EventRenewalFailed EventType = "renewal_failed"
	EventAlertSent EventType = "alert_sent"
)

// Event represents a single audit log entry.
type Event struct {
	Timestamp time.Time `json:"timestamp"`
	Type      EventType `json:"type"`
	Path      string    `json:"path"`
	Message   string    `json:"message"`
	Meta      map[string]string `json:"meta,omitempty"`
}

// Logger writes structured audit events to an output sink.
type Logger struct {
	out io.Writer
}

// New creates a new audit Logger. If w is nil, os.Stdout is used.
func New(w io.Writer) *Logger {
	if w == nil {
		w = os.Stdout
	}
	return &Logger{out: w}
}

// Log writes a structured JSON audit event.
func (l *Logger) Log(eventType EventType, path, message string, meta map[string]string) error {
	e := Event{
		Timestamp: time.Now().UTC(),
		Type:      eventType,
		Path:      path,
		Message:   message,
		Meta:      meta,
	}
	b, err := json.Marshal(e)
	if err != nil {
		return fmt.Errorf("audit: marshal event: %w", err)
	}
	_, err = fmt.Fprintln(l.out, string(b))
	return err
}
