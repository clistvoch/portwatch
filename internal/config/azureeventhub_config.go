package config

import "fmt"

type AzureEventHubConfig struct {
	ConnectionString string `toml:"connection_string"`
	EventHubName     string `toml:"event_hub_name"`
	Enabled          bool   `toml:"enabled"`
}

func defaultAzureEventHubConfig() AzureEventHubConfig {
	return AzureEventHubConfig{
		Enabled: false,
	}
}

func validateAzureEventHub(c AzureEventHubConfig) error {
	if !c.Enabled {
		return nil
	}
	if c.ConnectionString == "" {
		return fmt.Errorf("azure_event_hub: connection_string is required")
	}
	if c.EventHubName == "" {
		return fmt.Errorf("azure_event_hub: event_hub_name is required")
	}
	return nil
}
