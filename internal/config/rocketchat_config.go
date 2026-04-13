package config

import "fmt"

type RocketChatConfig struct {
	Enabled    bool   `toml:"enabled"`
	WebhookURL string `toml:"webhook_url"`
	Channel    string `toml:"channel"`
	Username   string `toml:"username"`
	IconEmoji  string `toml:"icon_emoji"`
}

func defaultRocketChatConfig() RocketChatConfig {
	return RocketChatConfig{
		Enabled:   false,
		Channel:   "#general",
		Username:  "portwatch",
		IconEmoji: ":satellite:",
	}
}

func validateRocketChat(cfg RocketChatConfig) error {
	if !cfg.Enabled {
		return nil
	}
	if cfg.WebhookURL == "" {
		return fmt.Errorf("rocketchat: webhook_url is required")
	}
	return nil
}
