package config

import "fmt"

type ZendutyConfig struct {
	Enabled    bool   `toml:"enabled"`
	APIKey     string `toml:"api_key"`
	ServiceID  string `toml:"service_id"`
	IntegrationID string `toml:"integration_id"`
	AlertType  string `toml:"alert_type"`
}

func defaultZendutyConfig() ZendutyConfig {
	return ZendutyConfig{
		Enabled:   false,
		AlertType: "critical",
	}
}

func validateZenduty(cfg ZendutyConfig) error {
	if !cfg.Enabled {
		return nil
	}
	if cfg.APIKey == "" {
		return &ValidationError{Field: "zenduty.api_key", Message: "api_key is required"}
	}
	if cfg.ServiceID == "" {
		return &ValidationError{Field: "zenduty.service_id", Message: "service_id is required"}
	}
	if cfg.IntegrationID == "" {
		return &ValidationError{Field: "zenduty.integration_id", Message: "integration_id is required"}
	}
	valid := map[string]bool{"critical": true, "warning": true, "info": true}
	if !valid[cfg.AlertType] {
		return &ValidationError{
			Field:   "zenduty.alert_type",
			Message: fmt.Sprintf("invalid alert_type %q: must be critical, warning, or info", cfg.AlertType),
		}
	}
	return nil
}
