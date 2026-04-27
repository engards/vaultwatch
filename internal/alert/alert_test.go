package alert

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/yourusername/vaultwatch/internal/monitor"
)

func TestToLevel(t *testing.T) {
	tests := []struct {
		input string
		want  Level
	}{
		{"ok", LevelOK},
		{"warning", LevelWarning},
		{"expired", LevelExpired},
		{"unknown", LevelOK},
	}
	for _, tc := range tests {
		got := toLevel(tc.input)
		if got != tc.want {
			t.Errorf("toLevel(%q) = %q; want %q", tc.input, got, tc.want)
		}
	}
}

func TestShouldAlert(t *testing.T) {
	tests := []struct {
		status string
		want   bool
	}{
		{"ok", false},
		{"warning", true},
		{"expired", true},
	}
	for _, tc := range tests {
		s := &monitor.SecretStatus{Status: tc.status}
		if got := ShouldAlert(s); got != tc.want {
			t.Errorf("ShouldAlert(%q) = %v; want %v", tc.status, got, tc.want)
		}
	}
}

func TestNotify_Output(t *testing.T) {
	var buf bytes.Buffer
	n := New(&buf)

	s := &monitor.SecretStatus{
		Path:      "secret/data/test",
		TTL:       6 * time.Hour,
		ExpiresAt: time.Now().Add(6 * time.Hour),
		Status:    "warning",
	}

	n.Notify(s)

	out := buf.String()
	if !strings.Contains(out, "WARNING") {
		t.Errorf("expected WARNING in output, got: %s", out)
	}
	if !strings.Contains(out, "secret/data/test") {
		t.Errorf("expected path in output, got: %s", out)
	}
}

func TestNew_DefaultsToStdout(t *testing.T) {
	n := New(nil)
	if n.out == nil {
		t.Error("expected non-nil writer")
	}
}
