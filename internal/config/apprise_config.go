package config

import "fmt"

type AppriseConfig struct {
	Enabled  bool   `toml:"enabled"`
	URL      string `toml:"url"`
	Tag      string `toml:"tag"`
	Title    string `toml:"title"`
}

func defaultAppriseConfig() AppriseConfig {
	return AppriseConfig{
		Enabled: false,
		Tag:     "portwatch",
		Title:   "PortWatch Alert",
	}
}

func validateApprise(c AppriseConfig) error {
	if !c.Enabled {
		return nil
	}
	if c.URL == "" {
		return fmt.Errorf("apprise: url is required")
	}
	return nil
}
