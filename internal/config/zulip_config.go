package config

import "fmt"

// ZulipConfig holds settings for the Zulip alert handler.
type ZulipConfig struct {
	Enabled  bool   `toml:"enabled"`
	BaseURL  string `toml:"base_url"`
	Email    string `toml:"email"`
	APIKey   string `toml:"api_key"`
	Stream   string `toml:"stream"`
	Topic    string `toml:"topic"`
}

func defaultZulipConfig() ZulipConfig {
	return ZulipConfig{
		Enabled: false,
		Stream:  "portwatch",
		Topic:   "port changes",
	}
}

func validateZulip(cfg ZulipConfig) error {
	if !cfg.Enabled {
		return nil
	}
	if cfg.BaseURL == "" {
		return fmt.Errorf("zulip: base_url is required")
	}
	if cfg.Email == "" {
		return fmt.Errorf("zulip: email is required")
	}
	if cfg.APIKey == "" {
		return fmt.Errorf("zulip: api_key is required")
	}
	if cfg.Stream == "" {
		return fmt.Errorf("zulip: stream is required")
	}
	return nil
}
