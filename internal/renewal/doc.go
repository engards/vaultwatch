// Package renewal provides automatic lease and token renewal logic for VaultWatch.
//
// It exposes a Renewer that evaluates whether a secret's remaining TTL has
// fallen within a configurable grace window and, if so, calls the Vault API
// to extend the lease.
//
// Typical usage:
//
//	r := renewal.New(vaultClient, logger, 0.2) // renew when 20% TTL remains
//	if r.ShouldRenew(secret.TTL, secret.MaxTTL) {
//		result := r.RenewLease(ctx, secret.LeaseID, 3600)
//		if result.Err != nil {
//			// handle error
//		}
//	}
package renewal
