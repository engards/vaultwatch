// Package audit provides structured audit logging for VaultWatch events.
//
// Each event is written as a single JSON line to a configurable io.Writer,
// making the output easy to pipe into log aggregators such as Loki, Splunk,
// or a simple file sink.
//
// Supported event types:
//
//   - EventSecretChecked    — a secret path was inspected during a monitor cycle
//   - EventSecretExpired    — a secret has passed its expiry threshold
//   - EventRenewalTriggered — the renewer attempted to extend a lease
//   - EventRenewalFailed    — a renewal attempt returned an error
//   - EventAlertSent        — an alert notification was dispatched
//
// Example usage:
//
//	l := audit.New(os.Stderr)
//	l.Log(audit.EventRenewalTriggered, "secret/myapp/db", "lease renewed",
//	    map[string]string{"lease_id": id, "ttl": "3600"})
package audit
