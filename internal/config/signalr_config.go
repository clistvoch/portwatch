package config

import "fmt"

// SignalRConfig holds configuration for Azure SignalR / ASP.NET SignalR alerting.
type SignalRConfig struct {
	Enabled    bool   `toml:"enabled"`
	Endpoint   string `toml:"endpoint"`
	AccessKey  string `toml:"access_key"`
	Hub        string `toml:"hub"`
	TimeoutSec int    `toml:"timeout_sec"`
}

func defaultSignalRConfig() SignalRConfig {
	return SignalRConfig{
		Enabled:    false,
		Hub:        "portwatch",
		TimeoutSec: 10,
	}
}

func validateSignalR(c SignalRConfig) error {
	if !c.Enabled {
		return nil
	}
	if c.Endpoint == "" {
		return fmt.Errorf("signalr: endpoint is required")
	}
	if c.AccessKey == "" {
		return fmt.Errorf("signalr: access_key is required")
	}
	if c.TimeoutSec <= 0 {
		return fmt.Errorf("signalr: timeout_sec must be positive")
	}
	return nil
}
