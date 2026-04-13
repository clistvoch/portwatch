package config

import "fmt"

// PushoverConfig holds configuration for the Pushover notification handler.
type PushoverConfig struct {
	Enabled  bool   `toml:"enabled"`
	APIToken string `toml:"api_token"`
	UserKey  string `toml:"user_key"`
	Title    string `toml:"title"`
	Priority int    `toml:"priority"`
	Sound    string `toml:"sound"`
}

func defaultPushoverConfig() PushoverConfig {
	return PushoverConfig{
		Enabled:  false,
		Title:    "portwatch alert",
		Priority: 0,
		Sound:    "pushover",
	}
}

func validatePushover(cfg PushoverConfig) error {
	if !cfg.Enabled {
		return nil
	}
	if cfg.APIToken == "" {
		return fmt.Errorf("pushover: api_token is required")
	}
	if cfg.UserKey == "" {
		return fmt.Errorf("pushover: user_key is required")
	}
	if cfg.Priority < -2 || cfg.Priority > 2 {
		return fmt.Errorf("pushover: priority must be between -2 and 2, got %d", cfg.Priority)
	}
	return nil
}
