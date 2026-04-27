package alert

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/yourusername/vaultwatch/internal/monitor"
)

// Level represents the severity of an alert.
type Level string

const (
	LevelOK      Level = "OK"
	LevelWarning Level = "WARNING"
	LevelExpired Level = "EXPIRED"
)

// Notifier sends alerts for secret status changes.
type Notifier struct {
	out io.Writer
}

// New creates a Notifier that writes to the given writer.
// If w is nil, os.Stdout is used.
func New(w io.Writer) *Notifier {
	if w == nil {
		w = os.Stdout
	}
	return &Notifier{out: w}
}

// Notify formats and writes an alert for the given secret status.
func (n *Notifier) Notify(s *monitor.SecretStatus) {
	level := toLevel(s.Status)
	timestamp := time.Now().Format(time.RFC3339)

	fmt.Fprintf(
		n.out,
		"[%s] %s | path=%s ttl=%s expires=%s\n",
		level,
		timestamp,
		s.Path,
		s.TTL.Round(time.Second),
		s.ExpiresAt.Format(time.RFC3339),
	)
}

// ShouldAlert returns true if the status warrants an alert (non-OK).
func ShouldAlert(s *monitor.SecretStatus) bool {
	return s.Status == "warning" || s.Status == "expired"
}

func toLevel(status string) Level {
	switch status {
	case "warning":
		return LevelWarning
	case "expired":
		return LevelExpired
	default:
		return LevelOK
	}
}
