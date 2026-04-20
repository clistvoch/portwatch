package config

import "fmt"

type SignalWireConfig struct {
	ProjectID string `toml:"project_id"`
	APIToken  string `toml:"api_token"`
	SpaceURL  string `toml:"space_url"`
	From      string `toml:"from"`
	To        string `toml:"to"`
}

func defaultSignalWireConfig() SignalWireConfig {
	return SignalWireConfig{
		SpaceURL: "https://example.signalwire.com",
	}
}

func validateSignalWire(c SignalWireConfig) error {
	if c.ProjectID == "" {
		return fmt.Errorf("signalwire: project_id is required")
	}
	if c.APIToken == "" {
		return fmt.Errorf("signalwire: api_token is required")
	}
	if c.SpaceURL == "" {
		return fmt.Errorf("signalwire: space_url is required")
	}
	if c.From == "" {
		return fmt.Errorf("signalwire: from number is required")
	}
	if c.To == "" {
		return fmt.Errorf("signalwire: to number is required")
	}
	return nil
}
