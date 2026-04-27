package config

import (
	"os"
	"testing"
	"time"
)

func writeTempConfig(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "vaultwatch-*.yaml")
	if err != nil {
		t.Fatalf("creating temp file: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("writing temp file: %v", err)
	}
	f.Close()
	return f.Name()
}

func TestLoad_ValidConfig(t *testing.T) {
	path := writeTempConfig(t, `
vault:
  address: "https://vault.example.com"
  token: "s.testtoken"
poll_interval: 30s
alerting:
  warn_threshold: 48h
  critical_threshold: 12h
`)

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Vault.Address != "https://vault.example.com" {
		t.Errorf("expected vault address, got %q", cfg.Vault.Address)
	}
	if cfg.PollInterval != 30*time.Second {
		t.Errorf("expected 30s poll interval, got %v", cfg.PollInterval)
	}
	if cfg.Alerting.WarnThreshold != 48*time.Hour {
		t.Errorf("expected 48h warn threshold, got %v", cfg.Alerting.WarnThreshold)
	}
}

func TestLoad_Defaults(t *testing.T) {
	path := writeTempConfig(t, `
vault:
  address: "https://vault.example.com"
  token: "s.testtoken"
`)

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.PollInterval != 60*time.Second {
		t.Errorf("expected default 60s poll interval, got %v", cfg.PollInterval)
	}
	if cfg.Alerting.CriticalThreshold != 24*time.Hour {
		t.Errorf("expected default 24h critical threshold, got %v", cfg.Alerting.CriticalThreshold)
	}
}

func TestLoad_MissingToken_FallsBackToEnv(t *testing.T) {
	t.Setenv("VAULT_TOKEN", "env-token")
	path := writeTempConfig(t, `
vault:
  address: "https://vault.example.com"
`)
	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Vault.Token != "env-token" {
		t.Errorf("expected token from env, got %q", cfg.Vault.Token)
	}
}

func TestLoad_MissingAddress_ReturnsError(t *testing.T) {
	path := writeTempConfig(t, `
vault:
  token: "s.testtoken"
`)
	_, err := Load(path)
	if err == nil {
		t.Fatal("expected error for missing vault address")
	}
}

func TestLoad_FileNotFound(t *testing.T) {
	_, err := Load("/nonexistent/path/config.yaml")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}
