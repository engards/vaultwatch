package monitor

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/yourusername/vaultwatch/internal/vault"
)

// SecretStatus represents the expiration status of a secret.
type SecretStatus struct {
	Path      string
	ExpiresAt time.Time
	TTL       time.Duration
	Status    string // "ok", "warning", "expired"
}

// Monitor watches Vault secrets and reports their expiration status.
type Monitor struct {
	client      *vault.Client
	warnBefore  time.Duration
	pollInterval time.Duration
}

// New creates a new Monitor with the given Vault client and thresholds.
func New(client *vault.Client, warnBefore, pollInterval time.Duration) *Monitor {
	return &Monitor{
		client:      client,
		warnBefore:  warnBefore,
		pollInterval: pollInterval,
	}
}

// CheckSecret fetches the secret at path and returns its status.
func (m *Monitor) CheckSecret(ctx context.Context, path string) (*SecretStatus, error) {
	secret, err := m.client.ReadSecret(ctx, path)
	if err != nil {
		return nil, fmt.Errorf("reading secret %q: %w", path, err)
	}
	if secret == nil {
		return nil, fmt.Errorf("secret not found at path %q", path)
	}

	ttl := time.Duration(secret.LeaseDuration) * time.Second
	expiresAt := time.Now().Add(ttl)

	status := classify(ttl, m.warnBefore)

	return &SecretStatus{
		Path:      path,
		ExpiresAt: expiresAt,
		TTL:       ttl,
		Status:    status,
	}, nil
}

// Run continuously monitors the given paths, calling onStatus for each result.
func (m *Monitor) Run(ctx context.Context, paths []string, onStatus func(*SecretStatus)) error {
	ticker := time.NewTicker(m.pollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			for _, path := range paths {
				status, err := m.CheckSecret(ctx, path)
				if err != nil {
					log.Printf("[error] %v", err)
					continue
				}
				onStatus(status)
			}
		}
	}
}

func classify(ttl, warnBefore time.Duration) string {
	switch {
	case ttl <= 0:
		return "expired"
	case ttl <= warnBefore:
		return "warning"
	default:
		return "ok"
	}
}
