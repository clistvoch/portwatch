package config

import "fmt"

// NtfyConfig holds configuration for the ntfy.sh push notification handler.
type NtfyConfig struct {
	Enabled  bool   `toml:"enabled"`
	ServerURL string `toml:"server_url"`
	Topic    string `toml:"topic"`
	Token    string `toml:"token"`
	Priority int    `toml:"priority"`
}

func defaultNtfyConfig() NtfyConfig {
	return NtfyConfig{
		Enabled:   false,
		ServerURL: "https://ntfy.sh",
		Topic:     "",
		Token:     "",
		Priority:  3,
	}
}

func validateNtfy(c NtfyConfig) error {
	if !c.Enabled {
		return nil
	}
	if c.Topic == "" {
		return &ValidationError{Field: "ntfy.topic", Msg: "topic is required when ntfy is enabled"}
	}
	if c.ServerURL == "" {
		return &ValidationError{Field: "ntfy.server_url", Msg: "server_url must not be empty"}
	}
	if c.Priority < 1 || c.Priority > 5 {
		return &ValidationError{Field: "ntfy.priority", Msg: fmt.Sprintf("priority must be between 1 and 5, got %d", c.Priority)}
	}
	return nil
}
