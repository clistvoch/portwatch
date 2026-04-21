package config

import "fmt"

// LineConfig holds configuration for the LINE Notify alert handler.
type LineConfig struct {
	Enabled     bool   `toml:"enabled"`
	Token       string `toml:"token"`
	APIURL      string `toml:"api_url"`
	MessagePrefix string `toml:"message_prefix"`
}

func defaultLineConfig() LineConfig {
	return LineConfig{
		Enabled:       false,
		Token:         "",
		APIURL:        "https://notify-api.line.me/api/notify",
		MessagePrefix: "[portwatch]",
	}
}

func validateLine(c LineConfig) error {
	if !c.Enabled {
		return nil
	}
	if c.Token == "" {
		return fmt.Errorf("line: token is required")
	}
	if c.APIURL == "" {
		return fmt.Errorf("line: api_url must not be empty")
	}
	return nil
}
