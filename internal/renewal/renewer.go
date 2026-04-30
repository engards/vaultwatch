package renewal

import (
	"context"
	"fmt"
	"log"
	"time"

	vaultapi "github.com/hashicorp/vault/api"
)

// RenewResult holds the outcome of a renewal attempt.
type RenewResult struct {
	Path      string
	Renewed   bool
	NewExpiry time.Time
	Err       error
}

// Renewer handles automatic renewal of Vault leases and tokens.
type Renewer struct {
	client    *vaultapi.Client
	logger    *log.Logger
	gracePct  float64 // renew when this fraction of TTL remains
}

// New creates a Renewer with the given Vault client.
// gracePct controls how early renewal is triggered (e.g. 0.2 = 20% of TTL left).
func New(client *vaultapi.Client, logger *log.Logger, gracePct float64) *Renewer {
	if gracePct <= 0 || gracePct >= 1 {
		gracePct = 0.2
	}
	if logger == nil {
		logger = log.Default()
	}
	return &Renewer{
		client:   client,
		logger:   logger,
		gracePct: gracePct,
	}
}

// ShouldRenew returns true when the remaining TTL is within the grace window.
func (r *Renewer) ShouldRenew(ttl, maxTTL time.Duration) bool {
	if maxTTL <= 0 {
		return false
	}
	threshold := time.Duration(float64(maxTTL) * r.gracePct)
	return ttl <= threshold
}

// RenewLease attempts to renew the lease identified by leaseID.
func (r *Renewer) RenewLease(ctx context.Context, leaseID string, increment int) RenewResult {
	res := RenewResult{Path: leaseID}

	secret, err := r.client.Sys().RenewWithContext(ctx, leaseID, increment)
	if err != nil {
		res.Err = fmt.Errorf("renew lease %q: %w", leaseID, err)
		r.logger.Printf("[renewal] failed to renew lease %q: %v", leaseID, err)
		return res
	}

	res.Renewed = true
	res.NewExpiry = time.Now().Add(time.Duration(secret.LeaseDuration) * time.Second)
	r.logger.Printf("[renewal] renewed lease %q, new expiry: %s", leaseID, res.NewExpiry.Format(time.RFC3339))
	return res
}
