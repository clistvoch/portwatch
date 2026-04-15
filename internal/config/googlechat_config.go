package config

import "fmt"

type GoogleChatConfig struct {
	Enabled    bool   `toml:"enabled"`
	WebhookURL string `toml:"webhook_url"`
	ThreadKey  string `toml:"thread_key"`
}

func defaultGoogleChatConfig() GoogleChatConfig {
	return GoogleChatConfig{
		Enabled:   false,
		ThreadKey: "",
	}
}

func validateGoogleChat(cfg GoogleChatConfig) error {
	if !cfg.Enabled {
		return nil
	}
	if cfg.WebhookURL == "" {
		return fmt.Errorf("googlechat: webhook_url is required")
	}
	return nil
}
