package config

import "fmt"

// MSTeamsConfig holds configuration for Microsoft Teams adaptive card alerts.
type MSTeamsConfig struct {
	Enabled    bool   `toml:"enabled"`
	WebhookURL string `toml:"webhook_url"`
	Title      string `toml:"title"`
	ThemeColor string `toml:"theme_color"`
}

func defaultMSTeamsConfig() MSTeamsConfig {
	return MSTeamsConfig{
		Enabled:    false,
		WebhookURL: "",
		Title:      "Portwatch Alert",
		ThemeColor: "FF0000",
	}
}

func validateMSTeams(c MSTeamsConfig) error {
	if !c.Enabled {
		return nil
	}
	if c.WebhookURL == "" {
		return fmt.Errorf("msteams: webhook_url is required when enabled")
	}
	if c.Title == "" {
		return fmt.Errorf("msteams: title must not be empty")
	}
	if c.ThemeColor == "" {
		return fmt.Errorf("msteams: theme_color must not be empty")
	}
	return nil
}
