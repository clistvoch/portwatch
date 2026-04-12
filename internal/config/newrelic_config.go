package config

import "fmt"

type NewRelicConfig struct {
	Enabled    bool   `toml:"enabled"`
	APIKey     string `toml:"api_key"`
	AccountID  string `toml:"account_id"`
	Region     string `toml:"region"`
	EventType  string `toml:"event_type"`
	TimeoutSec int    `toml:"timeout_sec"`
}

func defaultNewRelicConfig() NewRelicConfig {
	return NewRelicConfig{
		Enabled:    false,
		Region:     "US",
		EventType:  "PortWatchAlert",
		TimeoutSec: 10,
	}
}

func validateNewRelic(cfg NewRelicConfig) error {
	if !cfg.Enabled {
		return nil
	}
	if cfg.APIKey == "" {
		return ValidationError{Field: "newrelic.api_key", Message: "api_key is required"}
	}
	if cfg.AccountID == "" {
		return ValidationError{Field: "newrelic.account_id", Message: "account_id is required"}
	}
	if cfg.Region != "US" && cfg.Region != "EU" {
		return ValidationError{Field: "newrelic.region", Message: fmt.Sprintf("invalid region %q: must be US or EU", cfg.Region)}
	}
	if cfg.TimeoutSec <= 0 {
		return ValidationError{Field: "newrelic.timeout_sec", Message: "timeout_sec must be positive"}
	}
	return nil
}
