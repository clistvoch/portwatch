package config

import "fmt"

// WebhookTransformConfig holds settings for payload transformation on outbound webhooks.
type WebhookTransformConfig struct {
	Enabled    bool   `toml:"enabled"`
	Template   string `toml:"template"`
	ContentType string `toml:"content_type"`
	IncludeHost bool   `toml:"include_host"`
}

func defaultWebhookTransformConfig() WebhookTransformConfig {
	return WebhookTransformConfig{
		Enabled:     false,
		Template:    "",
		ContentType: "application/json",
		IncludeHost: true,
	}
}

func validateWebhookTransform(c WebhookTransformConfig) error {
	if !c.Enabled {
		return nil
	}
	if c.ContentType == "" {
		return fmt.Errorf("webhook_transform: content_type must not be empty")
	}
	validTypes := map[string]bool{
		"application/json": true,
		"application/x-www-form-urlencoded": true,
		"text/plain": true,
	}
	if !validTypes[c.ContentType] {
		return fmt.Errorf("webhook_transform: unsupported content_type %q", c.ContentType)
	}
	return nil
}
