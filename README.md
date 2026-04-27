# vaultwatch

A lightweight CLI tool to monitor and alert on HashiCorp Vault secret expiration and lease renewals.

---

## Installation

```bash
go install github.com/yourusername/vaultwatch@latest
```

Or download a pre-built binary from the [Releases](https://github.com/yourusername/vaultwatch/releases) page.

---

## Usage

Set your Vault address and token, then run `vaultwatch` against a path:

```bash
export VAULT_ADDR="https://vault.example.com"
export VAULT_TOKEN="s.your-token-here"

# Monitor secrets at a path and alert if expiring within 7 days
vaultwatch monitor --path secret/myapp --threshold 7d

# List all leases and their expiration times
vaultwatch leases --path aws/creds

# Renew a specific lease
vaultwatch renew --lease-id aws/creds/myapp/abc123
```

### Flags

| Flag | Description | Default |
|------|-------------|---------|
| `--path` | Vault secret path to monitor | _(required)_ |
| `--threshold` | Alert window before expiration | `24h` |
| `--interval` | How often to poll Vault | `5m` |
| `--output` | Output format: `text`, `json` | `text` |

---

## Configuration

`vaultwatch` respects standard Vault environment variables:

- `VAULT_ADDR` — Vault server address
- `VAULT_TOKEN` — Authentication token
- `VAULT_NAMESPACE` — Vault namespace (Enterprise)

---

## Requirements

- Go 1.21+
- HashiCorp Vault 1.10+

---

## License

[MIT](LICENSE) © 2024 yourusername