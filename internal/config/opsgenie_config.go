package config

import "fmt"

// OpsGenieConfig holds configuration for the OpsGenie alert handler.
type OpsGenieConfig struct {
	Enabled    bool   `toml:"enabled"`
	APIKey     string `toml:"api_key"`
	Team       string `toml:"team"`
	Priority   string `toml:"priority"`
	APIBaseURL string `toml:"api_base_url"`
}

func defaultOpsGenieConfig() OpsGenieConfig {
	return OpsGenieConfig{
		Enabled:    false,
		APIKey:     "",
		Team:       "",
		Priority:   "P3",
		APIBaseURL: "https://api.opsgenie.com",
	}
}

func validateOpsGenie(cfg OpsGenieConfig) error {
	if !cfg.Enabled {
		return nil
	}
	if cfg.APIKey == "" {
		return ValidationError{Field: "opsgenie.api_key", Message: "api_key is required when opsgenie is enabled"}
	}
	valid := map[string]bool{"P1": true, "P2": true, "P3": true, "P4": true, "P5": true}
	if !valid[cfg.Priority] {
		return ValidationError{
			Field:   "opsgenie.priority",
			Message: fmt.Sprintf("priority %q is invalid; must be one of P1-P5", cfg.Priority),
		}
	}
	if cfg.APIBaseURL == "" {
		return ValidationError{Field: "opsgenie.api_base_url", Message: "api_base_url must not be empty"}
	}
	return nil
}
