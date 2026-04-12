package config

import "fmt"

type LokiConfig struct {
	Enabled  bool   `toml:"enabled"`
	URL      string `toml:"url"`
	TenantID string `toml:"tenant_id"`
	Labels   map[string]string `toml:"labels"`
	Timeout  int    `toml:"timeout_seconds"`
}

func defaultLokiConfig() LokiConfig {
	return LokiConfig{
		Enabled:  false,
		URL:      "",
		TenantID: "",
		Labels:   map[string]string{"app": "portwatch"},
		Timeout:  5,
	}
}

func validateLoki(c LokiConfig) error {
	if !c.Enabled {
		return nil
	}
	if c.URL == "" {
		return &ValidationError{Field: "loki.url", Message: "url is required when loki is enabled"}
	}
	if c.Timeout <= 0 {
		return fmt.Errorf("loki.timeout_seconds must be greater than 0")
	}
	return nil
}
