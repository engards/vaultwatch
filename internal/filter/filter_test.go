package filter_test

import (
	"testing"
	"time"

	"github.com/yourusername/vaultwatch/internal/filter"
	"github.com/yourusername/vaultwatch/internal/monitor"
)

func sampleStatuses() []monitor.SecretStatus {
	now := time.Now()
	return []monitor.SecretStatus{
		{Path: "secret/app/db", State: monitor.StateOK, ExpiresAt: now.Add(72 * time.Hour)},
		{Path: "secret/app/api", State: monitor.StateWarning, ExpiresAt: now.Add(12 * time.Hour)},
		{Path: "secret/infra/tls", State: monitor.StateCritical, ExpiresAt: now.Add(2 * time.Hour)},
		{Path: "secret/infra/ssh", State: monitor.StateExpired, ExpiresAt: now.Add(-1 * time.Hour)},
	}
}

func TestApply_NoOptions_ReturnsAll(t *testing.T) {
	f := filter.New(filter.Options{})
	result := f.Apply(sampleStatuses())
	if len(result) != 4 {
		t.Fatalf("expected 4, got %d", len(result))
	}
}

func TestApply_PathPrefix(t *testing.T) {
	f := filter.New(filter.Options{PathPrefix: "secret/app"})
	result := f.Apply(sampleStatuses())
	if len(result) != 2 {
		t.Fatalf("expected 2, got %d", len(result))
	}
	for _, s := range result {
		if s.Path[:10] != "secret/app" {
			t.Errorf("unexpected path %s", s.Path)
		}
	}
}

func TestApply_StateFilter(t *testing.T) {
	f := filter.New(filter.Options{States: []string{"critical", "expired"}})
	result := f.Apply(sampleStatuses())
	if len(result) != 2 {
		t.Fatalf("expected 2, got %d", len(result))
	}
}

func TestApply_PathAndState(t *testing.T) {
	f := filter.New(filter.Options{
		PathPrefix: "secret/infra",
		States:     []string{"critical"},
	})
	result := f.Apply(sampleStatuses())
	if len(result) != 1 {
		t.Fatalf("expected 1, got %d", len(result))
	}
	if result[0].Path != "secret/infra/tls" {
		t.Errorf("unexpected path %s", result[0].Path)
	}
}

func TestApply_StateFilterCaseInsensitive(t *testing.T) {
	f := filter.New(filter.Options{States: []string{"WARNING"}})
	result := f.Apply(sampleStatuses())
	if len(result) != 1 {
		t.Fatalf("expected 1, got %d", len(result))
	}
}

func TestApply_EmptyInput(t *testing.T) {
	f := filter.New(filter.Options{PathPrefix: "secret/app"})
	result := f.Apply([]monitor.SecretStatus{})
	if len(result) != 0 {
		t.Fatalf("expected 0, got %d", len(result))
	}
}
