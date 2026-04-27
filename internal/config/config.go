package config

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// Config holds the top-level vaultwatch configuration.
type Config struct {
	Vault   VaultConfig   `yaml:"vault"`
	Alerting AlertingConfig `yaml:"alerting"`
	PollInterval time.Duration `yaml:"poll_interval"`
}

// VaultConfig holds Vault connection settings.
type VaultConfig struct {
	Address   string `yaml:"address"`
	Token     string `yaml:"token"`
	Namespace string `yaml:"namespace"`
	CACert    string `yaml:"ca_cert"`
}

// AlertingConfig holds alerting thresholds and notification settings.
type AlertingConfig struct {
	WarnThreshold     time.Duration `yaml:"warn_threshold"`
	CriticalThreshold time.Duration `yaml:"critical_threshold"`
	SlackWebhookURL   string        `yaml:"slack_webhook_url"`
}

// Load reads and parses the config file at the given path.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading config file %q: %w", path, err)
	}

	cfg := &Config{
		PollInterval:          60 * time.Second,
		Alerting: AlertingConfig{
			WarnThreshold:     7 * 24 * time.Hour,
			CriticalThreshold: 24 * time.Hour,
		},
	}

	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("parsing config file: %w", err)
	}

	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	return cfg, nil
}

func (c *Config) validate() error {
	if c.Vault.Address == "" {
		return fmt.Errorf("vault.address is required")
	}
	if c.Vault.Token == "" {
		c.Vault.Token = os.Getenv("VAULT_TOKEN")
	}
	if c.Vault.Token == "" {
		return fmt.Errorf("vault.token is required (or set VAULT_TOKEN env var)")
	}
	if c.PollInterval < 5*time.Second {
		return fmt.Errorf("poll_interval must be at least 5s")
	}
	return nil
}
