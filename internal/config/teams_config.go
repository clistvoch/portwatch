package config

import "fmt"

// TeamsConfig holds configuration for Microsoft Teams webhook alerts.
type TeamsConfig struct {
	Enabled    bool   `toml:"enabled"`
	WebhookURL string `toml:"webhook_url"`
	Title      string `toml:"title"`
}

func defaultTeamsConfig() TeamsConfig {
	return TeamsConfig{
		Enabled:    false,
		WebhookURL: "",
		Title:      "PortWatch Alert",
	}
}

func validateTeams(c TeamsConfig) error {
	if !c.Enabled {
		return nil
	}
	if c.WebhookURL == "" {
		return fmt.Errorf("teams: webhook_url is required when enabled")
	}
	if c.Title == "" {
		return fmt.Errorf("teams: title must not be empty")
	}
	return nil
}
