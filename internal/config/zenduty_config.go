package config

import "fmt"

type ZendutyConfig struct {
	Enabled   bool   `toml:"enabled"`
	APIKey    string `toml:"api_key"`
	ServiceID string `toml:"service_id"`
	AlertType string `toml:"alert_type"`
	Summary   string `toml:"summary"`
}

func defaultZendutyConfig() ZendutyConfig {
	return ZendutyConfig{
		Enabled:   false,
		AlertType: "critical",
		Summary:   "portwatch: port change detected",
	}
}

func validateZenduty(c ZendutyConfig) error {
	if !c.Enabled {
		return nil
	}
	if c.APIKey == "" {
		return fmt.Errorf("zenduty: api_key is required")
	}
	if c.ServiceID == "" {
		return fmt.Errorf("zenduty: service_id is required")
	}
	valid := map[string]bool{"critical": true, "warning": true, "info": true}
	if !valid[c.AlertType] {
		return fmt.Errorf("zenduty: invalid alert_type %q (must be critical, warning, or info)", c.AlertType)
	}
	return nil
}
