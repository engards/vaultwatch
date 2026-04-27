// Package vault provides a client wrapper around the HashiCorp Vault API
// for interacting with secrets and lease information.
package vault

import (
	"context"
	"fmt"
	"time"

	vaultapi "github.com/hashicorp/vault/api"

	"github.com/yourusername/vaultwatch/internal/config"
)

// Client wraps the Vault API client with vaultwatch-specific functionality.
type Client struct {
	api     *vaultapi.Client
	cfg     *config.Config
}

// SecretInfo holds metadata about a Vault secret relevant to expiration monitoring.
type SecretInfo struct {
	Path        string
	LeaseID     string
	LeaseDuration time.Duration
	Renewable   bool
	ExpireTime  time.Time
}

// New creates a new Vault client using the provided configuration.
func New(cfg *config.Config) (*Client, error) {
	vaultCfg := vaultapi.DefaultConfig()
	vaultCfg.Address = cfg.VaultAddr

	if err := vaultCfg.ConfigureTLS(&vaultapi.TLSConfig{
		Insecure: cfg.TLSSkipVerify,
	}); err != nil {
		return nil, fmt.Errorf("configuring vault TLS: %w", err)
	}

	apiClient, err := vaultapi.NewClient(vaultCfg)
	if err != nil {
		return nil, fmt.Errorf("creating vault api client: %w", err)
	}

	apiClient.SetToken(cfg.VaultToken)

	return &Client{
		api: apiClient,
		cfg: cfg,
	}, nil
}

// ReadSecret reads a secret at the given path and returns its metadata.
// For KV v2 secrets, the path should include the 'data/' prefix.
func (c *Client) ReadSecret(ctx context.Context, path string) (*SecretInfo, error) {
	secret, err := c.api.Logical().ReadWithContext(ctx, path)
	if err != nil {
		return nil, fmt.Errorf("reading secret at %q: %w", path, err)
	}
	if secret == nil {
		return nil, fmt.Errorf("secret not found at path %q", path)
	}

	info := &SecretInfo{
		Path:          path,
		LeaseID:       secret.LeaseID,
		LeaseDuration: time.Duration(secret.LeaseDuration) * time.Second,
		Renewable:     secret.Renewable,
	}

	if secret.LeaseDuration > 0 {
		info.ExpireTime = time.Now().Add(info.LeaseDuration)
	}

	return info, nil
}

// RenewSecret attempts to renew the lease for a secret by its lease ID.
// Returns the updated SecretInfo reflecting the new lease duration.
func (c *Client) RenewSecret(ctx context.Context, leaseID string, increment int) (*SecretInfo, error) {
	secret, err := c.api.Sys().RenewWithContext(ctx, leaseID, increment)
	if err != nil {
		return nil, fmt.Errorf("renewing lease %q: %w", leaseID, err)
	}
	if secret == nil {
		return nil, fmt.Errorf("no response when renewing lease %q", leaseID)
	}

	info := &SecretInfo{
		LeaseID:       secret.LeaseID,
		LeaseDuration: time.Duration(secret.LeaseDuration) * time.Second,
		Renewable:     secret.Renewable,
	}

	if secret.LeaseDuration > 0 {
		info.ExpireTime = time.Now().Add(info.LeaseDuration)
	}

	return info, nil
}

// Ping verifies connectivity to Vault by checking the health endpoint.
func (c *Client) Ping(ctx context.Context) error {
	_, err := c.api.Sys().HealthWithContext(ctx)
	if err != nil {
		return fmt.Errorf("vault health check failed: %w", err)
	}
	return nil
}
