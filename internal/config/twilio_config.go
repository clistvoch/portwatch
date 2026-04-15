package config

import "fmt"

type TwilioConfig struct {
	Enabled    bool   `toml:"enabled"`
	AccountSID string `toml:"account_sid"`
	AuthToken  string `toml:"auth_token"`
	From       string `toml:"from"`
	To         string `toml:"to"`
	BaseURL    string `toml:"base_url"`
}

func defaultTwilioConfig() TwilioConfig {
	return TwilioConfig{
		Enabled: false,
		BaseURL: "https://api.twilio.com",
	}
}

func validateTwilio(c TwilioConfig) error {
	if !c.Enabled {
		return nil
	}
	if c.AccountSID == "" {
		return fmt.Errorf("twilio: account_sid is required")
	}
	if c.AuthToken == "" {
		return fmt.Errorf("twilio: auth_token is required")
	}
	if c.From == "" {
		return fmt.Errorf("twilio: from number is required")
	}
	if c.To == "" {
		return fmt.Errorf("twilio: to number is required")
	}
	if c.BaseURL == "" {
		return fmt.Errorf("twilio: base_url must not be empty")
	}
	return nil
}
