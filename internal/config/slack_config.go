package config

import "fmt"

// SlackConfig holds configuration for the Slack alert handler.
type SlackConfig struct {
	Enabled    bool   `toml:"enabled"`
	WebhookURL string `toml:"webhook_url"`
	Channel    string `toml:"channel"`
	Username   string `toml:"username"`
	IconEmoji  string `toml:"icon_emoji"`
}

func defaultSlackConfig() SlackConfig {
	return SlackConfig{
		Enabled:   false,
		Channel:   "#alerts",
		Username:  "portwatch",
		IconEmoji: ":lock:",
	}
}

func validateSlack(cfg SlackConfig) error {
	if !cfg.Enabled {
		return nil
	}
	if cfg.WebhookURL == "" {
		return fmt.Errorf("%w: slack.webhook_url is required when slack is enabled", ErrValidation)
	}
	return nil
}
