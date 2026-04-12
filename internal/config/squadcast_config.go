package config

import "fmt"

type SquadcastConfig struct {
	Enabled    bool   `toml:"enabled"`
	WebhookURL string `toml:"webhook_url"`
	Environment string `toml:"environment"`
	Timeout    int    `toml:"timeout_seconds"`
}

func defaultSquadcastConfig() SquadcastConfig {
	return SquadcastConfig{
		Enabled:     false,
		WebhookURL:  "",
		Environment: "production",
		Timeout:     10,
	}
}

func validateSquadcast(cfg SquadcastConfig) error {
	if !cfg.Enabled {
		return nil
	}
	if cfg.WebhookURL == "" {
		return fmt.Errorf("squadcast: webhook_url is required")
	}
	if cfg.Timeout <= 0 {
		return fmt.Errorf("squadcast: timeout_seconds must be positive")
	}
	return nil
}
