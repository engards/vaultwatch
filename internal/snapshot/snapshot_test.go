package snapshot_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/yourusername/vaultwatch/internal/monitor"
	"github.com/yourusername/vaultwatch/internal/snapshot"
)

func sampleStatuses() []monitor.SecretStatus {
	return []monitor.SecretStatus{
		{
			Path:      "secret/db/password",
			ExpiresAt: time.Now().Add(48 * time.Hour),
			TTL:       48 * time.Hour,
			Level:     "ok",
		},
		{
			Path:      "secret/api/key",
			ExpiresAt: time.Now().Add(2 * time.Hour),
			TTL:       2 * time.Hour,
			Level:     "critical",
		},
	}
}

func TestNew_CreatesDirectory(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "snapshots")
	_, err := snapshot.New(dir)
	if err != nil {
		t.Fatalf("New: unexpected error: %v", err)
	}
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		t.Error("expected directory to be created")
	}
}

func TestSave_WritesFile(t *testing.T) {
	m, err := snapshot.New(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}

	path, err := m.Save(sampleStatuses())
	if err != nil {
		t.Fatalf("Save: %v", err)
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Errorf("expected file %q to exist", path)
	}
}

func TestLatest_ReturnsNilWhenEmpty(t *testing.T) {
	m, err := snapshot.New(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}

	snap, err := m.Latest()
	if err != nil {
		t.Fatalf("Latest: unexpected error: %v", err)
	}
	if snap != nil {
		t.Error("expected nil snapshot for empty directory")
	}
}

func TestLatest_ReturnsMostRecent(t *testing.T) {
	m, err := snapshot.New(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}

	statuses := sampleStatuses()
	if _, err := m.Save(statuses); err != nil {
		t.Fatal(err)
	}

	snap, err := m.Latest()
	if err != nil {
		t.Fatalf("Latest: %v", err)
	}
	if snap == nil {
		t.Fatal("expected non-nil snapshot")
	}
	if len(snap.Statuses) != len(statuses) {
		t.Errorf("got %d statuses, want %d", len(snap.Statuses), len(statuses))
	}
	if snap.Statuses[0].Path != statuses[0].Path {
		t.Errorf("got path %q, want %q", snap.Statuses[0].Path, statuses[0].Path)
	}
}
