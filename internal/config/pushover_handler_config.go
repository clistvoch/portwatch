package config

import "fmt"

// PushoverHandlerConfig holds configuration for the Pushover alert handler.
type PushoverHandlerConfig struct {
	Enabled  bool   `toml:"enabled"`
	APIKey   string `toml:"api_key"`
	UserKey  string `toml:"user_key"`
	Priority int    `toml:"priority"`
	Title    string `toml:"title"`
	BaseURL  string `toml:"base_url"`
}

func defaultPushoverHandlerConfig() PushoverHandlerConfig {
	return PushoverHandlerConfig{
		Enabled:  false,
		Priority: 0,
		Title:    "portwatch alert",
		BaseURL:  "https://api.pushover.net/1/messages.json",
	}
}

func validatePushoverHandler(c PushoverHandlerConfig) error {
	if !c.Enabled {
		return nil
	}
	if c.APIKey == "" {
		return &ValidationError{Field: "pushover.api_key", Message: "api_key is required"}
	}
	if c.UserKey == "" {
		return &ValidationError{Field: "pushover.user_key", Message: "user_key is required"}
	}
	if c.Priority < -2 || c.Priority > 2 {
		return &ValidationError{
			Field:   "pushover.priority",
			Message: fmt.Sprintf("priority must be between -2 and 2, got %d", c.Priority),
		}
	}
	if c.BaseURL == "" {
		return &ValidationError{Field: "pushover.base_url", Message: "base_url must not be empty"}
	}
	return nil
}
