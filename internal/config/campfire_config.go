package config

import "fmt"

// CampfireConfig holds settings for the Campfire alert handler.
type CampfireConfig struct {
	Enabled    bool   `toml:"enabled"`
	Token      string `toml:"token"`
	AccountID  string `toml:"account_id"`
	RoomID     string `toml:"room_id"`
	BaseURL    string `toml:"base_url"`
}

func defaultCampfireConfig() CampfireConfig {
	return CampfireConfig{
		Enabled:  false,
		BaseURL:  "https://api.campfirenow.com",
	}
}

func validateCampfire(cfg CampfireConfig) error {
	if !cfg.Enabled {
		return nil
	}
	if cfg.Token == "" {
		return fmt.Errorf("campfire: token is required")
	}
	if cfg.AccountID == "" {
		return fmt.Errorf("campfire: account_id is required")
	}
	if cfg.RoomID == "" {
		return fmt.Errorf("campfire: room_id is required")
	}
	if cfg.BaseURL == "" {
		return fmt.Errorf("campfire: base_url must not be empty")
	}
	return nil
}
