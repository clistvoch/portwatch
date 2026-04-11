package config

import "fmt"

// DiscordConfig holds configuration for the Discord webhook alert handler.
type DiscordConfig struct {
	Enabled    bool   `toml:"enabled"`
	WebhookURL string `toml:"webhook_url"`
	Username   string `toml:"username"`
	AvatarURL  string `toml:"avatar_url"`
}

func defaultDiscordConfig() DiscordConfig {
	return DiscordConfig{
		Enabled:  false,
		Username: "portwatch",
	}
}

func validateDiscord(c DiscordConfig) error {
	if !c.Enabled {
		return nil
	}
	if c.WebhookURL == "" {
		return fmt.Errorf("discord: webhook_url is required when enabled")
	}
	return nil
}
