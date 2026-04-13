package config

import "fmt"

type MattermostConfig struct {
	Enabled   bool   `toml:"enabled"`
	WebhookURL string `toml:"webhook_url"`
	Channel   string `toml:"channel"`
	Username  string `toml:"username"`
	IconEmoji string `toml:"icon_emoji"`
}

func defaultMattermostConfig() MattermostConfig {
	return MattermostConfig{
		Enabled:   false,
		Channel:   "",
		Username:  "portwatch",
		IconEmoji: ":shield:",
	}
}

func validateMattermost(cfg MattermostConfig) error {
	if !cfg.Enabled {
		return nil
	}
	if cfg.WebhookURL == "" {
		return fmt.Errorf("mattermost: webhook_url is required")
	}
	return nil
}
