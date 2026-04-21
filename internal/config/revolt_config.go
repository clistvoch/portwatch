package config

import "fmt"

// RevoltConfig holds configuration for the Revolt alert handler.
type RevoltConfig struct {
	Enabled   bool   `toml:"enabled"`
	WebhookURL string `toml:"webhook_url"`
	Username  string `toml:"username"`
	AvatarURL string `toml:"avatar_url"`
}

func defaultRevoltConfig() RevoltConfig {
	return RevoltConfig{
		Enabled:   false,
		WebhookURL: "",
		Username:  "portwatch",
		AvatarURL: "",
	}
}

func validateRevolt(c RevoltConfig) error {
	if !c.Enabled {
		return nil
	}
	if c.WebhookURL == "" {
		return fmt.Errorf("revolt: webhook_url is required")
	}
	return nil
}
