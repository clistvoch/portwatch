package config

import "fmt"

type GrafanaConfig struct {
	Enabled    bool   `toml:"enabled"`
	URL        string `toml:"url"`
	APIKey     string `toml:"api_key"`
	DashboardID string `toml:"dashboard_id"`
	Timeout    int    `toml:"timeout_seconds"`
}

func defaultGrafanaConfig() GrafanaConfig {
	return GrafanaConfig{
		Enabled: false,
		Timeout: 5,
	}
}

func validateGrafana(cfg GrafanaConfig) error {
	if !cfg.Enabled {
		return nil
	}
	if cfg.URL == "" {
		return fmt.Errorf("grafana: url is required")
	}
	if cfg.APIKey == "" {
		return fmt.Errorf("grafana: api_key is required")
	}
	if cfg.Timeout <= 0 {
		return fmt.Errorf("grafana: timeout_seconds must be positive")
	}
	return nil
}
