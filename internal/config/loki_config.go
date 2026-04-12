package config

import "fmt"

// LokiConfig holds configuration for the Loki alert handler.
type LokiConfig struct {
	Enabled bool              `toml:"enabled"`
	URL     string            `toml:"url"`
	Labels  map[string]string `toml:"labels"`
}

func defaultLokiConfig() LokiConfig {
	return LokiConfig{
		Enabled: false,
		URL:     "",
		Labels:  map[string]string{"job": "portwatch"},
	}
}

func validateLoki(c LokiConfig) error {
	if !c.Enabled {
		return nil
	}
	if c.URL == "" {
		return fmt.Errorf("loki: url is required when enabled")
	}
	return nil
}
