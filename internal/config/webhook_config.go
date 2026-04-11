package config

import "fmt"

// WebhookConfig holds configuration for the webhook alert handler.
type WebhookConfig struct {
	Enabled bool   `toml:"enabled"`
	URL     string `toml:"url"`
	Secret  string `toml:"secret"`
	Timeout int    `toml:"timeout_seconds"`
}

func defaultWebhookConfig() WebhookConfig {
	return WebhookConfig{
		Enabled: false,
		URL:     "",
		Secret:  "",
		Timeout: 10,
	}
}

func validateWebhook(cfg WebhookConfig) error {
	if !cfg.Enabled {
		return nil
	}
	if cfg.URL == "" {
		return &ValidationError{Field: "webhook.url", Message: "url is required when webhook is enabled"}
	}
	if cfg.Timeout <= 0 {
		return fmt.Errorf("webhook.timeout_seconds must be positive, got %d", cfg.Timeout)
	}
	return nil
}
