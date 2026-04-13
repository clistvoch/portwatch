package config

import (
	"fmt"
)

// RocketChatConfig holds settings for the Rocket.Chat alert handler.
type RocketChatConfig struct {
	Enabled    bool   `toml:"enabled"`
	WebhookURL string `toml:"webhook_url"`
}

func defaultRocketChatConfig() *RocketChatConfig {
	return &RocketChatConfig{
		Enabled:    false,
		WebhookURL: "",
	}
}

func validateRocketChat(c *RocketChatConfig) error {
	if c.WebhookURL == "" {
		return ValidationError{Field: "rocketchat.webhook_url", Err: fmt.Errorf("must not be empty when rocketchat is enabled")}
	}
	return nil
}
