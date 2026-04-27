package monitor

import (
	"testing"
	"time"
)

func TestClassify(t *testing.T) {
	warnBefore := 24 * time.Hour

	tests := []struct {
		name     string
		ttl      time.Duration
		want     string
	}{
		{"expired", 0, "expired"},
		{"negative ttl", -1 * time.Second, "expired"},
		{"within warning window", 12 * time.Hour, "warning"},
		{"at warning boundary", warnBefore, "warning"},
		{"ok", 48 * time.Hour, "ok"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := classify(tc.ttl, warnBefore)
			if got != tc.want {
				t.Errorf("classify(%v, %v) = %q; want %q", tc.ttl, warnBefore, got, tc.want)
			}
		})
	}
}

func TestSecretStatus_Fields(t *testing.T) {
	now := time.Now()
	ttl := 72 * time.Hour

	s := &SecretStatus{
		Path:      "secret/data/myapp/db",
		ExpiresAt: now.Add(ttl),
		TTL:       ttl,
		Status:    "ok",
	}

	if s.Path != "secret/data/myapp/db" {
		t.Errorf("unexpected Path: %q", s.Path)
	}
	if s.Status != "ok" {
		t.Errorf("unexpected Status: %q", s.Status)
	}
	if s.TTL != ttl {
		t.Errorf("unexpected TTL: %v", s.TTL)
	}
}

func TestNew_SetsFields(t *testing.T) {
	warn := 24 * time.Hour
	poll := 30 * time.Second

	m := New(nil, warn, poll)

	if m.warnBefore != warn {
		t.Errorf("warnBefore = %v; want %v", m.warnBefore, warn)
	}
	if m.pollInterval != poll {
		t.Errorf("pollInterval = %v; want %v", m.pollInterval, poll)
	}
}
