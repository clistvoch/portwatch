package config

import "fmt"

type AWSEventBridgeConfig struct {
	Enabled    bool   `toml:"enabled"`
	Region     string `toml:"region"`
	BusName    string `toml:"bus_name"`
	Source     string `toml:"source"`
	DetailType string `toml:"detail_type"`
	AccessKey  string `toml:"access_key"`
	SecretKey  string `toml:"secret_key"`
}

func defaultAWSEventBridgeConfig() AWSEventBridgeConfig {
	return AWSEventBridgeConfig{
		Enabled:    false,
		Region:     "us-east-1",
		BusName:    "default",
		Source:     "portwatch",
		DetailType: "PortChange",
	}
}

func validateAWSEventBridge(c AWSEventBridgeConfig) error {
	if !c.Enabled {
		return nil
	}
	if c.Region == "" {
		return fmt.Errorf("awseventbridge: region is required")
	}
	if c.BusName == "" {
		return fmt.Errorf("awseventbridge: bus_name is required")
	}
	if c.AccessKey == "" {
		return fmt.Errorf("awseventbridge: access_key is required")
	}
	if c.SecretKey == "" {
		return fmt.Errorf("awseventbridge: secret_key is required")
	}
	return nil
}
