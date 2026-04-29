package output

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/vaultwatch/internal/monitor"
)

func sampleStatuses() []monitor.SecretStatus {
	return []monitor.SecretStatus{
		{
			Path:      "secret/db/password",
			Status:    "warning",
			TTL:       45 * time.Minute,
			Renewable: true,
		},
		{
			Path:      "secret/api/key",
			Status:    "critical",
			TTL:       -1 * time.Second,
			Renewable: false,
		},
	}
}

func TestWrite_TableFormat(t *testing.T) {
	var buf bytes.Buffer
	f := New(FormatTable, &buf)
	if err := f.Write(sampleStatuses()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "PATH") {
		t.Error("expected table header PATH")
	}
	if !strings.Contains(out, "secret/db/password") {
		t.Error("expected path secret/db/password in output")
	}
	if !strings.Contains(out, "expired") {
		t.Error("expected 'expired' for negative TTL")
	}
	if !strings.Contains(out, "yes") {
		t.Error("expected 'yes' for renewable secret")
	}
}

func TestWrite_JSONFormat(t *testing.T) {
	var buf bytes.Buffer
	f := New(FormatJSON, &buf)
	if err := f.Write(sampleStatuses()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, `"Path"`) && !strings.Contains(out, `"path"`) {
		t.Error("expected JSON output to contain path field")
	}
	if !strings.Contains(out, "secret/api/key") {
		t.Error("expected path secret/api/key in JSON output")
	}
}

func TestFormatDuration(t *testing.T) {
	cases := []struct {
		d    time.Duration
		want string
	}{
		{2*time.Hour + 15*time.Minute, "2h15m"},
		{30*time.Minute + 10*time.Second, "30m10s"},
		{45 * time.Second, "45s"},
		{-1 * time.Second, "expired"},
		{0, "expired"},
	}
	for _, tc := range cases {
		got := formatDuration(tc.d)
		if got != tc.want {
			t.Errorf("formatDuration(%v) = %q, want %q", tc.d, got, tc.want)
		}
	}
}

func TestNew_DefaultsToStdout(t *testing.T) {
	f := New(FormatTable, nil)
	if f.w == nil {
		t.Error("expected non-nil writer when nil passed")
	}
}
