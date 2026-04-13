package config

import "fmt"

// FirehydrantConfig holds settings for the FireHydrant alert handler.
type FirehydrantConfig struct {
	Enabled    bool   `toml:"enabled"`
	APIKey     string `toml:"api_key"`
	ServiceID  string `toml:"service_id"`
	BaseURL    string `toml:"base_url"`
	Timeout    int    `toml:"timeout_seconds"`
}

func defaultFirehydrantConfig() FirehydrantConfig {
	return FirehydrantConfig{
		Enabled:   false,
		BaseURL:   "https://api.firehydrant.io/v1",
		Timeout:   10,
	}
}

func validateFirehydrant(cfg FirehydrantConfig) error {
	if !cfg.Enabled {
		return nil
	}
	if cfg.APIKey == "" {
		return fmt.Errorf("firehydrant: api_key is required")
	}
	if cfg.ServiceID == "" {
		return fmt.Errorf("firehydrant: service_id is required")
	}
	if cfg.BaseURL == "" {
		return fmt.Errorf("firehydrant: base_url must not be empty")
	}
	if cfg.Timeout <= 0 {
		return fmt.Errorf("firehydrant: timeout_seconds must be positive")
	}
	return nil
}
