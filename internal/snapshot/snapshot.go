// Package snapshot provides functionality to capture and persist the current
// state of monitored Vault secrets to disk for diffing and historical review.
package snapshot

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/yourusername/vaultwatch/internal/monitor"
)

// Snapshot represents a point-in-time capture of secret statuses.
type Snapshot struct {
	CapturedAt time.Time              `json:"captured_at"`
	Statuses   []monitor.SecretStatus `json:"statuses"`
}

// Manager handles reading and writing snapshots to a directory.
type Manager struct {
	dir string
}

// New creates a new snapshot Manager that stores files in dir.
// The directory is created if it does not exist.
func New(dir string) (*Manager, error) {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, fmt.Errorf("snapshot: create directory %q: %w", dir, err)
	}
	return &Manager{dir: dir}, nil
}

// Save writes a new snapshot file containing the provided statuses.
// Files are named by UTC timestamp: 20060102T150405Z.json
func (m *Manager) Save(statuses []monitor.SecretStatus) (string, error) {
	snap := Snapshot{
		CapturedAt: time.Now().UTC(),
		Statuses:   statuses,
	}

	data, err := json.MarshalIndent(snap, "", "  ")
	if err != nil {
		return "", fmt.Errorf("snapshot: marshal: %w", err)
	}

	filename := snap.CapturedAt.Format("20060102T150405Z") + ".json"
	path := filepath.Join(m.dir, filename)

	if err := os.WriteFile(path, data, 0o644); err != nil {
		return "", fmt.Errorf("snapshot: write %q: %w", path, err)
	}

	return path, nil
}

// Latest loads the most recently saved snapshot from the directory.
// Returns nil, nil when no snapshots exist yet.
func (m *Manager) Latest() (*Snapshot, error) {
	entries, err := os.ReadDir(m.dir)
	if err != nil {
		return nil, fmt.Errorf("snapshot: read dir: %w", err)
	}

	// ReadDir returns entries sorted by name; last entry is most recent.
	var last string
	for _, e := range entries {
		if !e.IsDir() && filepath.Ext(e.Name()) == ".json" {
			last = e.Name()
		}
	}
	if last == "" {
		return nil, nil
	}

	data, err := os.ReadFile(filepath.Join(m.dir, last))
	if err != nil {
		return nil, fmt.Errorf("snapshot: read %q: %w", last, err)
	}

	var snap Snapshot
	if err := json.Unmarshal(data, &snap); err != nil {
		return nil, fmt.Errorf("snapshot: unmarshal: %w", err)
	}
	return &snap, nil
}
