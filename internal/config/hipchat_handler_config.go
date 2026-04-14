package config

import "fmt"

// HipChatHandlerConfig holds runtime handler settings for HipChat alerts.
type HipChatHandlerConfig struct {
	RoomID  string `toml:"room_id"`
	Token   string `toml:"token"`
	Color   string `toml:"color"`
	Notify  bool   `toml:"notify"`
	BaseURL string `toml:"base_url"`
}

func defaultHipChatHandlerConfig() HipChatHandlerConfig {
	return HipChatHandlerConfig{
		Color:   "yellow",
		Notify:  false,
		BaseURL: "https://api.hipchat.com",
	}
}

func validateHipChatHandler(c HipChatHandlerConfig) error {
	if c.RoomID == "" {
		return fmt.Errorf("hipchat: room_id is required")
	}
	if c.Token == "" {
		return fmt.Errorf("hipchat: token is required")
	}
	validColors := map[string]bool{
		"yellow": true, "green": true, "red": true,
		"purple": true, "gray": true, "random": true,
	}
	if !validColors[c.Color] {
		return fmt.Errorf("hipchat: invalid color %q", c.Color)
	}
	if c.BaseURL == "" {
		return fmt.Errorf("hipchat: base_url must not be empty")
	}
	return nil
}
