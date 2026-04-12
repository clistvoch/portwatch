package config

import "fmt"

// DatadogConfig holds settings for the Datadog alert handler.
type DatadogConfig struct {
	Enabled bool   `toml:"enabled"`
	APIKey  string `toml:"api_key"`
	Site    string `toml:"site"`
	Service string `toml:"service"`
	Tags    []string `toml:"tags"`
}

func defaultDatadogConfig() DatadogConfig {
	return DatadogConfig{
		Enabled: false,
		Site:    "datadoghq.com",
		Service: "portwatch",
		Tags:    []string{},
	}
}

func validateDatadog(cfg DatadogConfig) error {
	if !cfg.Enabled {
		return nil
	}
	if cfg.APIKey == "" {
		return &ValidationError{Field: "datadog.api_key", Msg: "api_key is required when datadog is enabled"}
	}
	if cfg.Site == "" {
		return &ValidationError{Field: "datadog.site", Msg: "site must not be empty"}
	}
	if cfg.Service == "" {
		return fmt.Errorf("datadog.service must not be empty")
	}
	return nil
}
