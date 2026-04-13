package config

import "fmt"

// HipChatConfig holds configuration for the HipChat alert handler.
type HipChatConfig struct {
	Enabled    bool   `toml:"enabled"`
	RoomID     string `toml:"room_id"`
	AuthToken  string `toml:"auth_token"`
	BaseURL    string `toml:"base_url"`
	Color      string `toml:"color"`
	Notify     bool   `toml:"notify"`
}

func defaultHipChatConfig() HipChatConfig {
	return HipChatConfig{
		Enabled:   false,
		BaseURL:   "https://api.hipchat.com",
		Color:     "yellow",
		Notify:    false,
	}
}

func validateHipChat(cfg HipChatConfig) error {
	if !cfg.Enabled {
		return nil
	}
	if cfg.RoomID == "" {
		return fmt.Errorf("hipchat: room_id is required")
	}
	if cfg.AuthToken == "" {
		return fmt.Errorf("hipchat: auth_token is required")
	}
	validColors := map[string]bool{
		"yellow": true, "green": true, "red": true,
		"purple": true, "gray": true, "random": true,
	}
	if !validColors[cfg.Color] {
		return fmt.Errorf("hipchat: invalid color %q", cfg.Color)
	}
	return nil
}
